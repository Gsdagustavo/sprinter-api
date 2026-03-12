package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/logger"
	"github.com/Gsdagustavo/sprinter-api/infrastructure"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const shutdownTime = time.Second * 15

func main() {
	err := logger.InitLogger()
	if err != nil {
		slog.Error("failed to initialize logger", logger.Err(err))
	}

	err = start()
	if err != nil {
		slog.Error("failed to start server", slog.String("cause", err.Error()))
		log.Fatal(err)
	}
}

type program struct {
	server *http.Server
	cfg    entities.Config
}

func (p *program) SetupServer(r *mux.Router) {
	err := infrastructure.SetupModules(r, p.cfg)
	if err != nil {
		slog.Error("failed to setup modules", slog.String("cause", err.Error()))
		return
	}

	corsOptions := handlers.AllowedOriginValidator(func(s string) bool {
		// todo implement a proper cors validation here
		return true
	})

	p.server = &http.Server{
		Addr: fmt.Sprintf(":%d", p.cfg.Server.Port),
		Handler: handlers.CORS(
			corsOptions,
			handlers.AllowedHeaders([]string{"Authorization", "Content-Type", "Accept"}),
			handlers.AllowedMethods([]string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete}),
			handlers.AllowCredentials(),
		)(r),
		ReadTimeout:       time.Second * 120,
		WriteTimeout:      time.Second * 120,
		ReadHeaderTimeout: time.Second * 2,
		IdleTimeout:       time.Second * 60,
	}
}

func (p *program) StartListening() {
	go func() {
		slog.Info("starting server", slog.String("address", p.server.Addr))
		err := p.server.ListenAndServe()
		if err != nil {
			slog.Error("failed to start server", slog.String("cause", err.Error()))
		}
	}()
}

func (p *program) run(block bool) error {
	p.SetupServer(mux.NewRouter())
	p.StartListening()

	if !block {
		return nil
	}

	// Block to let the server shut down gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTime)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err := p.server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown http server", slog.String("cause", err.Error()))
		return err
	}

	return nil
}
