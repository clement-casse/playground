package web

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

const (
	defaultReadTimeout  = 30 * time.Second
	defaultWriteTimeout = 30 * time.Second
	defaultIdleTimeout  = 60 * time.Second
)

// Server manages HTTP requests
type Server struct {
	server      *http.Server
	mainHandler http.Handler

	logger *slog.Logger
}

// ServerOpt in an interface for applying Server options.
type ServerOpt interface {
	applyOpt(*Server) *Server
}

type serverOptFunc func(*Server) *Server

func (fn serverOptFunc) applyOpt(s *Server) *Server {
	return fn(s)
}

type Middleware interface {
	Handle(next http.Handler) http.Handler
}

// NewServer creates a simple HTTP server instance
func NewServer(addr string, mux http.Handler, opts ...ServerOpt) *Server {
	srv := &Server{
		server: &http.Server{
			Addr: addr,

			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
			IdleTimeout:  defaultIdleTimeout,
		},
		mainHandler: mux,
		logger:      slog.Default(),
	}
	for _, opt := range opts {
		srv = opt.applyOpt(srv)
	}
	return srv
}

// WithLogger modifies the current instance to add a custom logger, if the logger is nil, slog.Default() is used.
func WithLogger(l *slog.Logger) ServerOpt {
	return serverOptFunc(func(s *Server) *Server {
		if l == nil {
			l = slog.Default() // logger cannot be nil
		}
		s.logger = l
		s.server.ErrorLog = slog.NewLogLogger(l.Handler(), slog.LevelError)
		return s
	})
}

// WithReadTimeout applies a custom read timeout to the server.
func WithReadTimeout(t time.Duration) ServerOpt {
	return serverOptFunc(func(s *Server) *Server {
		s.server.ReadTimeout = t
		return s
	})
}

// WithWriteTimeout applies a custom write timeout to the server.
func WithWriteTimeout(t time.Duration) ServerOpt {
	return serverOptFunc(func(s *Server) *Server {
		s.server.WriteTimeout = t
		return s
	})
}

// WithIdleTimeout applies a custom idle timeout to the server.
func WithIdleTimeout(t time.Duration) ServerOpt {
	return serverOptFunc(func(s *Server) *Server {
		s.server.IdleTimeout = t
		return s
	})
}

// WithMiddlewares applies in order the middlewares provided as argument to the main Handler.
// These middlewares do not apply to the base functions exposed by the Server like `/health`.
func WithMiddlewares(mws ...Middleware) ServerOpt {
	return serverOptFunc(func(s *Server) *Server {
		for _, mw := range mws {
			origHandler := s.mainHandler
			s.mainHandler = mw.Handle(origHandler)
		}
		return s
	})
}

// StartServer starts the web server
func (s *Server) StartServer(ctx context.Context) error {
	s.server.Handler = s.makeHandler(s.mainHandler)
	s.logger.InfoContext(ctx, "Starting to serve requests on "+s.server.Addr)

	err := s.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		s.logger.WarnContext(ctx, "Server stopped")
		return nil
	}
	s.logger.ErrorContext(ctx, "Server error: "+err.Error())
	return err
}

// Shutdown stops the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) makeHandler(h http.Handler) http.Handler {
	mux := http.DefaultServeMux
	mux.HandleFunc("/health", s.healthHandlerFunc)
	mux.Handle("/", h)
	return mux
}

func (s *Server) healthHandlerFunc(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
