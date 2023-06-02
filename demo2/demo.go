package demo2

// 教程：https://linux.cn/article-8937-1.html

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/phprao/go-graphic/util"
)

const (
	width   = 500
	height  = 500
	rows    = 10
	columns = 10

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
	square = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
		-0.5, 0.5, 0,
		0.5, 0.5, 0,
		0.5, -0.5, 0,
	}
)

type cell struct {
	drawable uint32
	x        int
	y        int
}

func (c *cell) draw() {
	gl.BindVertexArray(c.drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
}

func Run() {
	runtime.LockOSThread()

	window := util.InitGlfw(width, height, "Conway's Game of Life")
	defer glfw.Terminate()

	var x0, y0, x1, x2, y1, y2 float64
	mouseCallback := func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		log.Printf("button:%d, action:%d, mod:%d\n", button, action, mod)

		if button == glfw.MouseButtonLeft && action == glfw.Press {
			x1, y1 = x0, y0
			log.Printf("x1:%f, y1:%f", x1, y1)
		}

		if button == glfw.MouseButtonLeft && action == glfw.Release {
			x2, y2 = x0, y0
			log.Printf("x2:%f, y2:%f", x2, y2)
			log.Printf("x move:%f, y move:%f", x2-x1, y2-y1)
		}
	}
	window.SetMouseButtonCallback(mouseCallback)

	cursorPosCallback := func(w *glfw.Window, xpos float64, ypos float64) {
		x0 = xpos
		y0 = ypos
	}
	window.SetCursorPosCallback(cursorPosCallback)

	scrollCallback := func(w *glfw.Window, xoff float64, yoff float64) {
		log.Printf("xoff:%f, yoff:%f", x2, y2)
	}
	window.SetScrollCallback(scrollCallback)

	dropCallback := func(w *glfw.Window, names []string) {
		// names:[D:\dev\php\magook\trunk\server\go-graphic\demo5\square.png]
		log.Printf("names:%v", names)
	}
	window.SetDropCallback(dropCallback)

	charCallback := func(w *glfw.Window, char rune) {
		log.Printf("char:%s", string(char))
	}
	window.SetCharCallback(charCallback)

	sizeCallback := func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	}
	window.SetSizeCallback(sizeCallback)

	program := util.InitOpenGL(vertexShaderSource, fragmentShaderSource)

	cells := makeCells()

	for !window.ShouldClose() {
		draw(cells, window, program)
	}
}

func makeCells() [][]*cell {
	cells := make([][]*cell, rows, columns)
	for x := 0; x < rows; x++ {
		for y := 0; y < columns; y++ {
			c := newCell(x, y)
			cells[x] = append(cells[x], c)
		}
	}
	return cells
}

func newCell(x, y int) *cell {
	points := make([]float32, len(square), len(square))
	copy(points, square)
	for i := 0; i < len(points); i++ {
		// var position float32
		// var size float32
		// switch i % 3 {
		// case 0: // 操作X坐标
		// 	size = 1.0 / float32(columns)
		// 	position = float32(x) * size
		// case 1: // 操作Y坐标
		// 	size = 1.0 / float32(rows)
		// 	position = float32(y) * size
		// default:
		// 	continue
		// }
		// if points[i] < 0 {
		// 	points[i] = (position * 2) - 1
		// } else {
		// 	points[i] = ((position + size) * 2) - 1
		// }

		var position float32
		var size float32
		switch i % 3 {
		case 0:
			// 操作X坐标
			// 取值范围是-1到1，因此长度为2，size为单个方格的x或y的长度
			// 由于x和y取值是大于等于0的，因此不妨现在[0,2]区间上来排，然后再减1，挪到[-1,1]区间
			// position就是偏移量
			size = 2.0 / float32(columns)
			position = float32(x) * size
		case 1:
			// 操作Y坐标
			size = 2.0 / float32(rows)
			position = float32(y) * size

		default:
			continue
		}
		if points[i] < 0 {
			points[i] = position - 1
		} else {
			points[i] = position + size - 1
		}
	}
	// fmt.Printf("x=%d, y=%d, points=%v\n", x, y, points)
	return &cell{
		drawable: util.MakeVao(points),
		x:        x,
		y:        y,
	}
}

func draw(cells [][]*cell, window *glfw.Window, prog uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(prog)

	for x := range cells {
		for _, c := range cells[x] {
			c.draw()
		}
	}

	glfw.PollEvents()
	window.SwapBuffers()
}
