package main

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"github.com/RezThePerson/js-not-go-brrr/handlers"
	"github.com/gorilla/mux"
)

//go:embed web/*
var webFS embed.FS

func Serve() error {
	if err := handlers.InitTemplates(webFS); err != nil {
		return err
	}

	// Strip the "web/" prefix so /static/style.css maps to web/static/style.css.
	staticFS, err := fs.Sub(webFS, "web")
	if err != nil {
		return err
	}

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.FileServer(http.FS(staticFS)))
	r.HandleFunc("/", handlers.HomeHandler)
	r.HandleFunc("/jump", handlers.JumpHandler).Methods(http.MethodPost)
	r.HandleFunc("/stream", handlers.StreamHandler).Methods(http.MethodGet)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("starting server", "port", port)
	return http.ListenAndServe(":"+port, r)
}

func main() {
	if err := Serve(); err != nil {
		slog.Error("server encountered an error", "err", err)
		os.Exit(1)
	}
}
