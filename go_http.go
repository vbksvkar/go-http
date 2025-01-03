package gohttp

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	gologger "github.com/vbksvkar/go-logger"
	"go.uber.org/zap"
)

var shutdown = make(chan os.Signal, 1)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func StartHttpServer(serviceName string, version string, logger *zap.SugaredLogger, server *http.Server) error {
	if logger == nil {
		l, err := gologger.New(serviceName, version)
		if err != nil {
			return err
		}
		logger = l
	}

	signal.Notify(shutdown, os.Interrupt)
	chiRouter, ok := server.Handler.(*chi.Mux)
	if !ok {
		return errors.New("server.Handler is not of type *chi.Mux")
	}

	chiRouter.Get("/ping", PingHandler)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(err)
			shutdown <- os.Interrupt
		}
	}()

	logger.Info("server started")
	oscall := <-shutdown
	logger.Infof("received signal: %v", oscall)

	ctxshut, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	err := server.Shutdown(ctxshut)
	if err != nil && err != http.ErrServerClosed {
		logger.Error("failed to shutdown server", "error", err)
		return err
	}
	logger.Info("server shutdown")
	return nil
}
