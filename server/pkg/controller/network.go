package controller

import (
	"github.com/elnur0000/fish-backend/pkg/transports"
)

type NetworkController struct {
	connections map[uint8]transports.Connection
}

func newNetworkController() NetworkController {
	return NetworkController{
		connections: map[uint8]transports.Connection{},
	}
}

func (nc *NetworkController) addConnection(playerId uint8, conn transports.Connection) {
	nc.connections[playerId] = conn
}

func (nc NetworkController) broadcast(data []byte) {
	for _, conn := range nc.connections {
		conn.Send(data)
	}
}

func (nc NetworkController) getPlayerConnection(playerId uint8) transports.Connection {
	return nc.connections[playerId]
}

func (nc *NetworkController) remove(playerId uint8) {
	delete(nc.connections, playerId)
}
