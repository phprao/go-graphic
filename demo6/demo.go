package demo6

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
)

func init() {
	runtime.LockOSThread()
}

func Run() {
	window := util.InitGlfw(width, height, "Conway's Game of Life")
	defer glfw.Terminate()

	initOpenGL()

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_BYTE, gl.Ptr(0))

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

// initOpenGL 初始化 OpenGL 并且返回一个初始化了的程序。
func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	return
}
