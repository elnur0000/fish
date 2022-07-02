package controller

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/elnur0000/fish-backend/pkg/components"
	"github.com/elnur0000/fish-backend/pkg/protocol"
	"github.com/elnur0000/fish-backend/pkg/transports"
)

const GameLoopInterval = 5 * time.Millisecond
const MaxPlayerPerGame = 64

type GameController struct {
	networkController NetworkController
	game              components.Game
	capacity          int
	playerCount       int
}

func newGameController() GameController {
	return GameController{
		networkController: newNetworkController(),
		game:              components.NewGame(MaxPlayerPerGame),
		capacity:          MaxPlayerPerGame,
	}
}

func (g *GameController) addPlayer(conn transports.Connection) (*components.Player, error) {

	p, err := g.game.SpawnNewPlayer()

	g.playerCount++

	if err != nil {
		return nil, err
	}

	g.networkController.addConnection(p.ID, conn)

	msg, err := protocol.Serialize(protocol.NewRegisteredMessage(p.ID))

	if err != nil {
		return nil, err
	}

	conn.Send(msg)

	return p, nil
}

func (g *GameController) removePlayer(playerId uint8) {
	g.game.RemovePlayer(playerId)
	g.playerCount--
	g.networkController.remove(playerId)

	msg, err := protocol.Serialize(protocol.NewPlayerDisconnectedMessage(playerId))
	if err != nil {
		log.Printf("Failed to serialize player disconnected message: %v", err)
		return
	}
	g.networkController.broadcast(msg)
}

func (g *GameController) onNewConnection(conn transports.Connection) {
	player, err := g.addPlayer(conn)

	if err != nil {
		log.Println(err)
		return
	}

	conn.OnDisconnect(func() {
		g.removePlayer(player.ID)
	})
}

func (g *GameController) startGameLoop() {
	for range time.Tick(GameLoopInterval) {
		if g.playerCount == 0 {
			continue
		}

		g.updateGameState()
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
			g.networkController.broadcast(msg)
		}
	}
}

type PlayerMessage struct {
	Message transports.Message
	Player  *components.Player
}

func (g GameController) updateGameState() {
	playerMessages := []PlayerMessage{}

	for _, player := range g.game.Players {
		if player != nil {
			g.extractPlayerMessages(player, &playerMessages)
		}
	}

	sort.SliceStable(playerMessages, func(i, j int) bool {
		return playerMessages[i].Message.ReceivedAt.Before(playerMessages[j].Message.ReceivedAt)
	})

	g.processPlayerMessages(playerMessages)
}

func (g GameController) processPlayerMessages(playerMessages []PlayerMessage) {
	for _, msg := range playerMessages {
		msgType := uint8(msg.Message.Body[0])
		switch msgType {
		case uint8(protocol.CONTROL):
			var controlMsg protocol.ControlMessage
			err := protocol.Deserialize(msg.Message.Body, &controlMsg)
			if err != nil {
				fmt.Println("Failed to decode player control msg")
				fmt.Println(err)
			}

			elapsed := time.Since(msg.Player.LastStateChange)
			if controlMsg.Duration > float32(elapsed*time.Millisecond) {
				g.removePlayer(msg.Player.ID)
				log.Printf("Removed player %d due to possible cheating attempt", msg.Player.ID)
				continue
			}
			msg.Player.HandleMovement(controlMsg.Rotation, controlMsg.Duration)
		}
	}
}

func (g GameController) extractPlayerMessages(player *components.Player, playerMessages *[]PlayerMessage) {
	recvBuffer := g.networkController.getPlayerConnection(player.ID).Recv()
	for len(recvBuffer) != 0 {
		message := <-recvBuffer
		*playerMessages = append(*playerMessages, PlayerMessage{
			Message: message,
			Player:  player,
		})
	}
}

type GamePool struct {
	gameControllers []*GameController
	transport       transports.Transport
}

type GamePoolOptions struct {
	Transport transports.Transport
}

func NewGamePool(options GamePoolOptions) GamePool {
	return GamePool{
		transport: options.Transport,
	}
}

func (gp *GamePool) Start() {
	gp.transport.RegisterConnHandler(func(conn transports.Connection) {
		gameController := gp.findJoinableGame()
		gameController.onNewConnection(conn)
	})
	gp.transport.WaitForConnection()
}

func (gp *GamePool) findJoinableGame() *GameController {
	for _, gameController := range gp.gameControllers {
		if gameController.capacity != gameController.playerCount {
			return gameController
		}
	}

	gameController := newGameController()
	gp.gameControllers = append(gp.gameControllers, &gameController)
	go gameController.startGameLoop()
	return &gameController
}
