package breakout

import "github.com/go-gl/mathgl/mgl32"

type Game struct {
	state         int
	Keys          [1024]bool
	Width, Height int
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

var Renderer *SpriteRenderer

func (g *Game) Init() {
	ResourceManager.LoadShader("breakout/vertex.glsl", "breakout/fragment.glsl", "sprite")

	projection := mgl32.Ortho(0, float32(g.Width), float32(g.Height), 0, -1, 1)
	ResourceManager.Shader("sprite").Use().SetInt("image", 0)
	ResourceManager.Shader("sprite").SetMat4("projection", projection)

	Renderer = NewSpriteRenderer(ResourceManager.Shader("sprite"))

	ResourceManager.LoadTexture("breakout/awesomeface.png", true, "face")

	g.state = GAME_ACTIVE
}

func (g *Game) ProcessInput(dt float64) {

}

func (g *Game) Update(dt float64) {

}

func (g *Game) Render() {
	Renderer.DrawSprite(ResourceManager.Texture("face"), mgl32.Vec2{200, 200}, mgl32.Vec2{300, 400}, 45, mgl32.Vec3{0, 1, 0})
}

func (g *Game) Pause() {
	g.state = GAME_MENU
}

func (g *Game) Unpause() {
	g.state = GAME_ACTIVE
}
