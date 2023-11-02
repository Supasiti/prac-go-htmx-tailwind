package main

import (
	"context"
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
	mux := router.NewRouter()

	// For static files
	mux.Handle("/css/output.css", http.FileServer(http.FS(css)))

	server := &http.Server{Addr: ":8080", Handler: mux}
	handleSigTerms(server)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("ListenAndServe() error", err)
		os.Exit(1)
	}
}

// For gracefully shutdown server
func handleSigTerms(server *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		slog.Info("received SIGTERM, exiting...")

		err := server.Shutdown(context.Background())
		if err != nil {
			slog.Error("server.Shutdown() error", err)
		}

		os.Exit(1)
	}()
}
