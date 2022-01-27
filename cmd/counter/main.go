package main

import (
	"challenge/internal/counter"
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
)

var (
	// program flags
	listen string
	cors   bool
)

func init() {
	// setup flags
	flag.StringVar(&listen, "listen", "", "TCP Address to listen to incoming connections (format IP:Port)")
	flag.BoolVar(&cors, "cors", false, "Enable CORS support")

	flag.Parse()
}

func main() {
	// Load service config
	cfg := counter.DefaultConfig
	if listen != "" {
		cfg.Listen = listen
	}
	cfg.Cors = cors
	// Instantiate service
	svc := counter.NewService(cfg)
	// setup stop goroutine
	go func() {
		// setup channel to receive interruption
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)
		// wait for system interruption
		<-stop
		// stop service
		svc.Shutdown(context.TODO())
	}()
	// Start service
	if err := svc.Run(); err != nil {
		log.Printf("HTTP server stopped. Reason: %v\n", err)
	}

}
