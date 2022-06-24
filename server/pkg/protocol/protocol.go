package protocol

import (
	"bytes"
	"encoding/binary"
)

type MessageType uint8

const (
	CONTROL      MessageType = iota // player input
	REGISTERED                      // when player is registered
	PLAYER_STATE                    // player state
)

type ControlMessage struct {
	Type     MessageType
	Rotation float32
	Duration float32
}

func NewControlMessage(rotation float32, duration float32) ControlMessage {
	return ControlMessage{
		Type:     CONTROL,
		Rotation: rotation,
		Duration: duration,
	}
}

type RegisteredMessage struct {
	Type MessageType
	ID   uint8
}

func NewRegisteredMessage(id uint8) RegisteredMessage {
	return RegisteredMessage{
		Type: REGISTERED,
		ID:   id,
	}
}

type PlayerStateMessage struct {
	Type     MessageType
	ID       uint8
	X        float32
	Y        float32
	Rotation float32
	Velocity float32
	Height   float32
	Width    float32
}

func NewPlayerStateMessage(id uint8, x, y, rotation, height, width, velocity float32) PlayerStateMessage {
	return PlayerStateMessage{
		Type:     PLAYER_STATE,
		ID:       id,
		X:        x,
		Y:        y,
		Rotation: rotation,
		Velocity: velocity,
		Height:   height,
		Width:    width,
	}
}

func Serialize(msgBody any) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, msgBody)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func Deserialize(binaryMsg []byte, result any) error {
	err := binary.Read(bytes.NewReader(binaryMsg), binary.BigEndian, result)

	if err != nil {
		return err
	}

	return nil
}

func GetMessageType(binaryMsg []byte) MessageType {
	return MessageType(binaryMsg[0])
}
