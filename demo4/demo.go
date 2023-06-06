package demo4

// 纹理贴图

import (
	"math"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
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

	vertexShaderSource2 = `
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
	fragmentShaderSource2 = `
        #version 410
        
		in vec3 fColor;
		in vec2 fTexCoord;

		out vec4 frag_colour;

		uniform sampler2D ourTexture1;
		uniform sampler2D ourTexture2;

        void main() {
			frag_colour = mix(texture(ourTexture1, fTexCoord), texture(ourTexture2, fTexCoord), 0.2);
        }
    ` + "\x00"

	vertexShaderSource3 = `
        #version 410

        in vec3 vPosition;
		in vec3 vColor;
		in vec2 vTexCoord;
		
		out vec3 fColor;
		out vec2 fTexCoord;

		uniform mat4 transe;

        void main() {
            gl_Position = transe * vec4(vPosition, 1.0);
			fColor = vColor;
			fTexCoord = vec2(vTexCoord.x, 1.0-vTexCoord.y);
        }
    ` + "\x00"
	fragmentShaderSource3 = `
        #version 410
        
		in vec3 fColor;
		in vec2 fTexCoord;

		out vec4 frag_colour;

		uniform sampler2D ourTexture1;
		uniform sampler2D ourTexture2;

        void main() {
			frag_colour = mix(texture(ourTexture1, fTexCoord), texture(ourTexture2, fTexCoord), 0.2);
        }
    ` + "\x00"
)

var (
	vertices = []float32{
		// Positions   // Colors      // Texture Coords
		0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, // Top Right
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, // Bottom Right
		-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // Bottom Left
		-0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, // Top Left
	}

	indices = []uint32{
		0, 1, 3, // First Triangle
		1, 2, 3, // Second Triangle
	}
)

// 一个纹理
func Run() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "texture2d")
	defer glfw.Terminate()

	program, _ := util.InitOpenGL(vertexShaderSource, fragmentShaderSource)
	vao := util.MakeVaoWithAttrib(program, vertices, indices, []util.VertAttrib{{Name: "vPosition", Size: 3}, {Name: "vColor", Size: 3}, {Name: "vTexCoord", Size: 2}})
	pointNum := int32(len(indices))
	texture1 := util.MakeTexture("demo4/container.jpg")

	for !window.ShouldClose() {
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.UseProgram(program)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)

		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, pointNum, gl.UNSIGNED_INT, gl.Ptr(indices))

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

// 两个纹理
func Run2() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "texture2d")
	defer glfw.Terminate()

	program, _ := util.InitOpenGL(vertexShaderSource2, fragmentShaderSource2)
	vao := util.MakeVaoWithAttrib(program, vertices, indices, []util.VertAttrib{{Name: "vPosition", Size: 3}, {Name: "vColor", Size: 3}, {Name: "vTexCoord", Size: 2}})
	pointNum := int32(len(indices))
	texture1 := util.MakeTexture("demo4/container.jpg")
	texture2 := util.MakeTexture("demo4/awesomeface.png")

	for !window.ShouldClose() {
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.UseProgram(program)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture1"+"\x00")), 0)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture2"+"\x00")), 1)

		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, pointNum, gl.UNSIGNED_INT, gl.Ptr(indices))

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

// 缩放，旋转，移动
func Run3() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "texture2d")
	defer glfw.Terminate()

	program, _ := util.InitOpenGL(vertexShaderSource3, fragmentShaderSource3)
	vao := util.MakeVaoWithAttrib(program, vertices, indices, []util.VertAttrib{{Name: "vPosition", Size: 3}, {Name: "vColor", Size: 3}, {Name: "vTexCoord", Size: 2}})
	pointNum := int32(len(indices))
	texture1 := util.MakeTexture("demo4/container.jpg")
	texture2 := util.MakeTexture("demo4/awesomeface.png")

	for !window.ShouldClose() {
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.UseProgram(program)

		rotate := mgl32.HomogRotate3D(mgl32.DegToRad(90), mgl32.Vec3{0, 0, 1})
		// rotate := mgl32.HomogRotate3D(float32(glfw.GetTime()), mgl32.Vec3{0, 0, 1}) // 旋转效果
		scale := mgl32.Scale3D(0.5, 0.5, 0.5)
		translate := mgl32.Translate3D(0.5, -0.5, 0)
		// 顺序要反着看：依次是 scale，rotate，translate
		transe := translate.Mul4(rotate).Mul4(scale)
		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("transe\x00")), 1, false, &transe[0])

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture1"+"\x00")), 0)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture2"+"\x00")), 1)

		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, pointNum, gl.UNSIGNED_INT, gl.Ptr(indices))

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

// 画两个箱子，缩放，旋转，移动
func Run4() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "texture2d")
	defer glfw.Terminate()

	program, _ := util.InitOpenGL(vertexShaderSource3, fragmentShaderSource3)
	vao := util.MakeVaoWithAttrib(program, vertices, indices, []util.VertAttrib{{Name: "vPosition", Size: 3}, {Name: "vColor", Size: 3}, {Name: "vTexCoord", Size: 2}})
	pointNum := int32(len(indices))
	texture1 := util.MakeTexture("demo4/container.jpg")
	texture2 := util.MakeTexture("demo4/awesomeface.png")

	for !window.ShouldClose() {
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.UseProgram(program)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture1"+"\x00")), 0)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture2"+"\x00")), 1)

		gl.BindVertexArray(vao)

		// 第一个箱子
		rotate := mgl32.HomogRotate3D(float32(glfw.GetTime()), mgl32.Vec3{0, 0, 1}) // 旋转效果
		scale := mgl32.Scale3D(0.5, 0.5, 0.5)
		translate := mgl32.Translate3D(0.5, -0.5, 0)
		transe := translate.Mul4(rotate).Mul4(scale)
		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("transe\x00")), 1, false, &transe[0])
		gl.DrawElements(gl.TRIANGLES, pointNum, gl.UNSIGNED_INT, gl.Ptr(indices))

		// 第二个箱子
		rotate2 := mgl32.HomogRotate3D(mgl32.DegToRad(90), mgl32.Vec3{0, 0, 1})
		s := float32(math.Sin(glfw.GetTime()))
		scale2 := mgl32.Scale3D(s, s, s)
		translate2 := mgl32.Translate3D(-0.5, 0.5, 0)
		transe2 := translate2.Mul4(rotate2).Mul4(scale2)
		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("transe\x00")), 1, false, &transe2[0])
		gl.DrawElements(gl.TRIANGLES, pointNum, gl.UNSIGNED_INT, gl.Ptr(indices))

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
