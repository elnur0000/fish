package components

import (
	"math"
)

const DefaultVelocity float32 = 0.5 // per milliseconds
const DefaultCellSize uint16 = 50

type World struct {
	Bounds   Rect
	Width    uint32
	Height   uint32
	cells    [][]*Object
	cellSize uint16
}

type Rect struct {
	A Vec
	B Vec
	C Vec
	D Vec
}

func (r *Rect) Subtract(v Vec) {
	r.A.Subtract(v)
	r.B.Subtract(v)
	r.C.Subtract(v)
	r.D.Subtract(v)
}

func (r Rect) Center() Vec {
	height := r.B.Y - r.A.Y
	width := r.C.X - r.B.X

	return Vec{
		X: r.A.X + width/2,
		Y: r.A.Y + height/2,
	}
}

type Vec struct {
	X float32
	Y float32
}

func (v *Vec) Subtract(targetVec Vec) {
	v.X -= targetVec.X
	v.Y -= targetVec.Y
}

func SumVecs(vec1, vec2 Vec) *Vec {
	newVec := Vec{
		X: vec1.X + vec2.X,
		Y: vec1.Y + vec2.Y,
	}
	return &newVec
}

func SubtractVecs(vec1, vec2 Vec) *Vec {
	newVec := Vec{
		X: vec1.X - vec2.X,
		Y: vec1.Y - vec2.Y,
	}
	return &newVec
}

type Object struct {
	Rect     Rect
	Position Vec
	world    *World
	Width    uint32
	Height   uint32
	Velocity float32
	Rotation float32
}

func NewWorld(x, x1, y, y1 float32) World {
	world := World{
		Bounds: Rect{
			A: Vec{
				X: x,
				Y: y,
			},
			B: Vec{
				X: x,
				Y: y1,
			},
			C: Vec{
				X: x1,
				Y: y1,
			},
			D: Vec{
				X: x1,
				Y: y,
			},
		},
		Width:    uint32(x1 - x),
		Height:   uint32(y1 - y),
		cellSize: DefaultCellSize,
	}

	world.buildCells()

	return world
}

func (w *World) buildCells() {
	w.cells = make([][]*Object, w.Width/uint32(w.cellSize)*w.Height/uint32(w.cellSize))
}

func (w *World) CreateObject(pos Vec, width, height uint32) *Object {
	object := Object{
		Position: pos,
		world:    w,
		Rect: Rect{
			A: Vec{
				X: pos.X - float32(width/2),
				Y: pos.Y - float32(height/2),
			},
			B: Vec{
				X: pos.X - float32(width/2),
				Y: pos.Y + float32(height/2),
			},
			C: Vec{
				X: pos.X + float32(width/2),
				Y: pos.Y + float32(height/2),
			},
			D: Vec{
				X: pos.X + float32(width/2),
				Y: pos.Y - float32(height/2),
			},
		},
		Width:    width,
		Height:   height,
		Velocity: DefaultVelocity,
	}

	w.insertToCells(&object)

	return &object
}

func (w *World) insertToCells(object *Object) {
	for yi := int(object.Rect.A.Y); yi < int(object.Rect.C.Y); yi += int(w.cellSize) {
		for xi := int(object.Rect.A.X); xi < int(object.Rect.C.X); xi += int(w.cellSize) {
			cellIndex := w.cellIndex(Vec{X: float32(xi), Y: float32(yi)})
			w.cells[cellIndex] = append(w.cells[cellIndex], object)
		}
	}
}

func (w *World) removeFromCells(object *Object) {
	for yi := int(object.Rect.A.Y); yi < int(object.Rect.C.Y); yi += int(w.cellSize) {
		for xi := int(object.Rect.A.X); xi < int(object.Rect.C.X); xi += int(w.cellSize) {
			cellIndex := w.cellIndex(Vec{X: float32(xi), Y: float32(yi)})
			cell := w.cells[cellIndex]
			for i, obj := range cell {
				if obj == object {
					if len(w.cells) < 2 {
						cell = []*Object{}
						continue
					}
					// swap current object with the last one and then remove it, for efficiency
					lastObj := cell[len(cell)-1]
					cell[i] = lastObj
					cell = cell[:len(cell)-1]
				}
			}
			w.cells[cellIndex] = cell
		}
	}

}

func (w *World) cellIndex(pos Vec) int {
	return int(pos.X)/int(w.cellSize) + (int(w.Width)/int(w.cellSize))*(int(pos.Y)/int(w.cellSize))
}

func (o *Object) Move(duration float32) {
	o.world.removeFromCells(o)
	dx := math.Cos(float64(o.Rotation)) * float64(o.Velocity*duration)
	dy := math.Sin(float64(o.Rotation)) * float64(o.Velocity*duration)
	dVec := Vec{
		X: -float32(dx),
		Y: -float32(dy),
	}

	newRect := Rect{
		A: *SumVecs(o.Rect.A, dVec),
		B: *SumVecs(o.Rect.B, dVec),
		C: *SumVecs(o.Rect.C, dVec),
		D: *SumVecs(o.Rect.D, dVec),
	}

	topBound := o.world.Bounds.B.Y
	leftBound := o.world.Bounds.B.X
	bottomBound := o.world.Bounds.A.Y
	rightBound := o.world.Bounds.C.X

	if newRect.B.Y > topBound {
		distance := newRect.B.Y - topBound
		newRect.Subtract(Vec{X: 0, Y: distance})
	}

	if newRect.B.X < leftBound {
		distance := newRect.B.X - leftBound
		newRect.Subtract(Vec{X: distance, Y: 0})
	}

	if newRect.A.Y < bottomBound {
		distance := newRect.A.Y - bottomBound
		newRect.Subtract(Vec{X: 0, Y: distance})
	}

	if newRect.C.X > rightBound {
		distance := newRect.C.X - rightBound
		newRect.Subtract(Vec{X: distance, Y: 0})
	}

	o.Rect = newRect
	o.Position = newRect.Center()

	o.world.insertToCells(o)
}

func (o *Object) SetVelocity(v float32) {
	o.Velocity = v
}

func (o *Object) SetRotation(r float32) {
	o.Rotation = r
}
