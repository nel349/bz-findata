package websocket

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/nel349/bz-findata/config"
	"github.com/nel349/bz-findata/internal/cex-collector/usecase"
	"github.com/nel349/bz-findata/pkg/entity"
	"github.com/nel349/bz-findata/pkg/exchange"
	"github.com/nel349/bz-findata/pkg/exchange/coinbase"
	"github.com/nel349/bz-findata/pkg/logger"
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

	// monitor heartbeat
	go c.conn.MonitorHeartbeat(ctx, 10*time.Second)

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
	var buffer []byte

	for {
		message, err := c.conn.ReadData()
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}

		// Append new data to buffer
		buffer = append(buffer, message...)

		// Process complete JSON objects from buffer
		for len(buffer) > 0 {
			// Find first opening brace
			start := bytes.IndexByte(buffer, '{')
			if start == -1 {
				buffer = nil
				break
			}

			// Find matching closing brace
			end := -1
			depth := 0
			for i := start; i < len(buffer); i++ {
				switch buffer[i] {
				case '{':
					depth++
				case '}':
					depth--
					if depth == 0 {
						end = i + 1
					}
				}
			}

			if end == -1 {
				// No complete JSON object found
				if start > 0 {
					buffer = buffer[start:]
				}
				break
			}

			// Extract and parse the complete JSON object
			jsonData := buffer[start:end]
			buffer = buffer[end:]

			var rawJSON json.RawMessage
			if err := json.Unmarshal(jsonData, &rawJSON); err != nil {
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
			case *coinbase.HeartbeatResponse:
				// update heartbeat
				c.conn.UpdateHeartbeat()
				heartbeat, err := r.ToHeartbeat()
				if err != nil {
					c.logger.Error(err)
					continue
				}
				mu.Lock()
				hMap[symbol] <- entity.Message{Heartbeat: heartbeat}
				mu.Unlock()
			case *coinbase.Response:
				if r.Type == coinbase.Error.String() {
					c.logger.Error(fmt.Errorf("API error: %s - %s", r.Message, r.Reason))
				}
			}
		}
	}
}
