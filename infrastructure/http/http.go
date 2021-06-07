package http

import (
	"github.com/Cliengo/acelle-mail/config"
	zap3 "github.com/Cliengo/acelle-mail/container/logger/factory/zap"
	middlewares2 "github.com/Cliengo/acelle-mail/infrastructure/http/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	conf "github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"time"
)

type Handlers struct {
	health      chi.Router
	integration chi.Router
}

func NewServerHandlers(health HealthHandler, integration IntegrationHandler) Handlers {
	return Handlers{
		health:      health.NewRouter(),
		integration: integration.NewRouter(),
	}
}

type Server struct {
	handlers Handlers
}

func New(handlers Handlers) Server {
	return Server{
		handlers: handlers,
	}
}

func (server Server) newRouter() *chi.Mux {
	return chi.NewRouter()
}

func (server Server) configMiddlewares(router *chi.Mux) {
	middlewares := []func(handler2 http.Handler) http.Handler{
		middleware.RequestID,
		middleware.RealIP,
		middleware.Recoverer,
		middleware.URLFormat,
		middlewares2.JwtAuth,
	}
	for _, item := range middlewares {
		router.Use(item)
	}
}

func requestLogger() func(next http.Handler) http.Handler {
	return chi.Chain(
		middleware.RequestID,
		handlerLogger(),
		middleware.Recoverer).Handler
}

func handlerLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			var requestID string
			if reqID := r.Context().Value(middleware.RequestIDKey); reqID != nil {
				requestID = reqID.(string)
			}

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			latency := time.Since(start)
			if conf.GetString(config.LoggerCode) == "zap" {
				fields := []zapcore.Field{
					zap.String("requestMethod", r.Method),
					zap.Int("requestStatus", ww.Status()),
					zap.String("requestPath", r.URL.Path),
					zap.String("remoteIP", r.RemoteAddr),
					zap.String("proto", r.Proto),
					zap.Int("bytes", ww.BytesWritten()),
					zap.Duration("elapsed", latency),
				}
				if requestID != "" {
					fields = append(fields, zap.String("requestID", requestID))
				}
				zap3.ZapLogger.Info("", fields...)
			}

		}
		return http.HandlerFunc(fn)
	}
}

func (server Server) addCors(router *chi.Mux) {
	//TODO: Se  puede estar modificando el allowed origin por una funcion mas pro
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Origin",
			"Authorization",
			"Content-Type",
		},
		ExposedHeaders: []string{},
		MaxAge:         300,
		Debug:          false,
	}))
}

func (server Server) GenerateServer() *http.Server {
	handler := server.newRouter()
	handler.Use(requestLogger())
	server.configMiddlewares(handler)
	server.addCors(handler)

	handler.Mount("/", server.handlers.health)
	handler.Mount("/integration", server.handlers.integration)

	//TODO: Modificar el puerto por una variable de entorno
	return &http.Server{
		Addr:         ":8081",
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

}

func (server Server) ListAndServe() error {
	srv := server.GenerateServer()
	err := srv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
