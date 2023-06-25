package demo4

// 纹理贴图

import (
	"image/gif"
	"math"
	"os"
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

	vertexShaderSource6 = `
        #version 410

        in vec3 vPosition;
		in vec2 vTexCoord;
		
		out vec2 fTexCoord;

		uniform mat4 transe;

        void main() {
            gl_Position = transe * vec4(vPosition, 1.0);
			fTexCoord = vec2(vTexCoord.x, 1.0-vTexCoord.y);
        }
    ` + "\x00"
	fragmentShaderSource6 = `
        #version 410
        
		in vec2 fTexCoord;

		out vec4 frag_colour;

		uniform sampler2D ourTexture1;
		uniform sampler2D ourTexture2;

        void main() {
			frag_colour = mix(texture(ourTexture1, fTexCoord), texture(ourTexture2, fTexCoord), 0.5);
        }
    ` + "\x00"

	vertexShaderSource7 = `
        #version 410

        in vec3 vPosition;
		in vec2 vTexCoord;
		
		out vec2 fTexCoord;

		uniform mat4 transe;

        void main() {
            gl_Position = transe * vec4(vPosition, 1.0);
			fTexCoord = vec2(vTexCoord.x, 1.0-vTexCoord.y);
        }
    ` + "\x00"
	fragmentShaderSource7 = `
        #version 410
        
		in vec2 fTexCoord;

		out vec4 frag_colour;

		uniform sampler2D ourTexture;

        void main() {
			frag_colour = texture(ourTexture, fTexCoord);
        }
    ` + "\x00"
	vertexShaderSource8 = `
        #version 410

        in vec3 vPosition;
		out vec3 textureDir;

		uniform mat4 transe;

        void main() {
            gl_Position = transe * vec4(vPosition, 1.0);
			textureDir = vPosition;
        }
    ` + "\x00"
	fragmentShaderSource8 = `
        #version 410
        
		in vec3 textureDir;

		out vec4 frag_colour;

		uniform samplerCube cubemap;

        void main() {
			frag_colour = texture(cubemap, textureDir);
        }
    ` + "\x00"

	vertexShaderSource9 = `
        #version 410

        in vec3 vPosition;
		in vec2 vTexCoord;
		
		out vec2 fTexCoord;

		uniform mat4 transe;

        void main() {
            gl_Position = transe * vec4(vPosition, 1.0);
			fTexCoord = vec2(vTexCoord.x, 1.0-vTexCoord.y);
        }
    ` + "\x00"
	fragmentShaderSource9 = `
        #version 410
        
		in vec2 fTexCoord;

		out vec4 frag_colour;

		uniform sampler2D ourTexture1;

        void main() {
			frag_colour = texture(ourTexture1, fTexCoord);
        }
    ` + "\x00"

	vertexShaderSource10 = `
        #version 410

        in vec2 vPosition;

        void main() {
            gl_Position = vec4(vPosition, 1.0, 1.0);
        }
    ` + "\x00"
	fragmentShaderSource10 = `
        #version 410
        
		out vec4 frag_colour;

        void main() {
			frag_colour = vec4(1, 1, 1, 1);
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

	vertices36 = []float32{
		// Right
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		// Left
		-0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0,
		// Top
		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		// Bottom
		-0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		// back
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,
		// Front
		-0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
	}

	vertices37 = []float32{
		-1.0, 1.0, -1.0,
		-1.0, -1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, -1.0,
		-1.0, 1.0, -1.0,

		-1.0, -1.0, 1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, -1.0,
		-1.0, 1.0, -1.0,
		-1.0, 1.0, 1.0,
		-1.0, -1.0, 1.0,

		1.0, -1.0, -1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, -1.0,

		-1.0, -1.0, 1.0,
		-1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,
		-1.0, -1.0, 1.0,

		-1.0, 1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,

		-1.0, -1.0, -1.0,
		-1.0, -1.0, 1.0,
		1.0, -1.0, -1.0,
		1.0, -1.0, -1.0,
		-1.0, -1.0, 1.0,
		1.0, -1.0, 1.0,
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

// model, view, projection
func Run5() {
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

		model := mgl32.HomogRotate3D(mgl32.DegToRad(-55), mgl32.Vec3{1, 0, 0})
		view := mgl32.Translate3D(0, 0, -3)
		projection := mgl32.Perspective(mgl32.DegToRad(45), float32(width)/float32(height), 0.1, 100)
		transe := projection.Mul4(view).Mul4(model)
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

// 立方体，model, view, projection，旋转
func Run6() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "texture2d")
	defer glfw.Terminate()

	program, _ := util.InitOpenGL(vertexShaderSource6, fragmentShaderSource6)
	vao := util.MakeVaoWithAttrib(program, vertices36, nil, []util.VertAttrib{{Name: "vPosition", Size: 3}, {Name: "vTexCoord", Size: 2}})
	pointNum := int32(len(vertices36)) / 5
	texture1 := util.MakeTexture("demo4/container.jpg")
	texture2 := util.MakeTexture("demo4/awesomeface.png")

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Enable(gl.DEPTH_TEST)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		model := mgl32.HomogRotate3D(float32(glfw.GetTime()), mgl32.Vec3{0.5, 1, 0})
		view := mgl32.Translate3D(0, 0, -3)
		projection := mgl32.Perspective(mgl32.DegToRad(45), float32(width)/height, 0.1, 100)
		transe := projection.Mul4(view).Mul4(model)
		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("transe\x00")), 1, false, &transe[0])

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture1"+"\x00")), 0)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture2"+"\x00")), 1)

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, pointNum)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

var more = []mgl32.Vec3{
	{0.0, 0.0, 0.0},
	{2.0, 5.0, -15.0},
	{-1.5, -2.2, -2.5},
	{-3.8, -2.0, -12.3},
	{2.4, -0.4, -3.5},
	{-1.7, 3.0, -7.5},
	{1.3, -2.0, -2.5},
	{1.5, 2.0, -2.5},
	{1.5, 0.2, -1.5},
	{-1.3, 1.0, -1.5},
}

// 10个立方体，model, view, projection
func Run7() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "texture2d")
	defer glfw.Terminate()

	program, _ := util.InitOpenGL(vertexShaderSource6, fragmentShaderSource6)
	vao := util.MakeVaoWithAttrib(program, vertices36, nil, []util.VertAttrib{{Name: "vPosition", Size: 3}, {Name: "vTexCoord", Size: 2}})
	pointNum := int32(len(vertices36)) / 5
	texture1 := util.MakeTexture("demo4/container.jpg")
	texture2 := util.MakeTexture("demo4/awesomeface.png")

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Enable(gl.DEPTH_TEST)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture1"+"\x00")), 0)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture2"+"\x00")), 1)

		gl.BindVertexArray(vao)

		view := mgl32.Translate3D(0, 0, -3)
		projection := mgl32.Perspective(mgl32.DegToRad(45), float32(width)/height, 0.1, 100)

		for k, v := range more {
			model := mgl32.HomogRotate3D(mgl32.DegToRad(float32(20*k)), mgl32.Vec3{1, 0.3, 0.5})
			model = mgl32.Translate3D(v.X(), v.Y(), v.Z()).Mul4(model)
			transe := projection.Mul4(view).Mul4(model)
			gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("transe\x00")), 1, false, &transe[0])
			gl.DrawArrays(gl.TRIANGLES, 0, pointNum)
		}

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

// 立方体旋转，摄像机不动
func Run8() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "texture2d")
	defer glfw.Terminate()

	program, _ := util.InitOpenGL(vertexShaderSource6, fragmentShaderSource6)
	vao := util.MakeVaoWithAttrib(program, vertices36, nil, []util.VertAttrib{{Name: "vPosition", Size: 3}, {Name: "vTexCoord", Size: 2}})
	pointNum := int32(len(vertices36)) / 5
	texture1 := util.MakeTexture("demo4/container.jpg")
	texture2 := util.MakeTexture("demo4/awesomeface.png")

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Enable(gl.DEPTH_TEST)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		model := mgl32.HomogRotate3D(float32(glfw.GetTime()), mgl32.Vec3{0, 1, 0})
		camera := mgl32.LookAtV(mgl32.Vec3{2, 2, 2}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
		projection := mgl32.Perspective(mgl32.DegToRad(45), float32(width)/height, 0.1, 100)
		transe := projection.Mul4(camera).Mul4(model)
		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("transe\x00")), 1, false, &transe[0])

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture1"+"\x00")), 0)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture2"+"\x00")), 1)

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, pointNum)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

// 立方体不动，摄像机围绕着一个圆旋转
func Run9() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "texture2d")
	defer glfw.Terminate()

	program, _ := util.InitOpenGL(vertexShaderSource6, fragmentShaderSource6)
	vao := util.MakeVaoWithAttrib(program, vertices36, nil, []util.VertAttrib{{Name: "vPosition", Size: 3}, {Name: "vTexCoord", Size: 2}})
	pointNum := int32(len(vertices36)) / 5
	texture1 := util.MakeTexture("demo4/container.jpg")
	texture2 := util.MakeTexture("demo4/awesomeface.png")

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Enable(gl.DEPTH_TEST)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		radius := 3.0
		cx := float32(math.Sin(glfw.GetTime()) * radius)
		cz := float32(math.Cos(glfw.GetTime()) * radius)
		camera := mgl32.LookAtV(mgl32.Vec3{cx, 2, cz}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
		projection := mgl32.Perspective(mgl32.DegToRad(45), float32(width)/height, 0.1, 100)
		transe := projection.Mul4(camera)
		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("transe\x00")), 1, false, &transe[0])

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture1"+"\x00")), 0)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture2"+"\x00")), 1)

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, pointNum)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

// 立方体不动，摄像机受到按键WSAD上下左右移动，按下W键不释放，就会一直移动。
func Run10() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "texture2d")
	defer glfw.Terminate()

	program, _ := util.InitOpenGL(vertexShaderSource6, fragmentShaderSource6)
	vao := util.MakeVaoWithAttrib(program, vertices36, nil, []util.VertAttrib{{Name: "vPosition", Size: 3}, {Name: "vTexCoord", Size: 2}})
	pointNum := int32(len(vertices36)) / 5
	texture1 := util.MakeTexture("demo4/container.jpg")
	texture2 := util.MakeTexture("demo4/awesomeface.png")

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Enable(gl.DEPTH_TEST)

	cameraPos := mgl32.Vec3{0, 0, 3}
	cameraFront := mgl32.Vec3{0, 0, -1}
	cameraUp := mgl32.Vec3{0, 1, 0}

	///////////////////////////////////////////////
	keyCallback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		cameraSpeed := float32(0.05)
		if window.GetKey(glfw.KeyW) == glfw.Press {
			cameraPos = cameraPos.Sub(cameraFront.Mul(cameraSpeed))
		}
		if window.GetKey(glfw.KeyS) == glfw.Press {
			cameraPos = cameraPos.Add(cameraFront.Mul(cameraSpeed))
		}
		if window.GetKey(glfw.KeyA) == glfw.Press {
			cameraPos = cameraPos.Add(cameraFront.Cross(cameraUp).Normalize().Mul(cameraSpeed))
		}
		if window.GetKey(glfw.KeyD) == glfw.Press {
			cameraPos = cameraPos.Sub(cameraFront.Cross(cameraUp).Normalize().Mul(cameraSpeed))
		}
		// log.Println(cameraPos, cameraPos.Add(cameraFront))
		// log.Printf("key:%d, scancode:%d, action:%d, mods:%v, cameraPos:%v\n", key, scancode, action, mods, cameraPos)
	}
	window.SetKeyCallback(keyCallback)
	///////////////////////////////////////////////

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		// 这样能保证无论我们怎么移动，摄像机都会注视着目标方向
		camera := mgl32.LookAtV(cameraPos, cameraPos.Add(cameraFront), cameraUp)
		projection := mgl32.Perspective(mgl32.DegToRad(45), float32(width)/height, 0.1, 100)
		transe := projection.Mul4(camera)
		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("transe\x00")), 1, false, &transe[0])

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture1"+"\x00")), 0)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture2"+"\x00")), 1)

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, pointNum)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

// 立方体不动，摄像机受到按键WSAD上下左右移动，按下W键不释放，就会一直移动。
// 实现方法2
func Run11() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "texture2d")
	defer glfw.Terminate()

	program, _ := util.InitOpenGL(vertexShaderSource6, fragmentShaderSource6)
	vao := util.MakeVaoWithAttrib(program, vertices36, nil, []util.VertAttrib{{Name: "vPosition", Size: 3}, {Name: "vTexCoord", Size: 2}})
	pointNum := int32(len(vertices36)) / 5
	texture1 := util.MakeTexture("demo4/container.jpg")
	texture2 := util.MakeTexture("demo4/awesomeface.png")

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Enable(gl.DEPTH_TEST)

	cameraPos := mgl32.Vec3{0, 0, 3}
	cameraFront := mgl32.Vec3{0, 0, -1}
	cameraUp := mgl32.Vec3{0, 1, 0}

	///////////////////////////////////////////////
	cameraSpeed := float32(0.01)
	var holdKey glfw.Key
	keyCallback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyW {
			if action == glfw.Press {
				holdKey = glfw.KeyW
			} else {
				holdKey = 0
			}
		}
		if key == glfw.KeyS {
			if action == glfw.Press {
				holdKey = glfw.KeyS
			} else {
				holdKey = 0
			}
		}
		if key == glfw.KeyA {
			if action == glfw.Press {
				holdKey = glfw.KeyA
			} else {
				holdKey = 0
			}
		}
		if key == glfw.KeyD {
			if action == glfw.Press {
				holdKey = glfw.KeyD
			} else {
				holdKey = 0
			}
		}
		// log.Println(cameraPos, cameraPos.Add(cameraFront))
	}
	window.SetKeyCallback(keyCallback)
	///////////////////////////////////////////////

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		if glfw.KeyW == holdKey {
			cameraPos = cameraPos.Sub(cameraFront.Mul(cameraSpeed))
		}
		if glfw.KeyS == holdKey {
			cameraPos = cameraPos.Add(cameraFront.Mul(cameraSpeed))
		}
		if glfw.KeyA == holdKey {
			cameraPos = cameraPos.Add(cameraFront.Cross(cameraUp).Normalize().Mul(cameraSpeed))
		}
		if glfw.KeyD == holdKey {
			cameraPos = cameraPos.Sub(cameraFront.Cross(cameraUp).Normalize().Mul(cameraSpeed))
		}

		// 这样能保证无论我们怎么移动，摄像机都会注视着目标方向
		camera := mgl32.LookAtV(cameraPos, cameraPos.Add(cameraFront), cameraUp)
		projection := mgl32.Perspective(mgl32.DegToRad(45), float32(width)/height, 0.1, 100)
		transe := projection.Mul4(camera)
		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("transe\x00")), 1, false, &transe[0])

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture1"+"\x00")), 0)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture2"+"\x00")), 1)

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, pointNum)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

// 滚轮控制场景缩放；鼠标左右上下滑动改变摄像机朝向
func Run12() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "texture2d")
	defer glfw.Terminate()

	program, _ := util.InitOpenGL(vertexShaderSource6, fragmentShaderSource6)
	vao := util.MakeVaoWithAttrib(program, vertices36, nil, []util.VertAttrib{{Name: "vPosition", Size: 3}, {Name: "vTexCoord", Size: 2}})
	pointNum := int32(len(vertices36)) / 5
	texture1 := util.MakeTexture("demo4/container.jpg")
	texture2 := util.MakeTexture("demo4/awesomeface.png")

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Enable(gl.DEPTH_TEST)

	cameraPos := mgl32.Vec3{0, 0, 3}
	cameraFront := mgl32.Vec3{0, 0, -1}
	cameraUp := mgl32.Vec3{0, 1, 0}

	// 鼠标左右上下移动，改变摄像机的朝向，实现视角移动
	//////////////////////////////////////////////////
	var firstMouse bool
	var cursorX float64 = 400
	var cursorY float64 = 300
	var yaw float64 = -90
	var pitch float64
	sensitivity := 0.05
	cursorPosCallback := func(w *glfw.Window, xpos float64, ypos float64) {
		if firstMouse {
			cursorX = xpos
			cursorY = ypos
			firstMouse = false
		}

		xoffset := sensitivity * (xpos - cursorX)
		yoffset := sensitivity * (cursorY - ypos)
		cursorX = xpos
		cursorY = ypos
		yaw += xoffset
		pitch += yoffset
		if pitch > 89 {
			pitch = 89
		}
		if pitch < -89 {
			pitch = -89
		}

		cameraFront = mgl32.Vec3{
			float32(math.Cos(float64(mgl32.DegToRad(float32(pitch)))) * math.Cos(float64(mgl32.DegToRad(float32(yaw))))),
			float32(math.Sin(float64(mgl32.DegToRad(float32(pitch))))),
			float32(math.Cos(float64(mgl32.DegToRad(float32(pitch)))) * math.Sin(float64(mgl32.DegToRad(float32(yaw))))),
		}.Normalize()
	}
	window.SetCursorPosCallback(cursorPosCallback)
	//////////////////////////////////////////////////

	// 让窗口完全捕获光标，并隐藏光标
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	//////////////////////////////////////////////////
	keyCallback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyEscape && action == glfw.Press {
			w.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		}
	}
	window.SetKeyCallback(keyCallback)
	//////////////////////////////////////////////////

	// 改变视口fov，实现滚动滚轮的缩放效果
	//////////////////////////////////////////////////
	var fov float64 = 45
	scrollCallback := func(w *glfw.Window, xoff float64, yoff float64) {
		if fov >= 1.0 && fov <= 45.0 {
			fov -= yoff
		}
		if fov <= 1.0 {
			fov = 1.0
		}
		if fov >= 45.0 {
			fov = 45.0
		}
	}
	window.SetScrollCallback(scrollCallback)
	//////////////////////////////////////////////////

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture1"+"\x00")), 0)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture2"+"\x00")), 1)

		gl.BindVertexArray(vao)

		camera := mgl32.LookAtV(cameraPos, cameraPos.Add(cameraFront), cameraUp)
		projection := mgl32.Perspective(mgl32.DegToRad(float32(fov)), float32(width)/height, 0.1, 100)
		for k, v := range more {
			model := mgl32.HomogRotate3D(mgl32.DegToRad(float32(20*k)), mgl32.Vec3{1, 0.3, 0.5})
			model = mgl32.Translate3D(v.X(), v.Y(), v.Z()).Mul4(model)
			transe := projection.Mul4(camera).Mul4(model)
			gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("transe\x00")), 1, false, &transe[0])
			gl.DrawArrays(gl.TRIANGLES, 0, pointNum)
		}

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

// 滚轮控制场景缩放；按下鼠标左键,然后左右上下移动改变摄像机朝向
func Run13() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "texture2d")
	defer glfw.Terminate()

	program, _ := util.InitOpenGL(vertexShaderSource6, fragmentShaderSource6)
	vao := util.MakeVaoWithAttrib(program, vertices36, nil, []util.VertAttrib{{Name: "vPosition", Size: 3}, {Name: "vTexCoord", Size: 2}})
	pointNum := int32(len(vertices36)) / 5
	texture1 := util.MakeTexture("demo4/container.jpg")
	texture2 := util.MakeTexture("demo4/awesomeface.png")

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Enable(gl.DEPTH_TEST)

	cameraPos := mgl32.Vec3{0, 0, 3}
	cameraFront := mgl32.Vec3{0, 0, -1}
	cameraUp := mgl32.Vec3{0, 1, 0}

	// 按下鼠标左键，左右上下移动，改变摄像机的朝向，实现视角移动
	//////////////////////////////////////////////////
	var cursorX float64 = 400
	var cursorY float64 = 300
	var yaw float64 = -90
	var pitch float64
	var leftMouseHold bool
	sensitivity := 0.05

	mouseCallback := func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if button == glfw.MouseButtonLeft {
			if action == glfw.Press {
				leftMouseHold = true
			} else {
				leftMouseHold = false
			}
		}
	}
	window.SetMouseButtonCallback(mouseCallback)

	cursorPosCallback := func(w *glfw.Window, xpos float64, ypos float64) {
		if !leftMouseHold {
			// 防止出现抖动
			cursorX = xpos
			cursorY = ypos
			return
		}

		xoffset := sensitivity * (xpos - cursorX)
		yoffset := sensitivity * (cursorY - ypos)
		cursorX = xpos
		cursorY = ypos
		yaw += xoffset
		pitch += yoffset
		if pitch > 89 {
			pitch = 89
		}
		if pitch < -89 {
			pitch = -89
		}

		cameraFront = mgl32.Vec3{
			float32(math.Cos(float64(mgl32.DegToRad(float32(pitch)))) * math.Cos(float64(mgl32.DegToRad(float32(yaw))))),
			float32(math.Sin(float64(mgl32.DegToRad(float32(pitch))))),
			float32(math.Cos(float64(mgl32.DegToRad(float32(pitch)))) * math.Sin(float64(mgl32.DegToRad(float32(yaw))))),
		}.Normalize()
	}
	window.SetCursorPosCallback(cursorPosCallback)
	//////////////////////////////////////////////////

	// 改变视口fov，实现滚动滚轮的缩放效果
	//////////////////////////////////////////////////
	var fov float64 = 45
	scrollCallback := func(w *glfw.Window, xoff float64, yoff float64) {
		if fov >= 1.0 && fov <= 45.0 {
			fov -= yoff
		}
		if fov <= 1.0 {
			fov = 1.0
		}
		if fov >= 45.0 {
			fov = 45.0
		}
	}
	window.SetScrollCallback(scrollCallback)
	//////////////////////////////////////////////////

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture1"+"\x00")), 0)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture2"+"\x00")), 1)

		gl.BindVertexArray(vao)

		camera := mgl32.LookAtV(cameraPos, cameraPos.Add(cameraFront), cameraUp)
		projection := mgl32.Perspective(mgl32.DegToRad(float32(fov)), float32(width)/height, 0.1, 100)
		for k, v := range more {
			model := mgl32.HomogRotate3D(mgl32.DegToRad(float32(20*k)), mgl32.Vec3{1, 0.3, 0.5})
			model = mgl32.Translate3D(v.X(), v.Y(), v.Z()).Mul4(model)
			transe := projection.Mul4(camera).Mul4(model)
			gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("transe\x00")), 1, false, &transe[0])
			gl.DrawArrays(gl.TRIANGLES, 0, pointNum)
		}

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

// 封装 Camera 类
func Run14() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "texture2d")
	defer glfw.Terminate()

	program, _ := util.InitOpenGL(vertexShaderSource6, fragmentShaderSource6)
	vao := util.MakeVaoWithAttrib(program, vertices36, nil, []util.VertAttrib{{Name: "vPosition", Size: 3}, {Name: "vTexCoord", Size: 2}})
	pointNum := int32(len(vertices36)) / 5
	texture1 := util.MakeTexture("demo4/container.jpg")
	texture2 := util.MakeTexture("demo4/awesomeface.png")

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Enable(gl.DEPTH_TEST)

	cameraPos := mgl32.Vec3{0, 0, 3}
	cameraFront := mgl32.Vec3{0, 0, -1}
	cameraUp := mgl32.Vec3{0, 1, 0}

	camera := util.NewCamera(cameraPos, cameraFront, cameraUp, width, height)
	camera.SetCursorPosCallback(window)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture1"+"\x00")), 0)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture2"+"\x00")), 1)

		gl.BindVertexArray(vao)

		cameraMate := camera.LookAtAndPerspective()
		for k, v := range more {
			model := mgl32.HomogRotate3D(mgl32.DegToRad(float32(20*k)), mgl32.Vec3{1, 0.3, 0.5})
			model = mgl32.Translate3D(v.X(), v.Y(), v.Z()).Mul4(model)
			transe := cameraMate.Mul4(model)
			gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("transe\x00")), 1, false, &transe[0])
			gl.DrawArrays(gl.TRIANGLES, 0, pointNum)
		}

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

// 天空盒
// 立方体贴图中，相机处于中心点，旋转鼠标即可发现自己在盒子里
func Run15() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "texture2d")
	defer glfw.Terminate()

	program, _ := util.InitOpenGL(vertexShaderSource6, fragmentShaderSource6)
	vao := util.MakeVaoWithAttrib(program, vertices36, nil, []util.VertAttrib{{Name: "vPosition", Size: 3}, {Name: "vTexCoord", Size: 2}})
	pointNum := int32(len(vertices36)) / 5
	texture1 := util.MakeTexture("demo4/container.jpg")
	texture2 := util.MakeTexture("demo4/awesomeface.png")

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Enable(gl.DEPTH_TEST)

	cameraPos := mgl32.Vec3{0, 0, 0}
	cameraFront := mgl32.Vec3{0, 0, -1}
	cameraUp := mgl32.Vec3{0, 1, 0}

	camera := util.NewCamera(cameraPos, cameraFront, cameraUp, width, height)
	camera.SetCursorPosCallback(window)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture1"+"\x00")), 0)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture2"+"\x00")), 1)

		gl.BindVertexArray(vao)
		transe := camera.LookAtAndPerspective().Mul4(mgl32.Ident4())
		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("transe\x00")), 1, false, &transe[0])
		gl.DrawArrays(gl.TRIANGLES, 0, pointNum)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

// 天空盒
// 天空盒大图下载
// 然后使用网站 https://jaxry.github.io/panorama-to-cubemap/ 或者 https://www.360toolkit.co/convert-spherical-equirectangular-tocubemap.html 来切割成6个小图
func Run16() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "textureCube")
	defer glfw.Terminate()

	program1, _ := util.InitOpenGL(vertexShaderSource8, fragmentShaderSource8)
	pointNum1 := int32(len(vertices37)) / 3
	vao1 := util.MakeVaoWithAttrib(program1, vertices37, nil, []util.VertAttrib{{Name: "vPosition", Size: 3}})
	texture := util.MakeTextureCube([]string{
		"demo4/skybox/right.jpg",
		"demo4/skybox/left.jpg",
		"demo4/skybox/top.jpg",
		"demo4/skybox/bottom.jpg",
		"demo4/skybox/front.jpg",
		"demo4/skybox/back.jpg",
	})

	program2, _ := util.MakeProgram(vertexShaderSource9, fragmentShaderSource9)
	vao2 := util.MakeVaoWithAttrib(program2, vertices36, nil, []util.VertAttrib{{Name: "vPosition", Size: 3}, {Name: "vTexCoord", Size: 2}})
	pointNum2 := int32(len(vertices36)) / 5
	texture1 := util.MakeTexture("demo4/container.jpg")

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Enable(gl.DEPTH_TEST)

	cameraPos := mgl32.Vec3{0, 0, 0} // 相机在原点
	cameraFront := mgl32.Vec3{0, 0, -1}
	cameraUp := mgl32.Vec3{0, 1, 0}

	camera := util.NewCamera(cameraPos, cameraFront, cameraUp, width, height)
	camera.SetCursorPosCallback(window)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// 天空盒
		gl.UseProgram(program1)

		gl.DepthMask(false) // 禁用深度缓存，这样天空盒就会永远被绘制在其它物体的背后
		gl.BindVertexArray(vao1)

		view := camera.LookAt()
		// 将4X4转换成3X3，去掉W分量，移除观察矩阵中的位移部分，这样天空盒就不会随着相机的
		// 移动而移动，也就不会跳出天空盒，但保留旋转变换，让玩家仍然能够环顾场景。
		view = view.Mat3().Mat4()
		projection := camera.Perspective()
		transe := projection.Mul4(view).Mul4(mgl32.Ident4())
		gl.UniformMatrix4fv(gl.GetUniformLocation(program1, gl.Str("transe\x00")), 1, false, &transe[0])

		gl.ActiveTexture(gl.TEXTURE0)
		gl.Uniform1i(gl.GetUniformLocation(program1, gl.Str("skybox"+"\x00")), 0)

		gl.BindTexture(gl.TEXTURE_CUBE_MAP, texture)
		gl.DrawArrays(gl.TRIANGLES, 0, pointNum1)

		// 其他物体
		gl.DepthMask(true) // 开启深度缓存
		gl.UseProgram(program2)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.Uniform1i(gl.GetUniformLocation(program2, gl.Str("ourTexture1"+"\x00")), 0)

		gl.BindVertexArray(vao2)
		// 此时摄像机还在中心，我们把物体向里推一点
		transe2 := camera.LookAtAndPerspective().Mul4(mgl32.Translate3D(0, 0, -3))
		gl.UniformMatrix4fv(gl.GetUniformLocation(program2, gl.Str("transe\x00")), 1, false, &transe2[0])
		gl.DrawArrays(gl.TRIANGLES, 0, pointNum2)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

// 天空盒
// 优化
func Run17() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "textureCube")
	defer glfw.Terminate()

	program1, _ := util.InitOpenGL(vertexShaderSource8, fragmentShaderSource8)
	pointNum1 := int32(len(vertices37)) / 3
	vao1 := util.MakeVaoWithAttrib(program1, vertices37, nil, []util.VertAttrib{{Name: "vPosition", Size: 3}})
	texture := util.MakeTextureCube([]string{
		"demo4/skybox/right.jpg",
		"demo4/skybox/left.jpg",
		"demo4/skybox/top.jpg",
		"demo4/skybox/bottom.jpg",
		"demo4/skybox/front.jpg",
		"demo4/skybox/back.jpg",
	})

	program2, _ := util.MakeProgram(vertexShaderSource9, fragmentShaderSource9)
	vao2 := util.MakeVaoWithAttrib(program2, vertices36, nil, []util.VertAttrib{{Name: "vPosition", Size: 3}, {Name: "vTexCoord", Size: 2}})
	pointNum2 := int32(len(vertices36)) / 5
	texture1 := util.MakeTexture("demo4/container.jpg")

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Enable(gl.DEPTH_TEST)

	cameraPos := mgl32.Vec3{0, 0, 0} // 相机在原点
	cameraFront := mgl32.Vec3{0, 0, -1}
	cameraUp := mgl32.Vec3{0, 1, 0}

	camera := util.NewCamera(cameraPos, cameraFront, cameraUp, width, height)
	camera.SetCursorPosCallback(window)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// 天空盒
		gl.UseProgram(program1)

		gl.DepthMask(false) // 禁用深度缓存，这样天空盒就会永远被绘制在其它物体的背后
		gl.BindVertexArray(vao1)

		view := camera.LookAt()
		// 移除观察矩阵中的位移部分，这样天空盒就不会随着相机的移动而移动，也就不会跳出天空盒，但保留旋转变换，让玩家仍然能够环顾场景。
		view = view.Mat3().Mat4()
		projection := camera.Perspective()
		transe := projection.Mul4(view).Mul4(mgl32.Ident4())
		gl.UniformMatrix4fv(gl.GetUniformLocation(program1, gl.Str("transe\x00")), 1, false, &transe[0])

		gl.ActiveTexture(gl.TEXTURE0)
		gl.Uniform1i(gl.GetUniformLocation(program1, gl.Str("skybox"+"\x00")), 0)

		gl.BindTexture(gl.TEXTURE_CUBE_MAP, texture)
		gl.DrawArrays(gl.TRIANGLES, 0, pointNum1)

		// 其他物体
		gl.DepthMask(true) // 开启深度缓存
		gl.UseProgram(program2)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.Uniform1i(gl.GetUniformLocation(program2, gl.Str("ourTexture1"+"\x00")), 0)

		gl.BindVertexArray(vao2)
		// 此时摄像机还在中心，我们把物体向里推一点
		transe2 := camera.LookAtAndPerspective().Mul4(mgl32.Translate3D(0, 0, -3))
		gl.UniformMatrix4fv(gl.GetUniformLocation(program2, gl.Str("transe\x00")), 1, false, &transe2[0])
		gl.DrawArrays(gl.TRIANGLES, 0, pointNum2)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

// 多张图片轮播图
// 每秒切换一次
func Run18() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "texture2d")
	defer glfw.Terminate()

	program, _ := util.InitOpenGL(vertexShaderSource, fragmentShaderSource)
	vao := util.MakeVaoWithAttrib(program, vertices, indices, []util.VertAttrib{{Name: "vPosition", Size: 3}, {Name: "vColor", Size: 3}, {Name: "vTexCoord", Size: 2}})
	pointNum := int32(len(indices))

	testureSlice := []uint32{
		util.MakeTexture("demo4/box2/front.jpg"),
		util.MakeTexture("demo4/box2/left.jpg"),
		util.MakeTexture("demo4/box2/right.jpg"),
		util.MakeTexture("demo4/box2/back.jpg"),
	}

	for !window.ShouldClose() {
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.UseProgram(program)

		gl.ActiveTexture(gl.TEXTURE0)
		t := int(glfw.GetTime()) % len(testureSlice)
		gl.BindTexture(gl.TEXTURE_2D, testureSlice[t])

		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, pointNum, gl.UNSIGNED_INT, gl.Ptr(indices))

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

// GIF图片贴图
func Run19() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "texture2d")
	defer glfw.Terminate()

	program, _ := util.InitOpenGL(vertexShaderSource, fragmentShaderSource)
	vao := util.MakeVaoWithAttrib(program, vertices, indices, []util.VertAttrib{{Name: "vPosition", Size: 3}, {Name: "vColor", Size: 3}, {Name: "vTexCoord", Size: 2}})
	pointNum := int32(len(indices))

	f, _ := os.Open("demo4/image15.gif")
	defer f.Close()

	gi, _ := gif.DecodeAll(f)

	type TextureSlice struct {
		TextureId uint32
		Delay     int
	}

	imageLen := len(gi.Image)

	textureSlice := make([]TextureSlice, imageLen)

	for k, v := range gi.Image {
		textureSlice[k].TextureId = util.MakeTextureByImage(v)
	}

	for k, d := range gi.Delay {
		textureSlice[k].Delay = d
	}

	// 更精细的参数
	var framTimeStart float64
	var frameIndex int
	var loop int

	for !window.ShouldClose() {
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.UseProgram(program)

		gl.ActiveTexture(gl.TEXTURE0)
		gltime := glfw.GetTime()

		// LoopCount 控制动画的重复播放规则
		// 0 表示无限循环
		// -1 表示只播放一次
		// 其它的播放 LoopCount+1 次
		//
		// Delay 连续的延迟时间，单位是百分之一秒，delay中数值表示展示的时间，10就表示0.1秒
		if gltime-framTimeStart > float64(textureSlice[frameIndex].Delay)/100 {
			if gi.LoopCount == 0 {
				frameIndex++
				frameIndex = frameIndex % imageLen
			} else if gi.LoopCount == -1 {
				if frameIndex < imageLen {
					frameIndex++
				}
			} else {
				if loop < gi.LoopCount+1 {
					frameIndex++
					frameIndex = frameIndex % imageLen
					if frameIndex == 0 {
						loop++
					}
				}
			}

			framTimeStart = gltime
		}

		gl.BindTexture(gl.TEXTURE_2D, textureSlice[frameIndex].TextureId)

		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, pointNum, gl.UNSIGNED_INT, gl.Ptr(indices))

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

// PNG 透明纹理
func Run20() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "texture2d")
	defer glfw.Terminate()

	program, _ := util.InitOpenGL(vertexShaderSource6, fragmentShaderSource6)
	vao := util.MakeVaoWithAttrib(program, vertices36, nil, []util.VertAttrib{{Name: "vPosition", Size: 3}, {Name: "vTexCoord", Size: 2}})
	pointNum := int32(len(vertices36)) / 5
	texture1 := util.MakeTexture("demo4/container.jpg")
	texture2 := util.MakeTexture("demo4/grass.png")

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Enable(gl.DEPTH_TEST)

	cameraPos := mgl32.Vec3{0, 0, 3}
	cameraFront := mgl32.Vec3{0, 0, -1}
	cameraUp := mgl32.Vec3{0, 1, 0}

	camera := util.NewCamera(cameraPos, cameraFront, cameraUp, width, height)
	camera.SetCursorPosCallback(window)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		transe := camera.LookAtAndPerspective()
		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("transe\x00")), 1, false, &transe[0])

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture1"+"\x00")), 0)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture2"+"\x00")), 1)

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, pointNum)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
