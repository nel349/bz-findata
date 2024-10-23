package entity

import "time"

// Heartbeat model data of response from exchange
type Heartbeat struct {
	Type        string    
	Sequence    int64     
	LastTradeID int64     
	ProductID   string    
	Time        time.Time 
}
