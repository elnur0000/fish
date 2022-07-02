package components

import "time"

type Player struct {
	ID              uint8
	Object          *Object
	LastStateChange time.Time
}

func NewPlayer(id uint8, object *Object) Player {
	return Player{
		ID:              id,
		Object:          object,
		LastStateChange: time.Now(),
	}
}

func (p *Player) HandleMovement(rotation float32, duration float32) {
	p.LastStateChange = time.Now()
	p.Object.SetRotation(rotation)
	p.Object.Move(duration)
}
