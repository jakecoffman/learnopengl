package main

import (
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/jakecoffman/learnopengl"
)

const (
	width  = 800
	height = 600
)

func main() {
	// Go will try to move this to another thread without this
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

	// build and compile our shader programs
	vertexShader := learnopengl.CompileShader(gl.VERTEX_SHADER, vertexShaderSource)
	fragmentShader := learnopengl.CompileShader(gl.FRAGMENT_SHADER, fragmentShaderSource)
	shaderProgram := learnopengl.LinkProgram(vertexShader, fragmentShader)
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	// set up vertex data (and buffer(s)) and configure vertex attributes
	var vertices = []float32{
		0.5, 0.5, 0, // top right
		0.5, -0.5, 0, // bottom right
		-0.5, -0.5, 0, // bottom left
		-0.5, 0.5, 0, // top left
	}
	var indices = []uint32{
		0, 1, 3, // first triangle
		1, 2, 3, // second triangle
	}

	// VBO stores vertices in memory on the GPU
	// VAO defines the data layout (*VertexAttrib* calls)
	// EBO (element buffer) allows storing indices of what to draw, to make the VBO smaller (*Element* calls)
	var vbo, vao, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)
	// bind the Vertex Array Object first, then bind and set vertex buffer(s), and then configure vertex attributes(s).
	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	// Can't take unsafe.Sizeof an array/slice in Go, it returns the size of the header.
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

	// Tells how to interpret the VBO data (stored in the VAO)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// note that this is allowed, the call to glVertexAttribPointer registered VBO as
	// the vertex attribute's bound vertex buffer object so afterwards we can safely unbind
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// remember: do NOT unbind the EBO while a VAO is active as the bound element buffer
	// object IS stored in the VAO; keep the EBO bound.
	//gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0);

	// You can unbind the VAO afterwards so other VAO calls won't accidentally modify this
	// VAO, but this rarely happens. Modifying other VAOs requires a call to glBindVertexArray
	// anyways so we generally don't unbind VAOs (nor VBOs) when it's not directly necessary.
	gl.BindVertexArray(0)

	// uncomment this call to draw in wireframe polygons.
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	for !window.ShouldClose() {
		processInput(window)

		gl.ClearColor(0.2, 0.3, 0.3, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.UseProgram(shaderProgram)
		// seeing as we only have a single VAO there's no need to bind it every time,
		// but we'll do so to keep things a bit more organized
		gl.BindVertexArray(vao)
		// draw the elements bound above (indices)
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

// The vertex shader takes a single vertex as input and translates it.
// OpenGL uses a range from -1 to 1 for all coordinates, which is why this shader is required.
var vertexShaderSource = `#version 330 core
layout (location = 0) in vec3 aPos;
void main()
{
   gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
}` + "\x00"

// The fragment shader calculates the final color of a pixel and applies lighting, shadows, etc.
var fragmentShaderSource = `#version 330 core
out vec4 FragColor;
void main()
{
   FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
}` + "\x00"
