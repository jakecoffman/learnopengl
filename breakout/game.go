package breakout

import (
	"github.com/jakecoffman/learnopengl/breakout/eng"
	"math"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var ResourceManager = eng.NewResourceManager()

type Game struct {
	state         int
	Keys          [1024]bool
	Width, Height int

	Levels []*Level
	Level  int

	Player *Object
	Ball   *Ball

	ParticleGenerator *eng.ParticleGenerator
	SpriteRenderer    *eng.SpriteRenderer
	TextRenderer      *eng.TextRenderer
}

// Game state
const (
	GAME_ACTIVE = iota
	GAME_MENU
	GAME_WIN
)

func NewGame(width, height int) *Game {
	return &Game{
		Width:  width,
		Height: height,
		Keys:   [1024]bool{},
	}
}

var (
	PLAYER_SIZE                   = mgl32.Vec2{100, 20}
	PLAYER_VELOCITY               = 500.0
	INITIAL_BALL_VELOCITY         = Vec2(100, -350)
	BALL_RADIUS           float32 = 12.5
)

func (g *Game) Init(window *glfw.Window) {
	window.SetKeyCallback(func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyEscape && action == glfw.Press {
			window.SetShouldClose(true)
		}
		if key >= 0 && key < 1024 {
			if action == glfw.Press {
				g.Keys[key] = true
			} else if action == glfw.Release {
				g.Keys[key] = false
			}
		}
	})

	width, height := float32(g.Width), float32(g.Height)

	ResourceManager.LoadShader("breakout/shaders/main.vs.glsl", "breakout/shaders/main.fs.glsl", "sprite")
	ResourceManager.LoadShader("breakout/shaders/particle.vs.glsl", "breakout/shaders/particle.fs.glsl", "particle")

	projection := mgl32.Ortho(0, width, height, 0, -1, 1)
	ResourceManager.Shader("sprite").
		Use().
		SetInt("sprite", 0).
		SetMat4("projection", projection)
	ResourceManager.Shader("particle").
		Use().
		SetInt("sprite", 0).
		SetMat4("projection", projection)

	ResourceManager.LoadTexture("breakout/textures/background.jpg", "background")
	ResourceManager.LoadTexture("breakout/textures/paddle.png", "paddle")
	ResourceManager.LoadTexture("breakout/textures/particle.png", "particle")
	ResourceManager.LoadTexture("breakout/textures/awesomeface.png", "face")
	ResourceManager.LoadTexture("breakout/textures/block.png", "block")
	ResourceManager.LoadTexture("breakout/textures/block_solid.png", "block_solid")

	shader, _ := ResourceManager.LoadShader("breakout/shaders/text.vs.glsl", "breakout/shaders/text.fs.glsl", "text")
	g.TextRenderer = eng.NewTextRenderer(shader, width, height, "breakout/textures/Roboto-Light.ttf", 24)
	g.TextRenderer.SetColor(1, 1, 1, 1)

	g.ParticleGenerator = eng.NewParticleGenerator(ResourceManager.Shader("particle"), ResourceManager.Texture("particle"), 500)
	g.SpriteRenderer = eng.NewSpriteRenderer(ResourceManager.Shader("sprite"))

	one := NewLevel()
	if err := one.Load("breakout/levels/1.txt", g.Width, int(float32(g.Height)*0.5)); err != nil {
		panic(err)
	}
	g.Levels = append(g.Levels, one)

	playerPos := mgl32.Vec2{float32(g.Width)/2.0 - PLAYER_SIZE.X()/2.0, float32(g.Height) - PLAYER_SIZE.Y()}
	g.Player = NewGameObject(playerPos, PLAYER_SIZE, ResourceManager.Texture("paddle"))

	ballPos := playerPos.Add(mgl32.Vec2{PLAYER_SIZE.X()/2.0 - float32(BALL_RADIUS), float32(-BALL_RADIUS * 2)})
	g.Ball = NewBall(ballPos, BALL_RADIUS, INITIAL_BALL_VELOCITY, ResourceManager.Texture("face"))

	g.state = GAME_ACTIVE
}

func (g *Game) ProcessInput(dt float64) {
	if g.state != GAME_ACTIVE {
		return
	}

	velocity := float32(PLAYER_VELOCITY * dt)

	if g.Keys[glfw.KeyA] {
		if g.Player.Position.X() >= 0 {
			g.Player.Position = mgl32.Vec2{g.Player.Position.X() - velocity, g.Player.Position.Y()}
			if g.Ball.Stuck {
				g.Ball.Position = mgl32.Vec2{g.Ball.Position.X() - velocity, g.Ball.Position.Y()}
			}
		}
	}
	if g.Keys[glfw.KeyD] {
		if g.Player.Position.X() <= float32(g.Width)-g.Player.Size.X() {
			g.Player.Position = mgl32.Vec2{g.Player.Position.X() + velocity, g.Player.Position.Y()}
			if g.Ball.Stuck {
				g.Ball.Position = mgl32.Vec2{g.Ball.Position.X() + velocity, g.Ball.Position.Y()}
			}
		}
	}
	if g.Keys[glfw.KeySpace] {
		g.Ball.Stuck = false
	}
}

func (g *Game) Update(dt float32) {
	g.Ball.Move(dt, float32(g.Width))
	g.DoCollisions()
	ball := g.Ball.Object
	g.ParticleGenerator.Update(dt, ball.Position, ball.Velocity, 2, mgl32.Vec2{g.Ball.Radius / 2, g.Ball.Radius / 2})
	if g.Ball.Position.Y() >= float32(g.Height) {
		g.ResetLevel()
		g.ResetPlayer()
	}
}

func (g *Game) Render() {
	if g.state == GAME_ACTIVE {
		g.SpriteRenderer.DrawSprite(ResourceManager.Texture("background"), Vec2(0, 0), Vec2(g.Width, g.Height), 0, eng.DefaultColor)
		g.Levels[g.Level].Draw(g.SpriteRenderer)
		g.Player.Draw(g.SpriteRenderer)
		g.ParticleGenerator.Draw()
		g.Ball.Draw(g.SpriteRenderer)
	}
	g.TextRenderer.Print("Hello, world!", 10, 25, 1)
}

func (g *Game) Pause() {
	g.state = GAME_MENU
}

func (g *Game) Unpause() {
	g.state = GAME_ACTIVE
}

func (g *Game) ResetLevel() {
	if g.Level == 0 {
		g.Levels[0].Load("breakout/level1.txt", g.Width, int(float32(g.Height)*0.5))
	}
	// TODO
}

func (g *Game) ResetPlayer() {
	g.Player.Size = PLAYER_SIZE
	g.Player.Position = mgl32.Vec2{float32(g.Width)/2 - PLAYER_SIZE.X()/2, float32(g.Height) - PLAYER_SIZE.Y()}
	g.Ball.Reset(g.Player.Position.Add(mgl32.Vec2{PLAYER_SIZE.X()/2 - BALL_RADIUS, -(BALL_RADIUS * 2)}), INITIAL_BALL_VELOCITY)
}

func (g *Game) DoCollisions() {
	for _, box := range g.Levels[g.Level].Bricks {
		if !box.Destroyed {
			collides, direction, diff := CheckBallCollision(g.Ball, box)
			if collides {
				if !box.IsSolid {
					box.Destroyed = true
				}

				if direction == DIRECTION_LEFT || direction == DIRECTION_RIGHT {
					g.Ball.Velocity = mgl32.Vec2{-g.Ball.Velocity.X(), g.Ball.Velocity.Y()}
					penetration := g.Ball.Radius - float32(math.Abs(float64(diff.X())))
					if direction == DIRECTION_LEFT {
						g.Ball.Position = mgl32.Vec2{g.Ball.Position.X() + penetration, g.Ball.Position.Y()}
					} else {
						g.Ball.Position = mgl32.Vec2{g.Ball.Position.X() - penetration, g.Ball.Position.Y()}
					}
				} else {
					g.Ball.Velocity = mgl32.Vec2{g.Ball.Velocity.X(), -g.Ball.Velocity.Y()}
					penetration := g.Ball.Radius - float32(math.Abs(float64(diff.Y())))
					if direction == DIRECTION_UP {
						g.Ball.Position = mgl32.Vec2{g.Ball.Position.X(), g.Ball.Position.Y() - penetration}
					} else {
						g.Ball.Position = mgl32.Vec2{g.Ball.Position.X(), g.Ball.Position.Y() + penetration}
					}
				}
			}
		}
	}
	if !g.Ball.Stuck {
		collides, _, _ := CheckBallCollision(g.Ball, g.Player)
		if collides {
			centerBoard := g.Player.Position.X() + g.Player.Size.X()/2
			distance := (g.Ball.Position.X() + g.Ball.Radius) - centerBoard
			percentage := distance / (g.Player.Size.X() / 2)

			var strength float32 = 2.0
			oldVelocity := g.Ball.Velocity
			g.Ball.Velocity = mgl32.Vec2{INITIAL_BALL_VELOCITY.X() * percentage * strength, g.Ball.Velocity.Y()}
			g.Ball.Velocity = g.Ball.Velocity.Normalize().Mul(oldVelocity.Len())
			g.Ball.Velocity = mgl32.Vec2{g.Ball.Velocity.X(), float32(-1 * math.Abs(float64(g.Ball.Velocity.Y())))}
		}
	}
}

func CheckCollision(one, two *Object) bool {
	collisionX := one.Position.X()+one.Size.X() >= two.Position.X() && two.Position.X()+two.Size.X() >= one.Position.X()
	collisionY := one.Position.Y()+one.Size.Y() >= two.Position.Y() && two.Position.Y()+two.Size.Y() >= one.Position.Y()
	return collisionX && collisionY
}

func CheckBallCollision(one *Ball, two *Object) (bool, Direction, mgl32.Vec2) {
	center := mgl32.Vec2{one.Position.X() + one.Radius, one.Position.Y() + one.Radius}

	aabbHalfExtents := mgl32.Vec2{two.Size.X() / 2, two.Size.Y() / 2}
	aabbCenter := mgl32.Vec2{two.Position.X() + aabbHalfExtents.X(), two.Position.Y() + aabbHalfExtents.Y()}
	difference := center.Sub(aabbCenter)
	clampedX := mgl32.Clamp(difference.X(), -aabbHalfExtents.X(), aabbHalfExtents.X())
	clampedY := mgl32.Clamp(difference.Y(), -aabbHalfExtents.Y(), aabbHalfExtents.Y())
	closest := aabbCenter.Add(mgl32.Vec2{clampedX, clampedY})
	difference = closest.Sub(center)
	if difference.Len() < one.Radius {
		return true, VectorDirection(difference), difference
	}
	return false, DIRECTION_UP, mgl32.Vec2{}
}

type Direction int

const (
	DIRECTION_UP Direction = iota
	DIRECTION_RIGHT
	DIRECTION_DOWN
	DIRECTION_LEFT
)

func VectorDirection(target mgl32.Vec2) Direction {
	compass := []mgl32.Vec2{
		{0, 1},
		{1, 0},
		{0, -1},
		{-1, 0},
	}
	var max float32 = 0.0
	bestMatch := -1
	for i := 0; i < 4; i++ {
		dotProduct := target.Normalize().Dot(compass[i])
		if dotProduct > max {
			max = dotProduct
			bestMatch = i
		}
	}
	return Direction(bestMatch)
}
