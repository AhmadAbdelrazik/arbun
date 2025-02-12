package handlers

import (
	"AhmadAbdelrazik/arbun/internal/jsonlog"
	"AhmadAbdelrazik/arbun/internal/services"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Application struct {
	services *services.Services
	logger   *jsonlog.Log
	wg       sync.WaitGroup
}

func New() *Application {
	return &Application{
		services: services.New(),
		logger:   jsonlog.New(os.Stdout, jsonlog.LevelInfo),
	}
}

func (app *Application) Serve() error {
	srv := http.Server{
		// TODO: Make the address configurable
		Addr:         "localhost:3000",
		Handler:      app.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  1 * time.Minute,
	}

	shutdownerr := make(chan error, 1)
	go func() {

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		sig := <-quit
		app.logger.PrintInfo("Shutting down the server", map[string]string{
			"signal": sig.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownerr <- err
		}

		app.logger.PrintInfo("waiting for the background jobs", nil)

		app.wg.Wait()

		shutdownerr <- nil
	}()

	// TODO: make the development configurable
	app.logger.PrintInfo("starting the server", map[string]string{
		"Addr": srv.Addr,
		"env":  "development",
	})

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	if err := <-shutdownerr; err != nil {
		return err
	}

	app.logger.PrintInfo("stopped the server", nil)

	return nil
}
