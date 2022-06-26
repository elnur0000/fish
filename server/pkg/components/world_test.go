package components

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestObjectCreate(t *testing.T) {
	world := NewWorld(0, 500, 0, 100)
	object := world.CreateObject(Vec{X: 50, Y: 50}, 100, 50)
	assert.Equal(t, world.cells[0][0], object)
}

func TestObjectRemove(t *testing.T) {
	world := NewWorld(0, 500, 0, 100)
	object := world.CreateObject(Vec{X: 50, Y: 50}, 100, 50)
	world.removeFromCells(object)
	assert.Equal(t, len(world.cells[0]), 0)
	assert.Equal(t, len(world.cells[1]), 0)
}

func BenchmarkInsert(b *testing.B) {
	world := NewWorld(0, 5000, 0, 5000)
	for i := 0; i < 50000; i++ {
		obj := world.CreateObject(Vec{X: 70, Y: 70}, 50, 50)
		world.removeFromCells(obj)
	}
	fmt.Println(world.cells)
}
