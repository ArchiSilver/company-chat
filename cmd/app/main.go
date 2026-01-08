package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"company-chat/internal/auth"
	"company-chat/internal/config"
	"company-chat/internal/domain"
	"company-chat/internal/httputil"
	"company-chat/internal/logging"
	"company-chat/internal/metrics"
	"company-chat/internal/middleware"
	"company-chat/internal/repository"
	"company-chat/internal/ws"

	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"strings"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.LoadConfig()

	// инициализация структурированного логгера
	isProd := os.Getenv("APP_ENV") == "production"
	logging.Init(isProd)
	defer logging.Sync()

	// регистрация метрик Prometheus и стандартных коллекторов
	// включаем per-chat метрику, только если явно указано в окружении
	enableWSByChat := os.Getenv("METRICS_WS_BY_CHAT") == "1"
	if enableWSByChat {
		metrics.RegisterWithOptions(true)
	} else {
		metrics.Register()
	}

	mux := http.NewServeMux()
	mux.Handle("/", metrics.WrapWithRoute("root.hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "🚀 Company Chat API — hello")
	})))

	mux.Handle("/health", metrics.WrapWithRoute("root.health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})))

	// Init DB
	db, err := repository.NewDB(cfg.GetDBConnectionString())
	if err != nil {
		logging.L.Fatalf("db connect: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)
	msgRepo := repository.NewMessageRepository(db)
	refreshRepo := repository.NewRefreshTokenRepository(db)

	hub := ws.NewHub()

	// валидатор для структур запросов
	v := validator.New()

	// отдача статического веб-клиента
	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))

	// endpoint метрик
	mux.Handle("/metrics", promhttp.Handler())

	// Загрузка секрета JWT: сперва пытаемся файл, затем переменную окружения
	jwtSecret := ""
	if path := os.Getenv("JWT_SECRET_FILE"); path != "" {
		if data, err := os.ReadFile(path); err == nil {
			jwtSecret = strings.TrimSpace(string(data))
		} else {
			logging.L.Fatalf("cannot read JWT secret file: %v", err)
		}
	}
	if jwtSecret == "" {
		jwtSecret = os.Getenv("JWT_SECRET")
	}
	if jwtSecret == "" && os.Getenv("APP_ENV") == "production" {
		logging.L.Fatalf("JWT secret is required in production (set JWT_SECRET or JWT_SECRET_FILE)")
	}
	if jwtSecret == "" {
		jwtSecret = "dev-secret"
	}

	// переключатель окружения для безопасности cookie (переменная isProd определена выше)

	// Эндпоинты аутентификации
	register := func(name, pattern string, h func(http.ResponseWriter, *http.Request)) {
		mux.Handle(pattern, metrics.WrapWithRoute(name, http.HandlerFunc(h)))
	}

	register("root.register", "/api/auth/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.JSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		var req struct {
			Email    string `json:"email" validate:"required,email"`
			Username string `json:"username" validate:"required,min=2"`
			Password string `json:"password" validate:"required,min=6"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httputil.JSONError(w, http.StatusBadRequest, "bad request")
			return
		}
		if err := v.Struct(&req); err != nil {
			if ve, ok := err.(validator.ValidationErrors); ok {
				fields := httputil.FormatValidationErrors(ve)
				httputil.JSON(w, http.StatusBadRequest, map[string]interface{}{"error": "validation", "fields": fields})
			} else {
				httputil.JSONError(w, http.StatusBadRequest, err.Error())
			}
			return
		}
		passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			httputil.JSONError(w, http.StatusInternalServerError, "internal error")
			return
		}
		u := &domain.User{
			Email:        req.Email,
			Username:     req.Username,
			PasswordHash: string(passHash),
			Role:         "user",
		}
		if err := userRepo.Create(r.Context(), u); err != nil {
			httputil.JSONError(w, http.StatusInternalServerError, "could not create user: "+err.Error())
			return
		}
		// hide password
		u.PasswordHash = ""
		httputil.JSON(w, http.StatusCreated, u)
	})

	mux.Handle("/api/auth/login", metrics.WrapWithRoute("auth.login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.JSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		var req struct {
			Email    string `json:"email" validate:"required,email"`
			Password string `json:"password" validate:"required"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httputil.JSONError(w, http.StatusBadRequest, "bad request")
			return
		}
		if err := v.Struct(&req); err != nil {
			if ve, ok := err.(validator.ValidationErrors); ok {
				fields := httputil.FormatValidationErrors(ve)
				httputil.JSON(w, http.StatusBadRequest, map[string]interface{}{"error": "validation", "fields": fields})
			} else {
				httputil.JSONError(w, http.StatusBadRequest, err.Error())
			}
			return
		}
		u, err := userRepo.FindByEmail(r.Context(), req.Email)
		if err != nil {
			httputil.JSONError(w, http.StatusInternalServerError, "internal error")
			return
		}
		if u == nil {
			httputil.JSONError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
			httputil.JSONError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		token, err := auth.GenerateAccessToken(jwtSecret, u.ID, time.Minute*15)
		if err != nil {
			httputil.JSONError(w, http.StatusInternalServerError, "could not generate token")
			return
		}
		// create refresh token and save hashed
		refreshToken, err := auth.GenerateRandomToken(32)
		if err == nil {
			_ = refreshRepo.Save(r.Context(), u.ID, refreshToken, time.Now().Add(7*24*time.Hour))
			// set cookie (consider adding Secure/SameSite in prod)
			cookie := &http.Cookie{
				Name:     "refresh_token",
				Value:    refreshToken,
				HttpOnly: true,
				Path:     "/",
			}
			if isProd {
				cookie.Secure = true
				cookie.SameSite = http.SameSiteLaxMode
			}
			http.SetCookie(w, cookie)
		}
		httputil.JSON(w, http.StatusOK, map[string]string{"access_token": token})
	})))

	mux.Handle("/api/auth/refresh", metrics.WrapWithRoute("auth.refresh", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// try cookie or body
		cookie, _ := r.Cookie("refresh_token")
		var rt string
		if cookie != nil {
			rt = cookie.Value
		} else {
			var body struct {
				RefreshToken string `json:"refresh_token"`
			}
			_ = json.NewDecoder(r.Body).Decode(&body)
			rt = body.RefreshToken
		}
		if rt == "" {
			httputil.JSONError(w, http.StatusBadRequest, "no refresh token")
			return
		}
		// lookup user by token
		uid, err := refreshRepo.GetUserIDByToken(r.Context(), rt)
		if err != nil || uid == "" {
			httputil.JSONError(w, http.StatusUnauthorized, "invalid refresh")
			return
		}
		// issue new access token
		token, err := auth.GenerateAccessToken(jwtSecret, uid, time.Minute*15)
		if err != nil {
			httputil.JSONError(w, http.StatusInternalServerError, "could not generate token")
			return
		}
		// rotate refresh token: revoke old and issue new
		_ = refreshRepo.RevokeByHash(r.Context(), rt)
		newRT, err := auth.GenerateRandomToken(32)
		if err == nil {
			_ = refreshRepo.Save(r.Context(), uid, newRT, time.Now().Add(7*24*time.Hour))
			cookie := &http.Cookie{Name: "refresh_token", Value: newRT, HttpOnly: true, Path: "/"}
			if isProd {
				cookie.Secure = true
				cookie.SameSite = http.SameSiteLaxMode
			}
			http.SetCookie(w, cookie)
		}
		httputil.JSON(w, http.StatusOK, map[string]string{"access_token": token})
	})))

	mux.Handle("/api/auth/logout", metrics.WrapWithRoute("auth.logout", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refresh_token")
		if err != nil || cookie.Value == "" {
			httputil.JSONError(w, http.StatusBadRequest, "no refresh token")
			return
		}
		_ = refreshRepo.RevokeByHash(r.Context(), cookie.Value)
		// clear cookie
		cookieClear := &http.Cookie{Name: "refresh_token", Value: "", Path: "/", MaxAge: -1}
		if isProd {
			cookieClear.Secure = true
			cookieClear.SameSite = http.SameSiteLaxMode
		}
		http.SetCookie(w, cookieClear)
		httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})))

	// хелпер для извлечения user id из заголовка Authorization или параметра token
	getUserID := func(r *http.Request) (string, error) {
		// try header
		authHeader := r.Header.Get("Authorization")
		var tokenStr string
		if authHeader != "" {
			fmt.Sscanf(authHeader, "Bearer %s", &tokenStr)
		}
		if tokenStr == "" {
			// try query param
			tokenStr = r.URL.Query().Get("token")
		}
		if tokenStr == "" {
			return "", fmt.Errorf("no token")
		}
		uid, err := auth.ParseAndValidate(jwtSecret, tokenStr)
		if err != nil {
			return "", err
		}
		return uid, nil
	}

	// WebSocket-эндпоинт
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	mux.Handle("/ws", metrics.WrapWithRoute("ws.connect", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get token from query param
		token := r.URL.Query().Get("token")
		var uid string
		if token != "" {
			var err error
			uid, err = auth.ParseAndValidate(jwtSecret, token)
			if err != nil {
				httputil.JSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
		} else {
			// fallback to header
			var err error
			uid, err = getUserID(r)
			if err != nil {
				httputil.JSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
		}
		chatID := r.URL.Query().Get("chat_id")
		if chatID == "" {
			httputil.JSONError(w, http.StatusBadRequest, "chat_id required")
			return
		}
		// check participant membership: allow only participants to connect
		var member bool
		err = db.Pool.QueryRow(r.Context(), "SELECT EXISTS (SELECT 1 FROM chat_participants WHERE chat_id=$1 AND user_id=$2)", chatID, uid).Scan(&member)
		if err != nil {
			httputil.JSONError(w, http.StatusInternalServerError, "could not verify chat membership")
			return
		}
		if !member {
			httputil.JSONError(w, http.StatusForbidden, "not a participant of chat")
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logging.L.Warnf("upgrade: %v", err)
			return
		}
		go ws.HandleWS(hub, db, msgRepo, jwtSecret, conn, uid, chatID)
	})))

	// Chats and messages REST endpoints (simple routing)
	mux.Handle("/api/chats", metrics.WrapWithRoute("api.chats", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chats" {
			// path like /api/chats/{id}/messages handled below
			// try messages path
			if strings.HasPrefix(r.URL.Path, "/api/chats/") {
				rest := strings.TrimPrefix(r.URL.Path, "/api/chats/")
				parts := strings.Split(rest, "/")
				if len(parts) >= 2 && parts[1] == "messages" {
					chatID := parts[0]
					if r.Method == http.MethodPost {
						// post message
						uid, err := getUserID(r)
						if err != nil {
							httputil.JSONError(w, http.StatusUnauthorized, "unauthorized")
							return
						}
						var body struct {
							Content string `json:"content" validate:"required"`
						}
						if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
							httputil.JSONError(w, http.StatusBadRequest, "bad request")
							return
						}
						if err := v.Struct(&body); err != nil {
							if ve, ok := err.(validator.ValidationErrors); ok {
								fields := httputil.FormatValidationErrors(ve)
								httputil.JSON(w, http.StatusBadRequest, map[string]interface{}{"error": "validation", "fields": fields})
							} else {
								httputil.JSONError(w, http.StatusBadRequest, err.Error())
							}
							return
						}
						m := &domain.Message{ChatID: chatID, SenderID: uid, Content: body.Content}
						if err := msgRepo.Create(r.Context(), m); err != nil {
							httputil.JSONError(w, http.StatusInternalServerError, "internal error")
							metrics.ErrorsTotal.WithLabelValues("rest_save").Inc()
							return
						}
						// metrics: count & size
						metrics.MessagesSent.WithLabelValues("rest").Inc()
						metrics.MessageSize.Observe(float64(len(body.Content)))
						httputil.JSON(w, http.StatusCreated, m)
						return
					}
					if r.Method == http.MethodGet {
						// list messages
						limit := 50
						msgs, err := msgRepo.ListByChat(r.Context(), chatID, limit, 0)
						if err != nil {
							httputil.JSONError(w, http.StatusInternalServerError, "internal error")
							return
						}
						httputil.JSON(w, http.StatusOK, msgs)
						return
					}
				}
			}
			http.NotFound(w, r)
			return
		}
		if r.Method == http.MethodPost {
			uid, err := getUserID(r)
			if err != nil {
				httputil.JSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			var body struct {
				Name string `json:"name" validate:"required"`
				Type string `json:"type"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				httputil.JSONError(w, http.StatusBadRequest, "bad request")
				return
			}
			if err := v.Struct(&body); err != nil {
				if ve, ok := err.(validator.ValidationErrors); ok {
					fields := httputil.FormatValidationErrors(ve)
					httputil.JSON(w, http.StatusBadRequest, map[string]interface{}{"error": "validation", "fields": fields})
				} else {
					httputil.JSONError(w, http.StatusBadRequest, err.Error())
				}
				return
			}
			c := &domain.Chat{Name: body.Name, Type: body.Type, CreatedBy: uid}
			if err := chatRepo.Create(r.Context(), c); err != nil {
				httputil.JSONError(w, http.StatusInternalServerError, "internal error")
				return
			}
			httputil.JSON(w, http.StatusCreated, c)
			return
		}
		if r.Method == http.MethodGet {
			http.Error(w, "not implemented: list user chats", http.StatusNotImplemented)
			return
		}
	})))

	// Protected route
	mux.HandleFunc("/api/users/me", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httputil.JSONError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		var tokenStr string
		fmt.Sscanf(authHeader, "Bearer %s", &tokenStr)
		uid, err := auth.ParseAndValidate(jwtSecret, tokenStr)
		if err != nil {
			httputil.JSONError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		u, err := userRepo.FindByID(r.Context(), uid)
		if err != nil || u == nil {
			httputil.JSONError(w, http.StatusNotFound, "not found")
			return
		}
		u.PasswordHash = ""
		httputil.JSON(w, http.StatusOK, u)
	})

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: metrics.Middleware(middleware.Middleware(mux)),
	}

	// Graceful shutdown on interrupt
	go func() {
		logging.L.Infof("server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.L.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logging.L.Fatalf("Server Shutdown Failed:%+v", err)
	}
	logging.L.Info("Server exited properly")
}
