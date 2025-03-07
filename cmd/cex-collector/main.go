package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nel349/bz-findata/config"
	"github.com/nel349/bz-findata/internal/cex-collector"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-ctx.Done()
		cancel()
	}()

	cfg, err := config.NewConfig(ctx)
	if err != nil {
		log.Fatalf("failed config init: %v", err)
	}

	app.Run(ctx, cfg)
}
