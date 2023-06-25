package circle

// 圆相关的

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
	height = 800
	PI     = 3.14159265358979323846

	vertexShaderSource1 = `
        #version 410

        in vec2 vPosition;

        void main() {
            gl_Position = vec4(vPosition, 1.0, 1.0);
        }
    ` + "\x00"
	fragmentShaderSource1 = `
        #version 410
        
		out vec4 frag_colour;

        void main() {
			frag_colour = vec4(1, 1, 1, 1);
        }
    ` + "\x00"

	vertexShaderSource2 = `
        #version 410

        in vec3 vPosition;

		uniform mat4 transe;

        void main() {
            gl_Position = transe * vec4(vPosition, 1.0);
        }
    ` + "\x00"
	fragmentShaderSource2 = `
        #version 410
        
		out vec4 frag_colour;

        void main() {
			frag_colour = vec4(1.0, 0.635, 0.345, 1.0);
        }
    ` + "\x00"

	vertexShaderSource3 = `
        #version 410

        in vec3 vPosition;
		in vec2 texCoord;

		out vec2 tc;

		uniform mat4 transe;

        void main() {
            gl_Position = transe * vec4(vPosition, 1.0);
			tc = vec2(texCoord.x, texCoord.y);
        }
    ` + "\x00"
	fragmentShaderSource3 = `
        #version 410
        
		out vec4 frag_colour;
		
		in vec2 tc;

		uniform sampler2D samp;

        void main() {
			frag_colour = texture(samp, tc);
        }
    ` + "\x00"
)

// 圆形
func Run1() {
	pointCount := 360
	vertices := make([]float32, (pointCount+1)*2)
	deg := 360.0 / float32(pointCount)

	vertices[0] = 0
	vertices[1] = 0
	for i := 1; i <= pointCount; i++ {
		index := i * 2
		radians := float64(mgl32.DegToRad(float32(i-1) * deg))
		vertices[index] = float32(math.Cos(radians)) / 4
		vertices[index+1] = float32(math.Sin(radians)) / 4
	}

	indices := make([]uint32, pointCount*3)
	for i := 1; i <= pointCount; i++ {
		index := (i - 1) * 3
		indices[index] = uint32(i)
		indices[index+1] = 0
		indices[index+2] = uint32(i + 1)
		if i+1 > pointCount {
			indices[index+2] = 1
		}
	}

	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "circle")
	defer glfw.Terminate()

	program, _ := util.InitOpenGL(vertexShaderSource1, fragmentShaderSource1)
	vao := util.MakeVaoWithAttrib(program, vertices, indices, []util.VertAttrib{{Name: "vPosition", Size: 2}})
	pointNum := int32(len(indices))

	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, pointNum, gl.UNSIGNED_INT, gl.Ptr(indices))

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

// 圆球
//
// 原理：https://www.cnblogs.com/8335IT/p/16290888.html
// 球面上的一个点P，与球心的连线L，L与Y轴正向的夹角为 α，L投影到 XOZ平面上L1，L1与X正向的夹角为 β，球面半径为 R
// 于是
//
//	y = R * cosα;
//	x = R * sinα * cosβ;
//	z = R * sinα * sinβ;
func Run2() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "circle")
	defer glfw.Terminate()

	// 将 0~π 的 α 角分为 Y_SEGMENTS 份，得到 Y_SEGMENTS + 1 个点
	Y_SEGMENTS := 10
	// 将 0~2π 的 β 角分为 X_SEGMENTS 份，得到 X_SEGMENTS + 1 个点
	X_SEGMENTS := 10

	vertices := make([]float32, (Y_SEGMENTS+1)*(X_SEGMENTS+1)*3)

	index := 0
	// 定义所有的顶点
	for y := 0; y <= Y_SEGMENTS; y++ {

		ySegment := float64(y) / float64(Y_SEGMENTS)
		α := ySegment * PI
		yPos := math.Cos(α)

		for x := 0; x <= X_SEGMENTS; x++ {
			xSegment := float64(x) / float64(X_SEGMENTS)
			β := xSegment * 2.0 * PI

			xPos := math.Cos(β) * math.Sin(α)
			zPos := math.Sin(β) * math.Sin(α)

			vertices[index] = float32(xPos)
			vertices[index+1] = float32(yPos)
			vertices[index+2] = float32(zPos)

			index += 3
		}
	}

	// 定义顶点索引，规定如何画三角形
	// 一个顶点会连接着6个三角形，需要定位它在每个三角形中的索引，三角形内部以顺时针来数点。
	indices := make([]uint32, Y_SEGMENTS*X_SEGMENTS*6)
	indix := 0
	for i := 0; i < Y_SEGMENTS; i++ {
		for j := 0; j < X_SEGMENTS; j++ {
			indices[indix] = uint32(i*(X_SEGMENTS+1) + j)
			indices[indix+1] = uint32((i+1)*(X_SEGMENTS+1) + j)
			indices[indix+2] = uint32((i+1)*(X_SEGMENTS+1) + j + 1)
			indices[indix+3] = uint32(i*(X_SEGMENTS+1) + j)
			indices[indix+4] = uint32((i+1)*(X_SEGMENTS+1) + j + 1)
			indices[indix+5] = uint32(i*(X_SEGMENTS+1) + j + 1)

			indix += 6
		}
	}

	program, _ := util.InitOpenGL(vertexShaderSource2, fragmentShaderSource2)
	vao := util.MakeVaoWithAttrib(program, vertices, indices, []util.VertAttrib{{Name: "vPosition", Size: 3}})
	pointNum := int32(len(indices))

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)

	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE) // 线框模式
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE) // 背面剔除，否则背面的线也能看到，杂乱无章
	gl.CullFace(gl.BACK)

	cameraPos := mgl32.Vec3{0, 0, 3}
	cameraFront := mgl32.Vec3{0, 0, -1}
	cameraUp := mgl32.Vec3{0, 1, 0}

	camera := util.NewCamera(cameraPos, cameraFront, cameraUp, width, height)
	camera.SetCursorPosCallback(window)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		transe2 := camera.LookAtAndPerspective().Mul4(mgl32.Translate3D(0, 0, -3))
		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("transe\x00")), 1, false, &transe2[0])

		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, pointNum, gl.UNSIGNED_INT, gl.Ptr(indices))

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

// 球面贴图
func Run3() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "circle")
	defer glfw.Terminate()

	// 将 0~π 的 α 角分为 Y_SEGMENTS 份，得到 Y_SEGMENTS + 1 个点
	Y_SEGMENTS := 50
	// 将 0~2π 的 β 角分为 X_SEGMENTS 份，得到 X_SEGMENTS + 1 个点
	X_SEGMENTS := 50

	vertices := make([]float32, (Y_SEGMENTS+1)*(X_SEGMENTS+1)*5)

	index := 0
	// 定义所有的顶点
	for y := 0; y <= Y_SEGMENTS; y++ {

		ySegment := float64(y) / float64(Y_SEGMENTS)
		α := ySegment * PI
		yPos := math.Cos(α)

		for x := 0; x <= X_SEGMENTS; x++ {
			xSegment := float64(x) / float64(X_SEGMENTS)
			β := xSegment * 2.0 * PI

			xPos := math.Cos(β) * math.Sin(α)
			zPos := math.Sin(β) * math.Sin(α)

			vertices[index] = float32(xPos)
			vertices[index+1] = float32(yPos)
			vertices[index+2] = float32(zPos)

			// 纹理坐标，按等间距划分
			vertices[index+3] = float32(xSegment)
			vertices[index+4] = float32(ySegment)

			index += 5
		}
	}

	// 定义顶点索引，规定如何画三角形
	// 一个顶点会连接着6个三角形，需要定位它在每个三角形中的索引，三角形内部以顺时针来数点。
	indices := make([]uint32, Y_SEGMENTS*X_SEGMENTS*6)
	indix := 0
	for i := 0; i < Y_SEGMENTS; i++ {
		for j := 0; j < X_SEGMENTS; j++ {
			indices[indix] = uint32(i*(X_SEGMENTS+1) + j)
			indices[indix+1] = uint32((i+1)*(X_SEGMENTS+1) + j)
			indices[indix+2] = uint32((i+1)*(X_SEGMENTS+1) + j + 1)
			indices[indix+3] = uint32(i*(X_SEGMENTS+1) + j)
			indices[indix+4] = uint32((i+1)*(X_SEGMENTS+1) + j + 1)
			indices[indix+5] = uint32(i*(X_SEGMENTS+1) + j + 1)

			indix += 6
		}
	}

	program, _ := util.InitOpenGL(vertexShaderSource3, fragmentShaderSource3)
	vao := util.MakeVaoWithAttrib(program, vertices, indices, []util.VertAttrib{{Name: "vPosition", Size: 3}, {Name: "texCoord", Size: 2}})
	pointNum := int32(len(indices))

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)

	// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE) // 线框模式
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE) // 背面剔除，否则背面的线也能看到，杂乱无章
	gl.CullFace(gl.BACK)

	cameraPos := mgl32.Vec3{0, 0, 3}
	cameraFront := mgl32.Vec3{0, 0, -1}
	cameraUp := mgl32.Vec3{0, 1, 0}

	camera := util.NewCamera(cameraPos, cameraFront, cameraUp, width, height)
	camera.SetCursorPosCallback(window)

	texture := util.MakeTexture("circle/earth2.jpg")

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)

		transe2 := camera.LookAtAndPerspective().Mul4(mgl32.Translate3D(0, 0, -3))
		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("transe\x00")), 1, false, &transe2[0])

		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, pointNum, gl.UNSIGNED_INT, gl.Ptr(indices))

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
