package coinbase

import (
	"context"
	"encoding/json"
	"fmt"
	// "time"

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
	*websocket.Conn
	heartbeatCounter int64
}

// NewCoinbaseClient init client for Coinbase
func NewCoinbaseClient(url, protocol, origin string) (*client, error) {
	if origin == "" || url == "" {
		return nil, fmt.Errorf("%s", ErrRequireConfigParameters)
	}

	conn, err := websocket.Dial(url, protocol, origin)
	if err != nil {
		return nil, err
	}

	return &client{conn, 0}, nil
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
				ProductIDs: []string{"ETH-BTC"},
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

	// Start a goroutine to handle incoming messages
	go func() {
		for {
			message, err := c.ReadData()
			if err != nil {
				fmt.Printf("Error reading data: %v\n", err)
				continue
			}

			// fmt.Println("Received message:", string(message))

			var response HeartbeatResponse
			if err := json.Unmarshal(message, &response); err != nil {
				fmt.Printf("Error unmarshaling response: %v\n", err)
				continue
			}

			switch response.Type {
			case Heartbeat.String():
				fmt.Printf("Received heartbeat: %v\n", response)
			default:
				// fmt.Printf("Received unknown message type: %v\n", response.Type)

			}
		}
	}()
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
