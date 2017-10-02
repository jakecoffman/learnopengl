package breakout

import "github.com/go-gl/mathgl/mgl32"

type Ball struct {
	*Object

	Radius float32
	Stuck bool
}

func NewBall(pos mgl32.Vec2, radius float32, velocity mgl32.Vec2, sprite *Texture2D) *Ball {
	ball := &Ball{}
	ball.Object = NewGameObject(pos, mgl32.Vec2{float32(radius*2), float32(radius*2)}, sprite)
	ball.Color = mgl32.Vec3{1, 1, 1}
	ball.Velocity = velocity
	ball.Radius = radius
	ball.Stuck = true
	return ball
}

func (b *Ball) Move(dt float32, windowWidth float32) mgl32.Vec2 {
	if b.Stuck {
		return b.Position
	}

	b.Position = b.Position.Add(b.Velocity.Mul(dt))
	if b.Position.X() <= 0 {
		b.Velocity = mgl32.Vec2{-b.Velocity.X(), b.Velocity.Y()}
		b.Position = mgl32.Vec2{0, b.Position.Y()}
	} else if b.Position.X() + b.Size.X() >= windowWidth {
		b.Velocity = mgl32.Vec2{-b.Velocity.X(), b.Velocity.Y()}
		b.Position = mgl32.Vec2{windowWidth-b.Size.X(), b.Position.Y()}
	}
	if b.Position.Y() <= 0 {
		b.Velocity = mgl32.Vec2{b.Velocity.X(), -b.Velocity.Y()}
		b.Position = mgl32.Vec2{b.Position.X(), 0}
	}

	return b.Position
}

func (b *Ball) Reset(position, velocity mgl32.Vec2) {
	b.Position = position
	b.Velocity = velocity
	b.Stuck = true
}