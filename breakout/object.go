package breakout

import (
	"fmt"
	"math"
	"github.com/jakecoffman/learnopengl/breakout/eng"
	"github.com/go-gl/mathgl/mgl32"
)

type Object struct {
	Position, Size, Velocity mgl32.Vec2
	Color                    mgl32.Vec3
	Rotation                 float64

	IsSolid, Destroyed bool

	Sprite *eng.Texture2D
}

func (o Object) String() string {
	return fmt.Sprintf("Object(@ %v - Color: %v)", o.Position, o.Color)
}

var (
	DefaultGameObjectColor = mgl32.Vec3{1, 1, 1}
)

func NewGameObject(pos, size mgl32.Vec2, sprite *eng.Texture2D) *Object {
	return &Object{pos, size, mgl32.Vec2{0, 0}, DefaultGameObjectColor, 0, false, false, sprite,}
}

func (o *Object) Draw(renderer *eng.SpriteRenderer, last *mgl32.Vec2, alpha float32) {
	pos := o.Position
	if last != nil {
		pos = SLerp(*last, o.Position, alpha)
	}
	renderer.DrawSprite(o.Sprite, pos, o.Size, o.Rotation, o.Color)
}

func SLerp(v, other mgl32.Vec2, t float32) mgl32.Vec2 {
	dot := v.Normalize().Dot(other.Normalize())
	omega := float32(math.Acos(Clamp(float64(dot), -1, 1)))

	if omega < 1e-3 {
		return Lerp(v, other, t)
	}

	denom := 1.0 / math.Sin(float64(omega))
	return v.Mul(float32(math.Sin(float64((1.0-t)*omega)) * denom)).Add(other.Mul(float32(math.Sin(float64(t*omega)) * denom)))
}

func Lerp(v, other mgl32.Vec2, t float32) mgl32.Vec2 {
	return v.Mul(1.0 - t).Add(other.Mul(t))
}

func Clamp(f, min, max float64) float64 {
	if f > min {
		return math.Min(f, max)
	} else {
		return math.Min(min, max)
	}
}
