package components

import (
	"errors"
	"log"
	"sync"
)

type Game struct {
	Players []*Player
	m       *sync.Mutex
	World   World
}

func NewGame(maxPlayers int) Game {
	return Game{
		m:       &sync.Mutex{},
		World:   NewWorld(0, 5000, 0, 5000),
		Players: make([]*Player, maxPlayers),
	}
}

func (g *Game) SpawnNewPlayer() (*Player, error) {
	g.m.Lock()
	defer g.m.Unlock()

	for i, player := range g.Players {
		var id uint8 = uint8(i + 1)
		if player == nil {
			p := NewPlayer(id, g.World.CreateObject(Vec{X: 50, Y: 50}, 70, 70))
			g.Players[i] = &p

			log.Printf("Created a new player with id %d", p.ID)

			return &p, nil
		}
	}

	return nil, errors.New("game is full")
}

func (g *Game) RemovePlayer(ID uint8) {
	player := g.Players[ID-1]
	player.Object.world.remove(player.Object)
	g.Players[ID-1] = nil
}

func (g *Game) UpdatePlayer(p Player) {
	g.m.Lock()
	defer g.m.Unlock()
	if g.Players[p.ID-1] != nil {
		g.RemovePlayer(p.ID)
	}

	g.Players[p.ID-1] = &p
}
