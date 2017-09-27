package main

import (
	_ "image/jpeg"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/jakecoffman/learnopengl"
	"github.com/runningwild/mathgl"
)

const (
	width  = 800
	height = 600
)

func main() {
	runtime.LockOSThread()

	glfw.Init()
	defer glfw.Terminate()
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	if runtime.GOOS == "darwin" {
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	}

	window, err := glfw.CreateWindow(width, height, "LearnOpenGL", nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	window.MakeContextCurrent()
	window.SetFramebufferSizeCallback(framebufferSizeCallback)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	shader := learnopengl.NewShader("4/vertex.glsl", "4/fragment.glsl")

	var vertices = []float32{
		0.5,  0.5, 0.0,   1.0, 1.0, // top right
		0.5, -0.5, 0.0,   1.0, 0.0, // bottom right
		-0.5, -0.5, 0.0,   0.0, 0.0, // bottom left
		-0.5,  0.5, 0.0,   0.0, 1.0,  // top left
	}
	var indices = []int32{
		0, 1, 3,
		1, 2, 3,
	}

	var vbo, vao, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)
	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

	// position
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// texture
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	// load image, create texture and generate mipmaps
	texture1, err := learnopengl.NewTexture("3/container.jpg")
	if err != nil {
		panic(err)
	}
	texture2, err := learnopengl.NewTexture("4/awesomeface.png")
	if err != nil {
		panic(err)
	}
	shader.Use()
	shader.SetInt("texture1", 0)
	shader.SetInt("texture2", 1)

	for !window.ShouldClose() {
		processInput(window)

		gl.ClearColor(0.2, 0.3, 0.3, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.ActiveTexture(gl.TEXTURE0)
		texture1.Bind()
		gl.ActiveTexture(gl.TEXTURE1)
		texture2.Bind()

		// create transforms
		transform := mathgl.Mat4{}
		transform.Translation(0.5, -0.5, 0)
		transform.RotationAxisAngle(mathgl.Vec3{0, 0, 1}, float32(glfw.GetTime()))

		shader.Use()
		transformLoc := gl.GetUniformLocation(shader.ID, gl.Str("transform\x00"))
		gl.UniformMatrix4fv(transformLoc, 1, false, &transform[0])

		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

		window.SwapBuffers()
		glfw.PollEvents()
	}

	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)
	gl.DeleteBuffers(1, &ebo)
}

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}

func framebufferSizeCallback(w *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}
