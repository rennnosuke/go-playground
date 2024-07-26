package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// create mux
	mux := &http.ServeMux{}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// ...
	})

	svr := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// exec server
	if err := svr.ListenAndServe(); err != nil {
		fmt.Printf("Server.ListenAndServe error: %v\n", err)
	}

	// catch SIGTERM, SIGINT, SIGKILL
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	// sync until signal received
	<-ctx.Done()

	// shutting down server within 1min - otherwise force shutdown
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	fmt.Printf("Shutdown server...\n")

	if err := svr.Shutdown(ctx); err != nil {
		fmt.Printf("failed to shutdown server: %v\n", err)
	}
}
