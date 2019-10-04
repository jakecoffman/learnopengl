package breakout

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/jakecoffman/learnopengl/breakout/eng"
	"math"
)

type Game struct {
	state         int
	Keys          [1024]bool
	Width, Height int

	Levels []*Level
	Level  int

	Player *Object
	Ball   *Ball

	*eng.ResourceManager
	ParticleGenerator *eng.ParticleGenerator
	SpriteRenderer    *eng.SpriteRenderer
	TextRenderer      *eng.TextRenderer
}

// Game state
const (
	stateActive = iota
	stateMenu
	stateWin
)

var (
	playerSize          = mgl32.Vec2{100, 20}
	playerVelocity      = float32(500.0)
	initialBallVelocity = Vec2(100, -350)
	ballRadius          = float32(12.5)
)

func (g *Game) New(w, h int, window *glfw.Window) {
	g.Width =  w
	g.Height = h
	g.Keys =[1024]bool{}
	g.ResourceManager = eng.NewResourceManager()

	width, height := float32(g.Width), float32(g.Height)

	g.LoadShader("breakout/shaders/main.vs.glsl", "breakout/shaders/main.fs.glsl", "sprite")
	g.LoadShader("breakout/shaders/particle.vs.glsl", "breakout/shaders/particle.fs.glsl", "particle")

	projection := mgl32.Ortho(0, width, height, 0, -1, 1)
	g.Shader("sprite").Use().
		SetInt("sprite", 0).
		SetMat4("projection", projection)
	g.Shader("particle").Use().
		SetInt("sprite", 0).
		SetMat4("projection", projection)

	g.LoadTexture("breakout/textures/background.jpg", "background")
	g.LoadTexture("breakout/textures/paddle.png", "paddle")
	g.LoadTexture("breakout/textures/particle.png", "particle")
	g.LoadTexture("breakout/textures/awesomeface.png", "face")
	block := g.LoadTexture("breakout/textures/block.png", "block")
	solid := g.LoadTexture("breakout/textures/block_solid.png", "block_solid")

	shader := g.LoadShader("breakout/shaders/text.vs.glsl", "breakout/shaders/text.fs.glsl", "text")
	g.TextRenderer = eng.NewTextRenderer(shader, width, height, "breakout/textures/Roboto-Light.ttf", 24)
	g.TextRenderer.SetColor(1, 1, 1, 1)

	g.ParticleGenerator = eng.NewParticleGenerator(g.Shader("particle"), g.Texture("particle"), 500)
	g.SpriteRenderer = eng.NewSpriteRenderer(g.Shader("sprite"))

	one := NewLevel(block, solid)
	one.Load("breakout/levels/1.txt", g.Width, int(float32(g.Height)*0.5))
	g.Levels = append(g.Levels, one)

	playerPos := mgl32.Vec2{float32(g.Width)/2.0 - playerSize.X()/2.0, float32(g.Height) - playerSize.Y()}
	g.Player = NewGameObject(playerPos, playerSize, g.Texture("paddle"))

	ballPos := playerPos.Add(mgl32.Vec2{playerSize.X()/2.0 - ballRadius, -ballRadius * 2})
	g.Ball = NewBall(ballPos, ballRadius, initialBallVelocity, g.Texture("face"))

	g.state = stateActive

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
}

func (g *Game) Update(dt float32) {
	g.processInput(dt)
	g.Ball.Move(dt, float32(g.Width))
	g.doCollisions()
	ball := g.Ball.Object
	g.ParticleGenerator.Update(dt, ball.Position, ball.Velocity, 2, mgl32.Vec2{g.Ball.Radius / 2, g.Ball.Radius / 2})
	if g.Ball.Position.Y() >= float32(g.Height) {
		g.resetLevel()
		g.resetPlayer()
	}
}

func (g *Game) Render() {
	if g.state == stateActive {
		g.SpriteRenderer.DrawSprite(g.Texture("background"), Vec2(0, 0), Vec2(g.Width, g.Height), 0, eng.DefaultColor)
		g.Levels[g.Level].Draw(g.SpriteRenderer)
		g.Player.Draw(g.SpriteRenderer)
		g.ParticleGenerator.Draw()
		g.Ball.Draw(g.SpriteRenderer)
	}
	g.TextRenderer.Print("Hello, world!", 10, 25, 1)
}

func (g *Game) Close() {
	g.Clear()
}

func (g *Game) processInput(dt float32) {
	if g.state != stateActive {
		return
	}

	velocity := playerVelocity * dt

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

func (g *Game) pause() {
	g.state = stateMenu
}

func (g *Game) unpause() {
	g.state = stateActive
}

func (g *Game) resetLevel() {
	if g.Level == 0 {
		g.Levels[0].Load("breakout/levels/1.txt", g.Width, int(float32(g.Height)*0.5))
	}
	// TODO
}

func (g *Game) resetPlayer() {
	g.Player.Size = playerSize
	g.Player.Position = mgl32.Vec2{float32(g.Width)/2 - playerSize.X()/2, float32(g.Height) - playerSize.Y()}
	g.Ball.Reset(g.Player.Position.Add(mgl32.Vec2{playerSize.X()/2 - ballRadius, -(ballRadius * 2)}), initialBallVelocity)
}

func (g *Game) doCollisions() {
	for _, box := range g.Levels[g.Level].Bricks {
		if !box.Destroyed {
			collides, direction, diff := checkBallCollision(g.Ball, box)
			if collides {
				if !box.IsSolid {
					box.Destroyed = true
				}

				if direction == directionLeft || direction == directionRight {
					g.Ball.Velocity = mgl32.Vec2{-g.Ball.Velocity.X(), g.Ball.Velocity.Y()}
					penetration := g.Ball.Radius - float32(math.Abs(float64(diff.X())))
					if direction == directionLeft {
						g.Ball.Position = mgl32.Vec2{g.Ball.Position.X() + penetration, g.Ball.Position.Y()}
					} else {
						g.Ball.Position = mgl32.Vec2{g.Ball.Position.X() - penetration, g.Ball.Position.Y()}
					}
				} else {
					g.Ball.Velocity = mgl32.Vec2{g.Ball.Velocity.X(), -g.Ball.Velocity.Y()}
					penetration := g.Ball.Radius - float32(math.Abs(float64(diff.Y())))
					if direction == directionUp {
						g.Ball.Position = mgl32.Vec2{g.Ball.Position.X(), g.Ball.Position.Y() - penetration}
					} else {
						g.Ball.Position = mgl32.Vec2{g.Ball.Position.X(), g.Ball.Position.Y() + penetration}
					}
				}
			}
		}
	}
	if !g.Ball.Stuck {
		collides, _, _ := checkBallCollision(g.Ball, g.Player)
		if collides {
			centerBoard := g.Player.Position.X() + g.Player.Size.X()/2
			distance := (g.Ball.Position.X() + g.Ball.Radius) - centerBoard
			percentage := distance / (g.Player.Size.X() / 2)

			var strength float32 = 2.0
			oldVelocity := g.Ball.Velocity
			g.Ball.Velocity = mgl32.Vec2{initialBallVelocity.X() * percentage * strength, g.Ball.Velocity.Y()}
			g.Ball.Velocity = g.Ball.Velocity.Normalize().Mul(oldVelocity.Len())
			g.Ball.Velocity = mgl32.Vec2{g.Ball.Velocity.X(), float32(-1 * math.Abs(float64(g.Ball.Velocity.Y())))}
		}
	}
}

func checkCollision(one, two *Object) bool {
	collisionX := one.Position.X()+one.Size.X() >= two.Position.X() && two.Position.X()+two.Size.X() >= one.Position.X()
	collisionY := one.Position.Y()+one.Size.Y() >= two.Position.Y() && two.Position.Y()+two.Size.Y() >= one.Position.Y()
	return collisionX && collisionY
}

func checkBallCollision(one *Ball, two *Object) (bool, Direction, mgl32.Vec2) {
	center := mgl32.Vec2{one.Position.X() + one.Radius, one.Position.Y() + one.Radius}

	aabbHalfExtents := mgl32.Vec2{two.Size.X() / 2, two.Size.Y() / 2}
	aabbCenter := mgl32.Vec2{two.Position.X() + aabbHalfExtents.X(), two.Position.Y() + aabbHalfExtents.Y()}
	difference := center.Sub(aabbCenter)
	clampedX := mgl32.Clamp(difference.X(), -aabbHalfExtents.X(), aabbHalfExtents.X())
	clampedY := mgl32.Clamp(difference.Y(), -aabbHalfExtents.Y(), aabbHalfExtents.Y())
	closest := aabbCenter.Add(mgl32.Vec2{clampedX, clampedY})
	difference = closest.Sub(center)
	if difference.Len() < one.Radius {
		return true, vectorDirection(difference), difference
	}
	return false, directionUp, mgl32.Vec2{}
}

type Direction int

const (
	directionUp Direction = iota
	directionRight
	directionDown
	directionLeft
)

func vectorDirection(target mgl32.Vec2) Direction {
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
