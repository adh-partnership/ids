package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adh-partnership/ids/backend/internal/middleware/logger"
	adhlog "github.com/adh-partnership/ids/backend/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	// HTTP/1.1 without TLS
	MODE_PLAIN = "plain"
	// HTTP/1.1 and HTTP/2 with TLS
	MODE_TLS = "tls"
	// HTTP/1.1 and HTTP/2 without TLS
	MODE_H2C = "h2c"
)

var log zerolog.Logger

type Server struct {
	Ctx    context.Context
	Router *chi.Mux
}

func New() *Server {
	router := chi.NewRouter()

	log = adhlog.ZL.With().Str("component", "server").Logger()
	router.Use(logger.Logger(adhlog.ZL.With().Str("component", "access").Logger())) // Includes RequestID and Recoverer middleware
	router.Use(middleware.RealIP)

	return &Server{
		Router: router,
	}
}

// Start the server and sets up graceful shutdown. Expects mode
// to be one of MODE_PLAIN, MODE_TLS, or MODE_H2C. If mode is invalid,
// will assume MODE_PLAIN which may not be the best default, but the preferred
// would require a TLS certificate and key to be set.
//
// Will block.
func (s *Server) Start(mode string, addr string) error {
	var cert string
	var key string
	var err error
	var serverStop context.CancelFunc

	srv := &http.Server{
		Addr:    addr,
		Handler: s.Router,
	}

	log.Info().Msgf("Starting server in %s mode on %s", mode, addr)

	s.Ctx, serverStop = context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig // Block until a signal is received.
		log.Info().Msg("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Warn().Msg("Graceful shutdown timed out. Forcing exit...")
			}
		}()
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			log.Error().Msgf("Error shutting down server: %s", err)
		}
		cancel()
		serverStop()
	}()

	if mode == MODE_TLS {
		cert = os.Getenv("SSL_CERT")
		key = os.Getenv("SSL_KEY")
		if cert == "" || key == "" {
			return errors.New("SSL_CERT and SSL_KEY environment variables must be set")
		}
		if _, err := os.Stat(cert); os.IsNotExist(err) {
			return errors.New("SSL_CERT file " + cert + " does not exist")
		}
		if _, err := os.Stat(key); os.IsNotExist(err) {
			return errors.New("SSL_KEY file " + key + " does not exist")
		}
	} else if mode == MODE_H2C {
		h2s := &http2.Server{}
		srv.Handler = h2c.NewHandler(s.Router, h2s)
	}

	if mode == MODE_TLS {
		err = srv.ListenAndServeTLS(cert, key)
	} else {
		err = srv.ListenAndServe()
	}
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error().Msgf("Error starting server: %s", err)
		return err
	}

	// Block until server is done
	<-s.Ctx.Done()

	return nil
}
