package metrics

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

// ResponseRecorder оборачивает http.ResponseWriter для сохранения кода статуса
type ResponseRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *ResponseRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

// Flush реализует интерфейс http.Flusher если базовый ResponseWriter поддерживает
func (r *ResponseRecorder) Flush() {
	if f, ok := r.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack реализует интерфейс http.Hijacker если базовый ResponseWriter поддерживает
func (r *ResponseRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("hijack not supported")
	}
	return hj.Hijack()
}

type routeKeyType string

var routeKey = routeKeyType("route_name")

// RouteContext возвращает context с установленным именем маршрута
func RouteContext(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, routeKey, name)
}

// getRouteFromContext возвращает имя маршрута, если оно присутствует
func getRouteFromContext(r *http.Request) string {
	if v := r.Context().Value(routeKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// WrapWithRoute возвращает handler, который устанавливает имя маршрута в контекст запроса
func WrapWithRoute(name string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := RouteContext(r.Context(), name)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Middleware инструментирует все HTTP-запросы метриками Prometheus.
// Он использует имя маршрута из контекста (если задано), иначе путь.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rr := &ResponseRecorder{ResponseWriter: w, Status: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(rr, r)
		dur := time.Since(start).Seconds()
		statusLabel := "unknown"
		if rr.Status != 0 {
			statusLabel = http.StatusText(rr.Status)
		}
		path := r.URL.Path
		if route := getRouteFromContext(r); route != "" {
			path = route
		}
		method := r.Method
		RequestsTotal.WithLabelValues(method, path, statusLabel).Inc()
		RequestDuration.WithLabelValues(method, path).Observe(dur)
	})
}
