package demo1

// 教程：https://linux.cn/article-8933-1.html

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/phprao/go-graphic/util"
)

const (
	width  = 500
	height = 500

	// 顶点着色器
	vertexShaderSource = `
        #version 410
        in vec3 vp;
        void main() {
            gl_Position = vec4(vp, 1.0);
        }
    ` + "\x00"

	// 片元着色器
	fragmentShaderSource = `
        #version 410
        out vec4 frag_colour;
        void main() {
            frag_colour = vec4(1, 1, 1, 1);
        }
    ` + "\x00"
)

var (
	// X,Y,Z
	// 窗口中心点为原点，向右为X正，上为Y正，取值 -1到1
	triangle = []float32{
		0, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
	}

	// 四边形
	square = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
		-0.5, 0.5, 0,
		0.5, 0.5, 0,
		0.5, -0.5, 0,
	}

	// 四边形2
	square2 = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
		0.5, 0.5, 0,
	}

	// 索引数据
	indexs = []uint32{
		0, 1, 2,
		0, 2, 3,
	}
)

func Run() {
	runtime.LockOSThread()

	window := util.InitGlfw(width, height, "Conway's Game of Life")

	defer glfw.Terminate()

	KeyPressAction(window)

	program := util.InitOpenGL(vertexShaderSource, fragmentShaderSource)

	vao := util.MakeVao(square)
	pointNum := int32(len(square))

	glfw.SwapInterval(1)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.UseProgram(program)

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, pointNum)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func Run2() {
	runtime.LockOSThread()

	window := util.InitGlfw(width, height, "Conway's Game of Life")

	defer glfw.Terminate()

	KeyPressAction(window)

	program := util.InitOpenGL(vertexShaderSource, fragmentShaderSource)

	vao := util.MakeVaoWithEbo(square2, indexs)
	pointNum := int32(len(indexs))

	glfw.SwapInterval(1)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.UseProgram(program)

		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, pointNum, gl.UNSIGNED_INT, gl.Ptr(indexs))

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func KeyPressAction(window *glfw.Window) {
	// action参数表示这个按键是被按下还是释放，按下的时候会触发action=1，如果不放会一直触发action=2，放开的时候会触发action=0事件
	// mods表示是否有Ctrl、Shift、Alt、Super四个按钮的操作，1-shift,2-ctrl,4-alt，8-win
	keyCallback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		log.Printf("key:%d, scancode:%d, action:%d, mods:%v\n", key, scancode, action, mods)
		// 如果按下了ESC键就关闭窗口
		if key == glfw.KeyEscape && action == glfw.Press {
			window.SetShouldClose(true)
		}
	}
	window.SetKeyCallback(keyCallback)
}
