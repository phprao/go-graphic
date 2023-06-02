package windowcolor2d

import (
	"log"
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/phprao/go-graphic/util"
)

const (
	width  = 500
	height = 500
)

var counter int64

func Run() {
	window := util.InitGlfw(width, height, "2d color")
	defer glfw.Terminate()

	gl.Init()

	KeyPressAction(window)
	glfw.SwapInterval(100)

	// n := 0
	// gl.ClearColor(1.0, 0.0, 0.0, 1.0)

	for !window.ShouldClose() {
		// gl.Clear(gl.COLOR_BUFFER_BIT)
		glfw.PollEvents()
		// n++
		// if n%50 == 0 {
		// 	gl.Clear(gl.COLOR_BUFFER_BIT)
		// }
		// if n >= 100 {
		// 	n = 0
		// }
		window.SwapBuffers()
		// log.Println("ok")
	}
}

func KeyPressAction(window *glfw.Window) {
	keyCallback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}

		if window.GetKey(glfw.KeyR) == glfw.Press {
			log.Println("R")
			gl.ClearColor(1.0, 0.0, 0.0, 1.0)
			gl.Clear(gl.COLOR_BUFFER_BIT)
		}

		if window.GetKey(glfw.KeyG) == glfw.Press {
			log.Println("G")
			gl.ClearColor(0.0, 1.0, 0.0, 1.0)
			gl.Clear(gl.COLOR_BUFFER_BIT)
		}

		if window.GetKey(glfw.KeyB) == glfw.Press {
			log.Println("B")
			gl.ClearColor(0.0, 0.0, 1.0, 1.0)
			gl.Clear(gl.COLOR_BUFFER_BIT)
		}

		if window.GetKey(glfw.KeyA) == glfw.Press {
			log.Println("A")
			counter = (counter + 1) % 100
			var u float64 = float64(counter) / float64(99.0)
			gl.ClearColor(0.5+float32(math.Cos(float64(2*3.14*u)))/2.0, 1.0, 0.5, 1.0)
		}

		if window.GetKey(glfw.KeyK) == glfw.Press {
			monitor := glfw.GetPrimaryMonitor()
			videoMode := monitor.GetVideoMode()
			window.SetMonitor(monitor, 0, 0, videoMode.Width, videoMode.Height, videoMode.RefreshRate)
		}
		if window.GetKey(glfw.KeyM) == glfw.Press {
			monitor := glfw.GetPrimaryMonitor()
			videoMode := monitor.GetVideoMode()
			window.SetMonitor(nil, 100, 100, 500, 500, videoMode.RefreshRate)
		}
	}

	window.SetKeyCallback(keyCallback)
}
