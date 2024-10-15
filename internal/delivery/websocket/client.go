package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dmitryburov/go-coinbase-socket/config"
	"github.com/dmitryburov/go-coinbase-socket/internal/entity"
	"github.com/dmitryburov/go-coinbase-socket/internal/usecase"
	"github.com/dmitryburov/go-coinbase-socket/pkg/exchange"
	"github.com/dmitryburov/go-coinbase-socket/pkg/exchange/coinbase"
	"github.com/dmitryburov/go-coinbase-socket/pkg/logger"
	"golang.org/x/sync/errgroup"
	"net"
	"strings"
	"sync"
)

type client struct {
	logger   logger.Logger
	conn     exchange.Manager
	uc       *usecase.Services
	products []string
	channels []string
}

// NewSocketClient init websocket client from delivery layout
func NewSocketClient(conn exchange.Manager, uc *usecase.Services, logger logger.Logger, cfg config.ExchangeConfig) (*client, error) {
	if len(cfg.Symbols) == 0 {
		return nil, errors.New("not found symbols for subscribes")
	}

	return &client{
		logger,
		conn,
		uc,
		cfg.Symbols,
		cfg.Channels,
	}, nil
}

// Run websocket listener
func (c *client) Run(ctx context.Context) error {

	var g = errgroup.Group{}
	var hMap = make(map[string]chan entity.Message)

	for _, symbol := range c.products {
		hMap[symbol] = make(chan entity.Message)

		// we should use different channels for ticker and order then merge them later

		// run readers
		g.Go(func() error {
			return c.uc.Exchange.Tick(ctx, hMap[symbol])
		})
	}

	// subscribe to products
	sData, _ := json.Marshal(map[string]interface{}{
		"type":        "subscribe",
		"product_ids": c.products,
		"channels":    c.channels,
	})
	_, err := c.conn.WriteData(sData)
	if err != nil {
		return err
	}

	message, err := c.conn.ReadData()
	if err != nil {
		c.logger.Error(err)
		return err
	}
	result, err := coinbase.ParseResponse(message)
	if err != nil {
		c.logger.Error(err)
		return err
	}

	switch v := result.(type) {
	case *coinbase.Response:
		if v.Type == coinbase.Error {
			c.logger.Fatal(fmt.Sprintf("Error: %s:%s", v.Message, v.Reason))
		}
		if v.Type == coinbase.Subscriptions {
			c.logger.Info(fmt.Sprintf("started subscription on products [%s]", strings.Join(c.products, ",")))
		}
	default:
		c.logger.Error(fmt.Sprintf("unknown response type: %T", result))
	}

	// Subscribe to heartbeats
	_, err = c.conn.SubscribeToHeartbeats()
	if err != nil {
		return fmt.Errorf("failed to subscribe to heartbeats: %w", err)
	}
	// writers
	for _, symbol := range c.products {
		g.Go(func() error {
			return c.responseReader(symbol, hMap)
		})
	}

	if err = g.Wait(); err != nil {
		return err
	}

	return nil
}

// responseReader write to symbol channel from response socket data
func (c *client) responseReader(symbol string, hMap map[string]chan entity.Message) error {
	var mu = sync.Mutex{}
	// var tickData *coinbase.Response

	for {
		message, err := c.conn.ReadData()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			c.logger.Error(err)
			continue
		}

		response, err := coinbase.ParseResponse(message)
		if err != nil {
			c.logger.Error(err)
			continue
		}

		switch r := response.(type) {
		case *coinbase.TickerResponse:
			ticker, err := r.ToTicker()
			if err != nil {
				c.logger.Error(err)
				continue
			}
			mu.Lock()
			hMap[symbol] <- entity.Message{Ticker: ticker}
			mu.Unlock()
		case *coinbase.ReceivedOrderResponse:
			order, err := r.ToReceivedOrder()
			if err != nil {
				c.logger.Error(err)
				continue
			}
			mu.Lock()
			hMap[symbol] <- entity.Message{Order: order}
			mu.Unlock()
			// Handle received order (you might need to create a new channel for orders)
		case *coinbase.Response:
			if r.Type == coinbase.Error {
				c.logger.Error(fmt.Errorf("API error: %s - %s", r.Message, r.Reason))
			}
		}
	}

	return nil
}
