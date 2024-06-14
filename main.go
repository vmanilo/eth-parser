package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/vmanilo/eth-parser/internal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go shutdownHook(cancel)

	srv := NewServer(internal.NewParser(ctx))

	log.Println("starting server on port 8080 ...")
	if err := srv.Serve(ctx, ":8080"); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("failed to start server", err)
	}

	log.Println("server stopped")
}

func shutdownHook(cancel context.CancelFunc) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-s

	cancel()
}
