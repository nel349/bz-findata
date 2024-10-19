package exchange

import "context"

// Manager is an interface exchange of application
type Manager interface {
	// SubscribeToHeartbeats is subscribing to heartbeat messages
	SubscribeToHeartbeats(ctx context.Context)
	// CloseConnection is closing connection
	CloseConnection() error
	// WriteData command write data to exchange connection
	WriteData(message []byte) (int, error)
	// ReadData command is reading from receiver data
	ReadData() ([]byte, error)
}
