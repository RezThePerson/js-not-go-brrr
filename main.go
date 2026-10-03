package main

import (
	"embed"
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

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.FileServer(http.FS(webFS)))
	r.HandleFunc("/", handlers.HomeHandler)
	// r.HandleFunc("/jump", handlers.JumpHandler)
	// r.HandleFunc("/game-stream", handlers.GameStreamHandler)
	// r.HandleFunc("/players-stream", handlers.PlayersStreamHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return http.ListenAndServe(":"+port, r)
}

func main() {
	if err := Serve(); err != nil {
		slog.Error("server encountered an error", "err", err)
		os.Exit(1)
	}
}
