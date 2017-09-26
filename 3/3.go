package main

import (
	_ "image/jpeg"
	"runtime"

	"image"
	"image/draw"
	"os"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/jakecoffman/learnopengl"
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

	shader := learnopengl.NewShader("3/vertex.glsl", "3/fragment.glsl")

	var vertices = []float32{
		0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, // top right
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, // bottom right
		-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom left
		-0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, // top left
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

	gl.BindBuffer(gl.ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

	// position
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// color
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	// texture
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	// load and create a texture
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture) // all upcoming GL_TEXTURE_2D operations now have effect on this texture object
	// set the texture wrapping parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	// set texture filtering parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	// load image, create texture and generate mipmaps
	f, err := os.Open("3/container.jpg")
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	size := rgba.Rect.Size()

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.SRGB, int32(size.X), int32(size.Y), 0, gl.RGB, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	f.Close()

	for !window.ShouldClose() {
		processInput(window)

		gl.ClearColor(0.2, 0.3, 0.3, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.BindTexture(gl.TEXTURE_2D, texture)

		shader.Use()
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
