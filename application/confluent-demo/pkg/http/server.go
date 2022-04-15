package http

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dueruen/learning-kubernetes/application/confluent-demo/pkg/middleware"
	"github.com/gorilla/mux"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	httpPort = flag.String("http-port", os.Getenv("HTTP_PORT"), "HTTP_PORT")
)

type Server struct {
	router         *mux.Router
	logger         *zap.Logger
	handler        http.Handler
	tracer         trace.Tracer
	tracerProvider *sdktrace.TracerProvider
}

func NewServer(logger *zap.Logger) (*Server, error) {
	srv := &Server{
		router: mux.NewRouter(),
		logger: logger,
	}

	return srv, nil
}

func (s *Server) registerHandlers() {
	s.router.HandleFunc("/healthz", s.healthzHandler).Methods("GET")
	s.router.HandleFunc("/readyz", s.readyzHandler).Methods("GET")
	s.router.HandleFunc("/long", s.longRunningProcess).Methods("GET")
}

func (s *Server) registerMiddlewares() {
	//otel := NewOpenTelemetryMiddleware()
	//s.router.Use(otel)
	httpLogger := middleware.NewLoggingMiddleware(s.logger)
	s.router.Use(httpLogger.Handler)
}

func (s *Server) ListenAndServe() chan bool {
	shutdown := make(chan bool)
	go func() {
		s.logger.Info("In listen and serve")

		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		// ctx := context.Background()

		tracerProvider, tracer := middleware.InitTracer(s.logger)
		s.tracerProvider = tracerProvider
		s.tracer = tracer

		s.registerHandlers()
		s.registerMiddlewares()

		s.handler = s.router

		srv := s.startServer()

		<-sigchan

		// s.logger.Info("Shutting down HTTP/HTTPS server", zap.Duration("timeout", s.config.HttpServerShutdownTimeout))

		// stop OpenTelemetry tracer provider
		// if s.tracerProvider != nil {
		// 	if err := s.tracerProvider.Shutdown(ctx); err != nil {
		// 		s.logger.Warn("stopping tracer provider", zap.Error(err))
		// 	}
		// }

		// determine if the http server was started
		if srv != nil {
			//if err := srv.Shutdown(ctx); err != nil {
			//s.logger.Warn("HTTP server graceful shutdown failed", zap.Error(err))
			s.logger.Warn("HTTP server graceful shutdown")
			//}
		}
		shutdown <- true
	}()

	s.logger.Info("HTTP shutdown")
	return shutdown
}

func (s *Server) startServer() *http.Server {

	// determine if the port is specified
	if httpPort == nil {

		// move on immediately
		return nil
	}

	srv := &http.Server{
		Addr:    ":" + *httpPort,
		Handler: s.handler,
	}

	// start the server in the background
	go func() {
		s.logger.Info("Starting HTTP Server.", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Fatal("HTTP server crashed", zap.Error(err))
		}
	}()

	// return the server and routine
	return srv
}
