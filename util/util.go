package util

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// initGlfw 初始化 glfw 并且返回一个可用的窗口。
func InitGlfw(width, height int, name string) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Visible, glfw.False)
	window, err := glfw.CreateWindow(width, height, name, nil, nil)
	if err != nil {
		panic(err)
	}

	// center
	sw := glfw.GetPrimaryMonitor().GetVideoMode().Width
	sh := glfw.GetPrimaryMonitor().GetVideoMode().Height
	window.SetPos((sw-width)/2, (sh-height)/2)
	window.Show()

	window.MakeContextCurrent()

	return window
}

func InitGlfw2(width, height int, name string) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	monitor := glfw.GetPrimaryMonitor()
	videoMode := monitor.GetVideoMode()
	glfw.WindowHint(glfw.RedBits, videoMode.RedBits)
	glfw.WindowHint(glfw.GreenBits, videoMode.GreenBits)
	glfw.WindowHint(glfw.BlueBits, videoMode.BlueBits)
	glfw.WindowHint(glfw.RefreshRate, videoMode.RefreshRate)
	window, err := glfw.CreateWindow(videoMode.Width, videoMode.Height, name, monitor, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

func CompileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func MakeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func MakeVaoWithEbo(points []float32, indexs []uint32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indexs), gl.Ptr(indexs), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	return vao
}

func MakeVaoWithEboAndAttrib(points []float32, indexs []uint32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indexs), gl.Ptr(indexs), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Position attribute
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 8*4, 0)
	gl.EnableVertexAttribArray(0)
	// Color attribute
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 8*4, 3*4)
	gl.EnableVertexAttribArray(1)
	// TexCoord attribute
	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, 8*4, 6*4)
	gl.EnableVertexAttribArray(2)

	return vao
}

func MakeTexture(filepath string) uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	MakeTextureTexParameteri()

	imgFile2, _ := os.Open(filepath)
	defer imgFile2.Close()
	img2, _, _ := image.Decode(imgFile2)
	rgba2 := image.NewRGBA(img2.Bounds())
	draw.Draw(rgba2, rgba2.Bounds(), img2, image.Point{0, 0}, draw.Src)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(rgba2.Rect.Size().X), int32(rgba2.Rect.Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba2.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return texture
}

func MakeTextureTexParameteri() {
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
}

func InitOpenGL(vertexShaderSource, fragmentShaderSource string) (program uint32) {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	program = gl.CreateProgram()

	if vertexShaderSource != "" {
		vertexShader, err := CompileShader(vertexShaderSource, gl.VERTEX_SHADER)
		if err != nil {
			panic(err)
		}
		gl.AttachShader(program, vertexShader)
	}

	if fragmentShaderSource != "" {
		fragmentShader, err := CompileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
		if err != nil {
			panic(err)
		}
		gl.AttachShader(program, fragmentShader)
	}

	gl.LinkProgram(program)
	return program
}
