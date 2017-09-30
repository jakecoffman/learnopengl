package breakout

import "github.com/go-gl/mathgl/mgl32"

type Object struct {
	Position, Size, Velocity mgl32.Vec2
	Color mgl32.Vec3
	Rotation float64

	IsSolid, Destroyed bool

	Sprite *Texture2D
}

var (
	DefaultGameObjectColor = mgl32.Vec3{1, 1, 1}
)

func NewGameObject(sprite *Texture2D) *Object {
	return &Object{
		Size: mgl32.Vec2{1, 1},
		Color: mgl32.Vec3{1, 1, 1},
		Sprite: sprite,
	}
}

func NewGameObject2(pos, size mgl32.Vec2, sprite *Texture2D, color mgl32.Vec3, velocity mgl32.Vec2) *Object {
	return &Object{
		pos,
		size,
		velocity,
		color,
		0,
		false,
		false,
		sprite,
	}
}

func (g *Object) Draw(renderer *SpriteRenderer) {
	renderer.DrawSprite(g.Sprite, g.Position, g.Size, g.Rotation, g.Color)
}
