package light

// 光照

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

	in vec3 aPos;
	in vec3 aNormal;

	out vec3 FragPos;
	out vec3 Normal;

	uniform mat4 model;
	uniform mat4 view;
	uniform mat4 projection;
	uniform mat3 normalModel;

	void main() {
		Normal = normalModel * aNormal;
		FragPos = vec3(model * vec4(aPos, 1.0));
		gl_Position = projection * view * model * vec4(aPos, 1.0);
	}
	` + "\x00"

	fragmentShaderSource = `
	#version 410
    
	out vec4 frag_colour;

	in vec3 Normal;
	in vec3 FragPos;

	uniform vec3 lightPos;
	uniform vec3 viewPos;
	uniform vec3 lightColor;
	uniform vec3 objectColor;

	void main() {
		// 环境光
		float ambientStrength = 0.1;
		vec3 ambient = ambientStrength * lightColor;

		// 漫反射光
		vec3 norm = normalize(Normal);
		vec3 lightDir = normalize(lightPos - FragPos);
		float diff = max(dot(norm, lightDir), 0.0);
		vec3 diffuse = diff * lightColor;

		// 镜面光
		float specularStrength = 0.5;
		vec3 viewDir = normalize(viewPos - FragPos);
		vec3 reflectDir = reflect(-lightDir, norm);
		float spec = pow(max(dot(viewDir, reflectDir), 0.0), 32);
		vec3 specular = specularStrength * spec * lightColor;

		frag_colour = vec4((ambient + diffuse + specular) * objectColor, 1.0);
	}
	` + "\x00"

	vertexShaderSource2 = `
	#version 410

	in vec3 aPos;

	uniform mat4 model;
	uniform mat4 view;
	uniform mat4 projection;

	void main() {
		gl_Position = projection * view * model * vec4(aPos, 1.0);
	}
	` + "\x00"

	fragmentShaderSource2 = `
	#version 410
    
	out vec4 frag_colour;

	void main() {
		frag_colour = vec4(1.0);
	}
	` + "\x00"
)

var (
	vertices = []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0,
		0.5, -0.5, -0.5, 0.0, 0.0, -1.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0,
		-0.5, 0.5, -0.5, 0.0, 0.0, -1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0,

		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0,

		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0,
		-0.5, 0.5, -0.5, -1.0, 0.0, 0.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0,
		-0.5, -0.5, 0.5, -1.0, 0.0, 0.0,
		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0,
		0.5, -0.5, -0.5, 0.0, -1.0, 0.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, -1.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0,

		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
		0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
	}
)

func Run() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "light")
	defer glfw.Terminate()

	pointNum := int32(len(vertices)) / 5

	program1, _ := util.InitOpenGL(vertexShaderSource, fragmentShaderSource)
	vao1 := util.MakeVaoWithAttrib(program1, vertices, nil, []util.VertAttrib{{Name: "aPos", Size: 3}, {Name: "aNormal", Size: 3}})
	program2, _ := util.MakeProgram(vertexShaderSource2, fragmentShaderSource2)
	vao2 := util.MakeVaoWithAttrib(program2, vertices, nil, []util.VertAttrib{{Name: "aPos", Size: 3}, {Name: "aNormal", Size: 3}})

	gl.ClearColor(0.1, 0.1, 0.1, 1.0)
	gl.Enable(gl.DEPTH_TEST)

	lightPos := mgl32.Vec3{2, 0, 0}
	lightColor := mgl32.Vec3{1, 1, 1}
	objectColor := mgl32.Vec3{1, 0.5, 0.31}

	cameraPos := mgl32.Vec3{0, 0, 3}
	cameraFront := mgl32.Vec3{0, 0, -1}
	cameraUp := mgl32.Vec3{0, 1, 0}

	camera := util.NewCamera(cameraPos, cameraFront, cameraUp, width, height)
	camera.SetCursorPosCallback(window)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		lightPos = mgl32.Vec3{2, float32(math.Sin(glfw.GetTime())), 0}

		// 画箱子
		gl.UseProgram(program1)

		gl.Uniform3fv(gl.GetUniformLocation(program1, gl.Str("lightColor\x00")), 1, &lightColor[0])
		gl.Uniform3fv(gl.GetUniformLocation(program1, gl.Str("lightPos\x00")), 1, &lightPos[0])
		gl.Uniform3fv(gl.GetUniformLocation(program1, gl.Str("objectColor\x00")), 1, &objectColor[0])
		gl.Uniform3fv(gl.GetUniformLocation(program1, gl.Str("viewPos\x00")), 1, &camera.CameraPos[0])
		view := camera.LookAt()
		projection := camera.Perspective()

		model1 := mgl32.Ident4()
		normalModel := model1.Inv().Transpose().Mat3()
		gl.UniformMatrix4fv(gl.GetUniformLocation(program1, gl.Str("model\x00")), 1, false, &model1[0])
		gl.UniformMatrix4fv(gl.GetUniformLocation(program1, gl.Str("view\x00")), 1, false, &view[0])
		gl.UniformMatrix4fv(gl.GetUniformLocation(program1, gl.Str("projection\x00")), 1, false, &projection[0])
		gl.UniformMatrix3fv(gl.GetUniformLocation(program1, gl.Str("normalModel\x00")), 1, false, &normalModel[0])
		gl.BindVertexArray(vao1)
		gl.DrawArrays(gl.TRIANGLES, 0, pointNum)

		// 画光源
		gl.UseProgram(program2)
		model2 := mgl32.Translate3D(lightPos.X(), lightPos.Y(), lightPos.Z()).Mul4(mgl32.Scale3D(0.2, 0.2, 0.2))
		gl.UniformMatrix4fv(gl.GetUniformLocation(program2, gl.Str("model\x00")), 1, false, &model2[0])
		gl.UniformMatrix4fv(gl.GetUniformLocation(program2, gl.Str("view\x00")), 1, false, &view[0])
		gl.UniformMatrix4fv(gl.GetUniformLocation(program2, gl.Str("projection\x00")), 1, false, &projection[0])
		gl.BindVertexArray(vao2)
		gl.DrawArrays(gl.TRIANGLES, 0, pointNum)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
