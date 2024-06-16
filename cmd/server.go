package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/wlcmtunknwndth/AtomicHackBackend/internal/broker/nats"
	"github.com/wlcmtunknwndth/AtomicHackBackend/internal/config"
	"github.com/wlcmtunknwndth/AtomicHackBackend/internal/handlers/chat"
	"github.com/wlcmtunknwndth/AtomicHackBackend/internal/storage/pg"
	"github.com/wlcmtunknwndth/AtomicHackBackend/lib/slogResp"
	"log/slog"
	"net/http"
)

const scope = "main"

func main() {
	cfg := config.MustLoad()
	slog.Info("read config", slog.Any("cfg", cfg))

	pg, err := pg.New(&cfg.Db)
	if err != nil {
		slog.Error("couldn't run storage", slogResp.Error(scope, err))
		return
	}
	slog.Info("successfully connected to storage")

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Logger)

	srv := http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	fileHandler := http.StripPrefix(cfg.FileServer.UrlPath, http.FileServer(http.Dir(cfg.FileServer.StorageFolder)))

	router.Handle(fmt.Sprintf("/%s/%s", cfg.FileServer.UrlPath, "*"), fileHandler)

	broker, err := nats.New(&cfg.Broker, pg)
	if err != nil {
		slog.Error("couldn't run broker", slogResp.Error(scope, err))
		return
	}

	reqFinder, err := broker.RequestFinder()
	if err != nil {
		slog.Error("couldn't run request finder", slogResp.Error(scope, err))
		return
	}
	defer reqFinder.Unsubscribe()

	reqSaver, err := broker.RequestSaver()
	if err != nil {
		slog.Error("couldn't run request saver", slogResp.Error(scope, err))
		return
	}
	defer reqSaver.Unsubscribe()

	respFinder, err := broker.ResponseFinder()
	if err != nil {
		slog.Error("couldn't run request finder", slogResp.Error(scope, err))
		return
	}
	defer respFinder.Unsubscribe()

	respSaver, err := broker.ResponseSaver()
	if err != nil {
		slog.Error("couldn't run request saver", slogResp.Error(scope, err))
		return
	}
	defer respSaver.Unsubscribe()

	handler := chat.New(2048, 2048, broker, cfg.ReceiverAddr)
	router.HandleFunc("/front-ws", handler.HandleConnections)

	if err = srv.ListenAndServe(); err != nil {
		slog.Error("failed to run server", slogResp.Error(scope, err))
	}
	slog.Info("server closed")
}
