package breakout

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Game struct {
	state         int
	Keys          [1024]bool
	Width, Height int

	Levels []*Level
	Level  int

	Player *Object

	renderer *SpriteRenderer
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
	PLAYER_SIZE     = mgl32.Vec2{100, 20}
	PLAYER_VELOCITY = 500.0
)

func (g *Game) Init() {
	ResourceManager.LoadShader("breakout/vertex.glsl", "breakout/fragment.glsl", "sprite")

	projection := mgl32.Ortho(0, float32(g.Width), float32(g.Height), 0, -1, 1)
	ResourceManager.Shader("sprite").
		Use().
		SetInt("image", 0).
		SetMat4("projection", projection)

	g.renderer = NewSpriteRenderer(ResourceManager.Shader("sprite"))

	ResourceManager.LoadTexture("breakout/background.jpg", "background")
	ResourceManager.LoadTexture("breakout/paddle.png", "paddle")
	ResourceManager.LoadTexture("breakout/awesomeface.png", "face")
	ResourceManager.LoadTexture("breakout/block.png", "block")
	ResourceManager.LoadTexture("breakout/block_solid.png", "block_solid")

	one := NewLevel()
	if err := one.Load("breakout/level1.txt", g.Width, int(float32(g.Height)*0.5)); err != nil {
		panic(err)
	}
	g.Levels = append(g.Levels, one)

	playerPos := mgl32.Vec2{float32(g.Width)/2.0 - PLAYER_SIZE.X()/2.0, float32(g.Height) - PLAYER_SIZE.Y()}
	g.Player = NewGameObject(playerPos, PLAYER_SIZE, ResourceManager.Texture("paddle"))

	g.state = GAME_ACTIVE
}

func (g *Game) ProcessInput(dt float64) {

}

func (g *Game) Update(dt float64) {

}

func (g *Game) Render() {
	if g.state == GAME_ACTIVE {
		g.renderer.DrawSprite(ResourceManager.Texture("background"), Vec2(0, 0), Vec2(g.Width, g.Height), 0, DefaultColor)
		g.Levels[g.Level].Draw(g.renderer)
		g.Player.Draw(g.renderer)
	}
}

func (g *Game) Pause() {
	g.state = GAME_MENU
}

func (g *Game) Unpause() {
	g.state = GAME_ACTIVE
}
