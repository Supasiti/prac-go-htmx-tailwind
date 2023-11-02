package main

import (
	"embed"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/supasiti/prac-go-htmx-tailwind/internal/router"
)

var (
	//go:embed css/output.css
	css embed.FS
)

func main() {
	handleSigTerms()

	mux := router.NewRouter()

	// For static files
	mux.Handle("/css/output.css", http.FileServer(http.FS(css)))

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		slog.Error("Unable to create server", err)
		os.Exit(1)
	}
}

func handleSigTerms() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		slog.Info("received SIGTERM, exiting...")
		os.Exit(1)
	}()
}
