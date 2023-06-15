package util

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
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

type VertAttrib struct {
	Name string
	Size int32
}

func MakeVaoWithAttrib(program uint32, points []float32, indexs []uint32, vertAttribSlice []VertAttrib) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	if indexs != nil {
		var ebo uint32
		gl.GenBuffers(1, &ebo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indexs), gl.Ptr(indexs), gl.STATIC_DRAW)
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var stride int32
	for _, v := range vertAttribSlice {
		stride += v.Size
	}
	var offset uintptr
	for _, v := range vertAttribSlice {
		vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str(v.Name+"\x00")))
		gl.EnableVertexAttribArray(vertAttrib)
		gl.VertexAttribPointerWithOffset(vertAttrib, v.Size, gl.FLOAT, false, stride*4, offset*4)
		offset += uintptr(v.Size)
	}

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

func InitOpenGL(vertexShaderSource, fragmentShaderSource string) (program uint32, err error) {
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

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link shader program: %v", log)
	}

	return program, nil
}

func MakeProgram(vertexShaderSource, fragmentShaderSource string) (program uint32, err error) {
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

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link shader program: %v", log)
	}

	return program, nil
}

type Camera struct {
	CameraPos    mgl32.Vec3
	CameraFront  mgl32.Vec3
	CameraUp     mgl32.Vec3
	Fov          float64
	WindowWidth  int
	WindowHeight int
}

// cameraFront 为相机的朝向.
// cameraPos 为相机的位置.
// windowWidth 窗口宽度.
// windowHeight 窗口高度.
func NewCamera(cameraPos mgl32.Vec3, cameraFront mgl32.Vec3, cameraUp mgl32.Vec3, windowWidth int, windowHeight int) *Camera {
	return &Camera{cameraPos, cameraFront, cameraUp, 45, windowWidth, windowHeight}
}

func (c *Camera) LookAt() mgl32.Mat4 {
	return mgl32.LookAtV(c.CameraPos, c.CameraPos.Add(c.CameraFront), c.CameraUp)
}

func (c *Camera) Perspective() mgl32.Mat4 {
	return mgl32.Perspective(mgl32.DegToRad(float32(c.Fov)), float32(c.WindowWidth)/float32(c.WindowHeight), 0.1, 100)
}

func (c *Camera) SetCursorPosCallback(window *glfw.Window) {
	// WSAD按键【平移相机位置】
	keyCallback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		cameraSpeed := float32(0.05)
		moveUp := mgl32.Vec3{0, 1, 0}
		// Z轴前进（向里）
		if window.GetKey(glfw.KeyW) == glfw.Press {
			c.CameraPos = c.CameraPos.Sub(c.CameraFront.Mul(cameraSpeed))
		}
		// Z轴后退（向外）
		if window.GetKey(glfw.KeyS) == glfw.Press {
			c.CameraPos = c.CameraPos.Add(c.CameraFront.Mul(cameraSpeed))
		}
		// X轴向左
		if window.GetKey(glfw.KeyA) == glfw.Press {
			c.CameraPos = c.CameraPos.Add(c.CameraFront.Cross(c.CameraUp).Normalize().Mul(cameraSpeed))
		}
		// X轴向右
		if window.GetKey(glfw.KeyD) == glfw.Press {
			c.CameraPos = c.CameraPos.Sub(c.CameraFront.Cross(c.CameraUp).Normalize().Mul(cameraSpeed))
		}
		// Y轴向上
		if window.GetKey(glfw.KeyQ) == glfw.Press {
			c.CameraPos = c.CameraPos.Add(moveUp.Mul(cameraSpeed))
		}
		// Y轴向下
		if window.GetKey(glfw.KeyE) == glfw.Press {
			c.CameraPos = c.CameraPos.Sub(moveUp.Mul(cameraSpeed))
		}
	}

	window.SetKeyCallback(keyCallback)

	// 鼠标滚轮实现【缩放】
	scrollCallback := func(w *glfw.Window, xoff float64, yoff float64) {
		if c.Fov >= 1.0 && c.Fov <= 45.0 {
			c.Fov -= yoff
		}
		if c.Fov <= 1.0 {
			c.Fov = 1.0
		}
		if c.Fov >= 45.0 {
			c.Fov = 45.0
		}
	}
	window.SetScrollCallback(scrollCallback)

	// 按下鼠标左键，移动鼠标来改变【相机朝向】
	cursorX := float64(c.WindowWidth / 2)
	cursorY := float64(c.WindowHeight / 2)
	var yaw float64 = -90
	var pitch float64
	var leftMouseHold bool
	var sensitivity float64 = 0.05
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

		c.CameraFront = mgl32.Vec3{
			float32(math.Cos(float64(mgl32.DegToRad(float32(pitch)))) * math.Cos(float64(mgl32.DegToRad(float32(yaw))))),
			float32(math.Sin(float64(mgl32.DegToRad(float32(pitch))))),
			float32(math.Cos(float64(mgl32.DegToRad(float32(pitch)))) * math.Sin(float64(mgl32.DegToRad(float32(yaw))))),
		}.Normalize()
	}
	window.SetCursorPosCallback(cursorPosCallback)
}
