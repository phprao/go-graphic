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

	vertices36 = []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,

		-0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
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

// 立方体不动，摄像机受到按键WSAD上下左右移动
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
