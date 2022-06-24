package main

import (
	"syscall/js"

	"github.com/elnur0000/fish-backend/pkg/components"
	"github.com/elnur0000/fish-backend/pkg/protocol"
)

var game = components.NewGame()
var localPlayerId uint8

func processServerMessage(this js.Value, args []js.Value) interface{} {
	msg := arrayBufferToBytes(args[0])

	switch protocol.GetMessageType(msg) {

	case protocol.PLAYER_STATE:
		var playerState protocol.PlayerStateMessage
		err := protocol.Deserialize(msg, &playerState)

		if err != nil {
			println(err.Error())
			return nil
		}

		p := components.Player{
			ID: playerState.ID,
			Object: game.World.CreateObject(components.Vec{
				X: playerState.X,
				Y: playerState.Y,
			},
				playerState.Width,
				playerState.Height,
			),
		}
		p.Object.SetRotation(playerState.Rotation)

		game.UpdatePlayer(p)
	case protocol.REGISTERED:
		println("Registered")
		var registeredMessage protocol.RegisteredMessage
		err := protocol.Deserialize(msg, &registeredMessage)

		if err != nil {
			println(err.Error())
			return nil
		}
		localPlayerId = registeredMessage.ID
	}

	return nil
}

func getPlayerStates(this js.Value, args []js.Value) interface{} {
	var result []interface{}
	for _, p := range game.Players {
		if p != nil {
			result = append(result, map[string]interface{}{
				"ID":       p.ID,
				"rotation": p.Object.Rotation,
				"x":        p.Object.Position.X,
				"y":        p.Object.Position.Y,
				"velocity": p.Object.Velocity,
				"height":   p.Object.Height,
				"width":    p.Object.Width,
			})
		}
	}

	return js.ValueOf(result)
}

func processLocalPlayerInput(this js.Value, args []js.Value) interface{} {
	rotation := args[0].Float()
	duration := args[1].Float()
	sendPlayerInput := args[2]

	if localPlayerId == 0 {
		return nil
	}
	localPlayer := game.Players[localPlayerId-1]

	localPlayer.HandleMovement(float32(rotation), float32(duration))

	controlMsg := protocol.NewControlMessage(localPlayer.Object.Rotation, float32(duration))
	serializedControlMsg, err := protocol.Serialize(controlMsg)

	if err != nil {
		println(err.Error())
		return nil
	}

	uint8Array := js.Global().Get("Uint8Array")
	dst := uint8Array.New(len(serializedControlMsg))
	js.CopyBytesToJS(dst, serializedControlMsg)

	sendPlayerInput.Invoke(dst)

	return nil
}

func arrayBufferToBytes(arrayBuffer js.Value) []byte {
	result := make([]uint8, arrayBuffer.Get("byteLength").Int())
	js.CopyBytesToGo(result, arrayBuffer)
	return result
}

func getLocalPlayerId(this js.Value, args []js.Value) interface{} {
	if localPlayerId == 0 {
		return nil
	}

	return js.ValueOf(localPlayerId)
}

func main() {
	js.Global().Set("processLocalPlayerInput", js.FuncOf(processLocalPlayerInput))
	js.Global().Set("processServerMessage", js.FuncOf(processServerMessage))
	js.Global().Set("getPlayerStates", js.FuncOf(getPlayerStates))
	js.Global().Set("getLocalPlayerId", js.FuncOf(getLocalPlayerId))
	select {}
}
