package eng

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Scene interface {
	New(width, height int, window *glfw.Window)
	Render()
	Update(float32)
	Close()
}

func Run(scene Scene, width, height int) {
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

	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	deltaTime := 0.5
	lastFrame := 0.0

	frames := 0
	showFps := time.Tick(1 * time.Second)

	scene.New(width, height, window)

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

		scene.Update(float32(deltaTime))

		gl.ClearColor(0, 0, 0, 0.5)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		scene.Render()
		window.SwapBuffers()
	}

	scene.Close()
}
