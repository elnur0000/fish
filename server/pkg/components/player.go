package components

type Client interface {
	Send([]byte)
	Recv() chan []byte
}

type Player struct {
	ID     uint8
	Object Object
	Client Client
}

func NewPlayer(id uint8, object Object, client Client) Player {
	return Player{
		ID:     id,
		Object: object,
		Client: client,
	}
}

func (p *Player) HandleMovement(rotation float32, duration float32) {
	p.Object.SetRotation(rotation)
	p.Object.Move(duration)
}
