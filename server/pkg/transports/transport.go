package transports

type Transport interface {
	WaitForConnection()
	RegisterConnHandler(handler func(conn Connection))
}

type Connection interface {
	Send([]byte)
	OnDisconnect(cb func())
	Recv() chan []byte
}
