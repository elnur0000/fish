package components

import "math"

const DEFAULT_VELOCITY float32 = 0.5 // per milliseconds

type World struct {
	Bounds  Rect
	Objects []Object
}

type Rect struct {
	A Vec
	B Vec
	C Vec
	D Vec
}

func (r *Rect) Subtract(v Rect) {
	r.A.Subtract(v.A)
	r.B.Subtract(v.B)
	r.C.Subtract(v.C)
	r.D.Subtract(v.D)
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
	Width    float32
	Height   float32
	Velocity float32
	Rotation float32
}

func NewWorld(x, x1, y, y1 float32) World {
	return World{
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
	}
}

func (w *World) CreateObject(pos Vec, width, height float32) Object {
	object := Object{
		Position: pos,
		world:    w,
		Rect: Rect{
			A: Vec{
				X: pos.X - width/2,
				Y: pos.Y - height/2,
			},
			B: Vec{
				X: pos.X - width/2,
				Y: pos.Y + height/2,
			},
			C: Vec{
				X: pos.X + width/2,
				Y: pos.Y + height/2,
			},
			D: Vec{
				X: pos.X + width/2,
				Y: pos.Y - height/2,
			},
		},
		Width:    width,
		Height:   height,
		Velocity: DEFAULT_VELOCITY,
	}

	w.Objects = append(w.Objects, object)

	return object
}

func (o *Object) Move(duration float32) {
	dx := math.Cos(float64(o.Rotation)) * float64((o.Velocity * duration))
	dy := math.Sin(float64(o.Rotation)) * float64((o.Velocity * duration))

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
		newRect.Subtract(Rect{
			A: Vec{X: 0, Y: distance},
			B: Vec{X: 0, Y: distance},
			C: Vec{X: 0, Y: distance},
			D: Vec{X: 0, Y: distance},
		})
	}

	if newRect.B.X < leftBound {
		distance := newRect.B.X - leftBound
		newRect.Subtract(Rect{
			A: Vec{X: distance, Y: 0},
			B: Vec{X: distance, Y: 0},
			C: Vec{X: distance, Y: 0},
			D: Vec{X: distance, Y: 0},
		})
	}

	if newRect.A.Y < bottomBound {
		distance := newRect.A.Y - bottomBound
		newRect.Subtract(Rect{
			A: Vec{X: 0, Y: distance},
			B: Vec{X: 0, Y: distance},
			C: Vec{X: 0, Y: distance},
			D: Vec{X: 0, Y: distance},
		})
	}

	if newRect.C.X > rightBound {
		distance := newRect.C.X - rightBound
		newRect.Subtract(Rect{
			A: Vec{X: distance, Y: 0},
			B: Vec{X: distance, Y: 0},
			C: Vec{X: distance, Y: 0},
			D: Vec{X: distance, Y: 0},
		})
	}

	o.Rect = newRect
	o.Position = newRect.Center()
}

func (o *Object) SetVelocity(v float32) {
	o.Velocity = v
}

func (o *Object) SetRotation(r float32) {
	o.Rotation = r
}
