package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"challenge/internal/counter"
	"challenge/pkg/logging"
)

var (
	// program flags
	listen       string
	cors         bool
	debug, trace bool
)

func init() {
	// setup flags
	flag.StringVar(&listen, "listen", "", "TCP Address to listen to incoming connections (format IP:Port)")
	flag.BoolVar(&cors, "cors", false, "Enable CORS support")
	flag.BoolVar(&debug, "debug", false, "Increases logging verbosity to DEBUG level")
	flag.BoolVar(&trace, "trace", false, "Increases logging verbosity to TRACE level")

	flag.Parse()

	// setup logging
	// Higher verbosity level takes precedence
	switch {
	case trace:
		logging.Setup(logging.TRACE)
	case debug:
		logging.Setup(logging.DEBUG)
	}
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
	logging.Infof("Starting server listening on address: %s, cors %t", cfg.Listen, cfg.Cors)
	// Start service
	if err := svc.Run(); err != nil {
		log.Printf("HTTP server stopped. Reason: %v\n", err)
	}
}
