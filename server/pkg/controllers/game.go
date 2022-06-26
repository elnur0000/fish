package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/elnur0000/fish-backend/pkg/components"
	"github.com/elnur0000/fish-backend/pkg/protocol"
	"github.com/elnur0000/fish-backend/pkg/transports"
)

const GameLoopInterval = 5 * time.Millisecond

type GameController struct {
	networkController NetworkController
	game              components.Game
}

type GameOptions struct {
	Transport transports.Transport
}

func NewGameController(options GameOptions) GameController {
	return GameController{
		networkController: NewNetworkController(options.Transport),
		game:              components.NewGame(),
	}
}

func (g *GameController) AddPlayer(playerConn transports.Connection) (*components.Player, error) {

	p, err := g.game.SpawnNewPlayer(playerConn)

	if err != nil {
		return nil, err
	}

	msg, err := protocol.Serialize(protocol.NewRegisteredMessage(p.ID))
	if err != nil {
		return nil, err
	}
	playerConn.Send(msg)

	return p, nil
}

func (g GameController) StartGameLoop() {
	go g.networkController.ListenForConnections(func(conn transports.Connection) {
		player, err := g.AddPlayer(conn)

		if err != nil {
			log.Println(err)
		}

		conn.OnDisconnect(func() {
			g.game.RemovePlayer(player.ID)
			// remove conn
		})
	})

	for range time.Tick(GameLoopInterval) {
		g.processPlayerCommands()
		g.broadcastGameState()
	}
}

func (g GameController) broadcastGameState() {
	for _, player := range g.game.Players {
		if player != nil {
			msg, err := protocol.Serialize(protocol.NewPlayerStateMessage(player.ID, player.Object.Position.X, player.Object.Position.Y, player.Object.Rotation, player.Object.Velocity, player.Object.Height, player.Object.Width))
			if err != nil {
				continue
			}
			g.networkController.Broadcast(msg)
		}
	}
}

func (g GameController) processPlayerCommands() {
	for _, player := range g.game.Players {
		if player != nil {
			recvBuffer := player.Client.Recv()
			for len(recvBuffer) != 0 {
				message := <-recvBuffer

				msgType := uint8(message[0])

				switch msgType {
				case uint8(protocol.CONTROL):
					var controlMsg protocol.ControlMessage
					err := protocol.Deserialize(message, &controlMsg)
					if err != nil {
						fmt.Println("Failed to decode player control msg")
						fmt.Println(err)
					}
					player.HandleMovement(controlMsg.Rotation, controlMsg.Duration)
				}

			}
		}
	}
}
