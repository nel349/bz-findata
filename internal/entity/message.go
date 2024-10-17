package entity

type Message struct {
	Ticker *Ticker
	Order  *Order
	Heartbeat *Heartbeat
}