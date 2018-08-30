package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
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

	text := breakout.NewTextRenderer("breakout/shaders/text.vs.glsl", "breakout/shaders/text.fs.glsl", width, height)
	if err = text.Load("breakout/textures/Roboto-Light.ttf", 24); err != nil {
		panic(err)
	}
	text.SetColor(1, 1, 1, 1)

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
		text.Print("Hello, world!", 10, 25, 1)
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
