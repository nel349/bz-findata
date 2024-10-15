package usecase

import (
	"context"
	"fmt"
	"time"
	"github.com/dmitryburov/go-coinbase-socket/internal/entity"
	"github.com/dmitryburov/go-coinbase-socket/internal/repository"
	"github.com/dmitryburov/go-coinbase-socket/pkg/logger"
)

type exchangeService struct {
	exchange repository.Exchange
	logger   logger.Logger
}

// NewExchangeService created exchange usecase
func NewExchangeService(
	exchange repository.Exchange,
	logger logger.Logger,
) *exchangeService {
	return &exchangeService{exchange, logger}
}

func (e *exchangeService) Tick(ctx context.Context, ch <-chan entity.Message) error {

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			e.logger.Info("Context cancelled, exiting Tick")
			return ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				e.logger.Info("Channel closed, exiting Tick")
				return nil
			}

			if err := e.exchange.CreateTick(ctx, msg); err != nil {
				//TODO [critical] block - what's need?
				return err
			}

			e.logger.Info(
				fmt.Sprintf(
					"writed ticker %s > time:%d, bid:%f, ask:%f",
					msg.Ticker.Symbol,
					msg.Ticker.Timestamp,
					msg.Ticker.Bid,
					msg.Ticker.Ask,
				),
			)
		case <-ticker.C:
			e.logger.Info("Still alive, waiting for data")
		}

	}
}
