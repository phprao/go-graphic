package demo6

// 纹理环绕效果

import (
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/phprao/go-graphic/util"
)

const (
	width  = 800
	height = 600

	vertexShaderSource = `
        #version 410

        in vec3 vPosition;
		in vec3 vColor;
		in vec2 vTexCoord;
		
		out vec3 fColor;
		out vec2 fTexCoord;

        void main() {
            gl_Position = vec4(vPosition, 1.0);
			fColor = vColor;
			fTexCoord = vec2(vTexCoord.x, 1.0-vTexCoord.y);
        }
    ` + "\x00"
	fragmentShaderSource = `
        #version 410
        
		in vec3 fColor;
		in vec2 fTexCoord;

		out vec4 frag_colour;

		uniform sampler2D ourTexture;

        void main() {
			frag_colour = texture(ourTexture, fTexCoord);
        }
    ` + "\x00"
)

var (
	vertices = []float32{
		// Positions   // Colors      // Texture Coords
		0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 2.0, 2.0, // Top Right
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 2.0, 0.0, // Bottom Right
		-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // Bottom Left
		-0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.0, 2.0, // Top Left
	}

	indices = []uint32{
		0, 1, 3, // First Triangle
		1, 2, 3, // Second Triangle
	}
)

func Run() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "texture2d")
	defer glfw.Terminate()

	program, _ := util.InitOpenGL(vertexShaderSource, fragmentShaderSource)
	vao := util.MakeVaoWithAttrib(program, vertices, indices, []util.VertAttrib{{Name: "vPosition", Size: 3}, {Name: "vColor", Size: 3}, {Name: "vTexCoord", Size: 2}})
	pointNum := int32(len(indices))
	texture1 := util.MakeTexture("demo6/round.jpg")

	for !window.ShouldClose() {
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.UseProgram(program)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)

		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, pointNum, gl.UNSIGNED_INT, gl.Ptr(indices))

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
