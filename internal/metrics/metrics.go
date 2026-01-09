package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	MessagesSent = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "company_chat_messages_total",
			Help: "Total number of messages created",
		}, []string{"channel"},
	)
	MessageSize = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "company_chat_message_size_bytes",
			Help:    "Histogram of message sizes in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 2, 8),
		},
	)
	ErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "company_chat_errors_total",
			Help: "Total number of internal errors",
		}, []string{"area"},
	)
	// HTTP-метрики
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "company_chat_http_requests_total",
			Help: "Total number of HTTP requests",
		}, []string{"method", "path", "status"},
	)
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "company_chat_http_request_duration_seconds",
			Help:    "Histogram of HTTP request durations",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "path"},
	)

	// Метрики WebSocket
	ActiveWSConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "company_chat_ws_active_connections",
			Help: "Number of active WebSocket connections",
		},
	)
	// опциональная метрика по чатам (осторожно - высокая кардинальность)
	ActiveWSConnectionsByChat = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "company_chat_ws_active_connections_by_chat",
			Help: "Number of active WebSocket connections per chat",
		}, []string{"chat_id"},
	)
)

var enableWSByChat bool
var registered bool

// Register регистрирует метрики. По умолчанию per-chat метрика отключена.
func Register() {
	RegisterWithOptions(false)
}

// RegisterWithOptions регистрирует метрики; если enablePerChat=true,
// будет зарегистрирована метрика ActiveWSConnectionsByChat.
func RegisterWithOptions(enablePerChat bool) {
	if registered {
		// already registered in this process, skip
		return
	}
	// register collectors but ignore AlreadyRegisteredError
	tryRegister := func(c prometheus.Collector) {
		if err := prometheus.DefaultRegisterer.Register(c); err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); ok {
				// ignore duplicate
				return
			}
			// unexpected error
			panic(err)
		}
	}
	tryRegister(MessagesSent)
	tryRegister(MessageSize)
	tryRegister(ErrorsTotal)
	tryRegister(RequestsTotal)
	tryRegister(RequestDuration)
	tryRegister(ActiveWSConnections)
	if enablePerChat {
		tryRegister(ActiveWSConnectionsByChat)
	}
	enableWSByChat = enablePerChat
	// стандартные коллекторы
	tryRegister(collectors.NewGoCollector())
	tryRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registered = true
}

// WSByChatEnabled возвращает флаг, включена ли метрика по чатам
func WSByChatEnabled() bool {
	return enableWSByChat
}
