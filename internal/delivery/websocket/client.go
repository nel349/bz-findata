package websocket

import (
	// "bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"github.com/dmitryburov/go-coinbase-socket/config"
	"github.com/dmitryburov/go-coinbase-socket/internal/entity"
	"github.com/dmitryburov/go-coinbase-socket/internal/usecase"
	"github.com/dmitryburov/go-coinbase-socket/pkg/exchange"
	"github.com/dmitryburov/go-coinbase-socket/pkg/exchange/coinbase"
	"github.com/dmitryburov/go-coinbase-socket/pkg/logger"
	"golang.org/x/sync/errgroup"
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

		// run readers
		g.Go(func() error {
			return c.uc.Exchange.ProcessStream(ctx, hMap[symbol])
		})
	}

	auth := coinbase.NewAuth()
	signature, timestamp, err := auth.GenerateSignature()
	if err != nil {
		c.logger.Error("error generate signature: ", err)
		return err
	}
	// subscribe to products
	sData, _ := json.Marshal(map[string]interface{}{
		"type":        "subscribe",
		"product_ids": c.products,
		"channels":    c.channels,
		"signature":   signature,
		"timestamp":   timestamp,
		"key":         auth.Key,
		"passphrase":  auth.Passphrase,
	})

	_, err = c.conn.WriteData(sData)
	if err != nil {
		c.logger.Error(err)
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
		if v.Type == coinbase.Error.String() {
			c.logger.Fatal(fmt.Sprintf("Error: %s:%s", v.Message, v.Reason))
		}
		if v.Type == coinbase.Subscriptions.String() {
			c.logger.Info(fmt.Sprintf("started subscription on products [%s]", strings.Join(c.products, ",")))
		}
	default:
		c.logger.Error(fmt.Sprintf("unknown response type: %T", result))
	}

	// Subscribe to heartbeats
	c.conn.SubscribeToHeartbeats(ctx)

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

func (c *client) responseReader(symbol string, hMap map[string]chan entity.Message) error {
    var mu = sync.Mutex{}
    var accumulator string

    for {
        message, err := c.conn.ReadData()
        if err != nil {
            return fmt.Errorf("failed to read message: %w", err)
        }

        // Append the new message to the accumulator
        accumulator += string(message)

        // Try to extract and process complete JSON objects
        for {
            start := strings.Index(accumulator, "{")
            end := strings.LastIndex(accumulator, "}")

            if start == -1 || end == -1 || end < start {
                // No complete JSON object found
                break
            }

            // Extract the JSON object
            jsonStr := accumulator[start : end+1]
            accumulator = accumulator[end+1:]

            // Process the complete JSON message
            var rawJSON json.RawMessage
            err := json.Unmarshal([]byte(jsonStr), &rawJSON)
            if err != nil {
                c.logger.Error("Error unmarshalling JSON: ", err)
                continue
            }

            response, err := coinbase.ParseResponse(rawJSON)
            if err != nil {
                c.logger.Error("Failed to parse response: ", err)
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
            case *coinbase.OrderResponse:
                order, err := r.ToOrderResponse()
                if err != nil {
                    c.logger.Error(err)
                    continue
                }
                mu.Lock()
                hMap[symbol] <- entity.Message{Order: order}
                mu.Unlock()
            case *coinbase.Response:
                if r.Type == coinbase.Error.String() {
                    c.logger.Error(fmt.Errorf("API error: %s - %s", r.Message, r.Reason))
                }
            }
        }
    }
}
