package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/renjugeo/go-server/api"
	"github.com/renjugeo/go-server/config"
	"github.com/renjugeo/go-server/server"
	"github.com/renjugeo/go-server/util"
	"go.uber.org/zap"
)

func main() {
	os.Exit(run())
}

func run() int {
	var logger *zap.Logger
	var (
		term = make(chan os.Signal, 1)
		srvc = make(chan int, 1)
	)

	// Parse config file
	cfgPath := flag.String("config-path", "config.yaml", "configuration file")
	flag.Parse()

	cfg, err := config.ParseConfig(*cfgPath)
	if err != nil {
		panic(err)
	}

	// Set up logger
	logger, err = util.GetLogger(cfg)
	if err != nil {
		panic(err)
	}

	// Listen for interrupts
	signal.Notify(term, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	api := api.NewAPI(cfg, logger)

	// create server instance
	srv := server.NewServer(api, cfg, logger)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			close(srvc)
		}
	}()

	for {
		select {
		case <-term:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer func() {
				// extra handling here
				cancel()
			}()
			if err := srv.Shutdown(ctx); err != nil {
				logger.Debug("error shuting down server")
				return 1
			}
			logger.Debug("shutting down server...")
			return 0
		case <-srvc:
			if err := srv.Close(); err != nil {
				logger.Debug("error closing server")
			}
			return 1
		}
	}
}
