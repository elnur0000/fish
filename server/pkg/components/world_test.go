package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectRemove(t *testing.T) {
	width := 50000
	height := 50000
	world := NewWorld(0, float32(width), 0, float32(height))
	objects := []*Object{}
	for i := 150; i < width-150; i++ {
		obj := world.CreateObject(Vec{X: float32(i), Y: float32(i)}, 150, 150)
		objects = append(objects, obj)
	}

	for _, obj := range objects {
		obj.world.remove(obj)
	}

	var objectCount int
	for _, cell := range world.cells {
		objectCount += cell.Length
	}

	assert.Equal(t, 0, objectCount)
}

func BenchmarkInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		width := 50000
		height := 50000
		world := NewWorld(0, float32(width), 0, float32(height))
		for i := 150; i < width-150; i++ {
			obj := world.CreateObject(Vec{X: float32(i), Y: float32(i)}, 150, 150)
			world.remove(obj)
		}
	}
}
