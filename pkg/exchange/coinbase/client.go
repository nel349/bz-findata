package coinbase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nel349/bz-findata/config"
	"golang.org/x/net/websocket"
)

const (
	ErrRequireConfigParameters = "not correct input parameters"
)

type Subscribe struct {
	Type       string   `json:"type"`
	ProductIDs []string `json:"product_ids"`
	Channels   []string `json:"channels"`
}

type SubscribeHeartbeat struct {
	Type     string `json:"type"`
	Channels []struct {
		Name       string   `json:"name"`
		ProductIDs []string `json:"product_ids"`
	} `json:"channels"`
}

type client struct {
	cfg *config.Config
	*websocket.Conn
	lastHeartbeat     time.Time
	reconnectAttempts int
	lastReconnectTime time.Time
}

// NewCoinbaseClient init client for Coinbase
func NewCoinbaseClient(cfg *config.Config) (*client, error) {
	if cfg.Exchange.Origin == "" || cfg.Exchange.Url == "" {
		return nil, fmt.Errorf("%s", ErrRequireConfigParameters)
	}

	conn, err := websocket.Dial(cfg.Exchange.Url, cfg.Exchange.Protocol, cfg.Exchange.Origin)
	if err != nil {
		return nil, err
	}

	return &client{cfg, conn, time.Now(), 0, time.Now()}, nil
}

func (c *client) SubscribeToHeartbeats(ctx context.Context) {

	fmt.Println("Subscribing to heartbeats...")
	subscribeMsg := SubscribeHeartbeat{
		Type: "subscribe",
		Channels: []struct {
			Name       string   `json:"name"`
			ProductIDs []string `json:"product_ids"`
		}{
			{
				Name:       "heartbeat",
				ProductIDs: c.cfg.Exchange.Symbols,
			},
		},
	}

	subscribeBytes, err := json.Marshal(subscribeMsg)
	if err != nil {
		fmt.Println("error marshaling subscribe message: %w", err)
		return
	}

	// Send initial subscription message
	_, err = c.WriteData(subscribeBytes)
	if err != nil {
		fmt.Printf("error writing initial subscribe message: %v\n", err)
		return
	}
}

// Add a method to check and handle heartbeat timeout
func (c *client) MonitorHeartbeat(ctx context.Context, timeout time.Duration) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if time.Since(c.lastHeartbeat) > timeout {

				// Increment attempts before trying to reconnect
				c.reconnectAttempts++

				// check if we've exceeded the retry limit first
				if c.reconnectAttempts >= 5 {
					fmt.Println("Too many reconnection attempts, circuit breaker activated")
					time.Sleep(100 * time.Millisecond)
					panic("Circuit breaker activated due to too many reconnection attempts")
				}

				fmt.Printf("Heartbeat timeout detected, time since last heartbeat: %s\n", time.Since(c.lastHeartbeat))
				fmt.Printf("Time since last reconnection: %s\n", time.Since(c.lastReconnectTime))
				fmt.Printf("Reconnection attempts: %d\n", c.reconnectAttempts)

				// Add backoff between reconnection attempts
				backoff := time.Duration(c.reconnectAttempts+1) * time.Second
				time.Sleep(backoff)

				if err := c.reconnect(ctx, c.cfg); err != nil {
					c.reconnectAttempts++
					c.lastReconnectTime = time.Now()
					log.Printf("Reconnection failed: %v", err)
					continue
				}

				// Only reset attempts on successful reconnection
				c.lastReconnectTime = time.Now()
			}
		}
	}
}

// Add method to update last heartbeat time
func (c *client) UpdateHeartbeat() {
	c.lastHeartbeat = time.Now()
	// Reset attempts after successful heartbeat
	if c.reconnectAttempts > 0 {
		log.Printf("Connection stabilized, resetting reconnection attempts")
		c.reconnectAttempts = 0
	}
}

// Add reconnect method
func (c *client) reconnect(ctx context.Context, cfg *config.Config) error {
	if err := c.CloseConnection(); err != nil {
		log.Printf("Warning: error closing connection: %v", err)
	}

	fmt.Println("Reconnecting...")

	conn, err := websocket.Dial(cfg.Exchange.Url, cfg.Exchange.Protocol, cfg.Exchange.Origin)
	if err != nil {
		return fmt.Errorf("reconnection failed: %w", err)
	}

	c.Conn = conn

	// Pass context through subscription
	c.SubscribeToHeartbeats(ctx)

	// Update heartbeat time after successful connection
	c.lastHeartbeat = time.Now()

	return nil
}

func (c *client) WriteData(message []byte) (int, error) {
	return c.Write(message)
}

func (c *client) ReadData() ([]byte, error) {
	var message = make([]byte, 512) //TODO need change global? 1MB

	n, err := c.Read(message)
	if err != nil {
		return nil, err
	}

	return message[:n], nil
}

func (c *client) CloseConnection() error {
	return c.Close()
}
