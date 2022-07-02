package transports

import "time"

type Transport interface {
	WaitForConnection()
	RegisterConnHandler(handler func(conn Connection))
}

type Connection interface {
	Send([]byte)
	OnDisconnect(cb func())
	GetConnId() uint16
	Recv() chan Message
}

type Message struct {
	Body       []byte
	ReceivedAt time.Time
}
