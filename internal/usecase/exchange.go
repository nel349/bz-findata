package usecase

import (
	"context"
	"fmt"
	"github.com/dmitryburov/go-coinbase-socket/internal/entity"
	"github.com/dmitryburov/go-coinbase-socket/internal/repository"
	"github.com/dmitryburov/go-coinbase-socket/pkg/logger"
	"time"
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
			e.logger.Info("Context cancelled, exiting Message")
			return ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				e.logger.Info("Channel closed, exiting Message")
				return nil
			}

			switch {
			case msg.Ticker != nil:
				if err := e.exchange.CreateTick(ctx, msg); err != nil {
					//TODO [critical] block - what's need?
					return err
				}

				e.logger.Info(
					fmt.Sprintf(
						"Inserted ticker %s > time:%d, bid:%f, ask:%f",
						msg.Ticker.Symbol,
						msg.Ticker.Timestamp,
						msg.Ticker.Bid,
						msg.Ticker.Ask,
					),
				)

			case msg.Order != nil:
				e.logger.Info(fmt.Sprintf("Received order in : %+v", msg.Order))
			default:
				e.logger.Info("Unknown message type")
			}
		case <-ticker.C:
			e.logger.Info("Still alive, waiting for data")
		}
	}
}
