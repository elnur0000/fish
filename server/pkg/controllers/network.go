package controllers

import "github.com/elnur0000/fish-backend/pkg/transports"

type NetworkController struct {
	transport   transports.Transport
	connections []transports.Connection
}

func NewNetworkController(transport transports.Transport) NetworkController {
	return NetworkController{
		transport: transport,
	}
}

func (nc *NetworkController) ListenForConnections(onConnect func(conn transports.Connection)) {
	nc.transport.RegisterConnHandler(func(conn transports.Connection) {
		nc.connections = append(nc.connections, conn)
		if onConnect != nil {
			onConnect(conn)
		}
	})
	nc.transport.WaitForConnection()
}

func (nc NetworkController) Broadcast(data []byte) {
	for _, conn := range nc.connections {
		conn.Send(data)
	}
}
