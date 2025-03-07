package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/nel349/bz-findata/internal/cex-collector/repository"
	"github.com/nel349/bz-findata/pkg/entity"
	"github.com/nel349/bz-findata/pkg/logger"
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

type ProductThreshold struct {
	ProductID string
	Threshold float64
}

var allowedOrderTypes = []string{
	"match",
	// "open",
	// "done",
	// "received",
}

var thresholds = []ProductThreshold{
	{ProductID: "ETH-USD", Threshold: 20000}, // 20k
	{ProductID: "BTC-USD", Threshold: 20000}, // 20k
}

func (e *exchangeService) shouldProcessOrder(order *entity.Order) bool {
	// Check if order type is allowed
	orderTypeAllowed := false
	for _, allowedType := range allowedOrderTypes {
		if strings.ToLower(order.Type) == allowedType {
			orderTypeAllowed = true
			break
		}
	}
	if !orderTypeAllowed {
		return false
	}

	orderValue := order.Size * order.Price

	for _, t := range thresholds {
		if order.ProductID == t.ProductID && orderValue > t.Threshold {
			return true
		}
	}
	return false
}

func (e *exchangeService) ProcessStream(ctx context.Context, ch <-chan entity.Message) error {

	// ticker := time.NewTicker(5 * time.Second)
	// defer ticker.Stop()
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
				if e.shouldProcessOrder(msg.Order) {
					e.logger.Info(fmt.Sprintf(
						"Received order: total_value:%f, type:%s, product_id:%s, size:%f, price:%f",
						msg.Order.Size*msg.Order.Price,
						msg.Order.Type,
						msg.Order.ProductID,
						msg.Order.Size,
						msg.Order.Price,
					))
					if err := e.exchange.CreateOrder(ctx, msg); err != nil {
						e.logger.Error(fmt.Sprintf("Failed to create order: %v", err))
						continue
					}
				}
			case msg.Heartbeat != nil:
				e.logger.Info(fmt.Sprintf("Received heartbeat in : %+v", msg.Heartbeat))
			default:
				e.logger.Info("Unknown message type")
			}
		// case <-ticker.C:
		// 	e.logger.Info("Still alive, waiting for data")
		}
	}
}
