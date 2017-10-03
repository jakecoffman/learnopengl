package main

import (
	"runtime"

	"fmt"
	"os"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/jakecoffman/learnopengl"
	"github.com/jakecoffman/learnopengl/breakout"
)

const (
	width  = 800
	height = 600
)

var Breakout = breakout.NewGame(width, height)

func main() {
	runtime.LockOSThread()

	// glfw: initialize and configure
	glfw.Init()
	defer glfw.Terminate()
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	if runtime.GOOS == "darwin" {
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	}

	// glfw window creation
	window, err := glfw.CreateWindow(width, height, "Breakout", nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	window.MakeContextCurrent()
	window.SetKeyCallback(keyCallback)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	Breakout.Init()

	shader := breakout.NewShader(fontVs, fontFs)
	projection := mgl32.Ortho2D(0, width, height, 0)
	shader.Use().SetMat4("projection", projection).SetInt("text", 0)

	fd, err := os.Open("text/Roboto-Light.ttf")
	if err != nil {
		panic(err)
	}
	font, err := learnopengl.LoadTrueTypeFont(shader.ID, fd, 50, 32, 127)
	if err != nil {
		panic(err)
	}
	fd.Close()

	deltaTime := 0.5
	lastFrame := 0.0

	frames := 0
	showFps := time.Tick(1 * time.Second)

	for !window.ShouldClose() {
		currentFrame := glfw.GetTime()
		frames++
		select {
		case <-showFps:
			window.SetTitle(fmt.Sprintf("Breakout | %d FPS", frames))
			frames = 0
		default:
		}
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame
		glfw.PollEvents()

		Breakout.ProcessInput(deltaTime)
		Breakout.Update(float32(deltaTime))

		gl.ClearColor(0, 0, 0, 0.5)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		Breakout.Render()
		font.SetColor(1, 1, 1, 1)
		font.Printf(100, 100, 1, "Hello, world!")
		window.SwapBuffers()
	}

	breakout.ResourceManager.Clear()
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
	if key >= 0 && key < 1024 {
		if action == glfw.Press {
			Breakout.Keys[key] = true
		} else if action == glfw.Release {
			Breakout.Keys[key] = false
		}
	}
}

const fontVs = `
#version 330 core
layout (location = 0) in vec4 vertex; // <vec2 pos, vec2 tex>
out vec2 TexCoords;

uniform mat4 projection;

void main()
{
    gl_Position = projection * vec4(vertex.xy, 0.0, 1.0);
    TexCoords = vertex.zw;
}
`

const fontFs = `
#version 330 core
in vec2 TexCoords;
out vec4 color;

uniform sampler2D text;
uniform vec3 textColor;

void main()
{
    vec4 sampled = vec4(1.0, 1.0, 1.0, texture(text, TexCoords).r);
    color = vec4(textColor, 1.0) * sampled;
}
`
