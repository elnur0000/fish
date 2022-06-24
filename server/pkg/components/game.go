package components

import (
	"errors"
	"log"
	"sync"
)

const MAX_PLAYER_PER_GAME = 64

type Game struct {
	Players [MAX_PLAYER_PER_GAME]*Player
	m       *sync.Mutex
	World   World
}

func NewGame() Game {
	return Game{
		m:     &sync.Mutex{},
		World: NewWorld(0, 5000, 0, 5000),
	}
}

func (g *Game) SpawnNewPlayer(conn Client) (*Player, error) {
	g.m.Lock()
	defer g.m.Unlock()

	for i, player := range g.Players {
		var id uint8 = uint8(i + 1)
		if player == nil {
			p := NewPlayer(id, g.World.CreateObject(Vec{X: 50, Y: 50}, 70, 70), conn)
			g.Players[i] = &p

			log.Printf("Created a new player with id %d", p.ID)

			return &p, nil
		}
	}

	return nil, errors.New("Game is full")
}

func (g *Game) RemovePlayer(ID uint8) {
	log.Printf("Removing a player with id %d", ID)
	g.m.Lock()
	defer g.m.Unlock()
	g.Players[ID-1] = nil

}

func (g *Game) UpdatePlayer(p Player) {
	g.m.Lock()
	defer g.m.Unlock()
	g.Players[p.ID-1] = &p
}