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

	vertexShaderSourceLight = `
	#version 410

	in vec3 aPos;

	uniform mat4 model;
	uniform mat4 view;
	uniform mat4 projection;

	void main() {
		gl_Position = projection * view * model * vec4(aPos, 1.0);
	}
	` + "\x00"

	fragmentShaderSourceLight = `
	#version 410
    
	out vec4 frag_colour;

	void main() {
		frag_colour = vec4(1.0);
	}
	` + "\x00"

	vertexShaderSource2 = `
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

	fragmentShaderSource2 = `
	#version 410
    
	out vec4 frag_colour;

	in vec3 Normal;
	in vec3 FragPos;

	uniform vec3 viewPos;

	struct Material {
		vec3 ambient;
		vec3 diffuse;
		vec3 specular;
		float shininess;
	}; 

	uniform Material material;

	struct Light {
		vec3 position;
	
		vec3 ambient;
		vec3 diffuse;
		vec3 specular;
	};
	
	uniform Light light;

	void main() {
		// 环境光
		vec3 ambient = light.ambient * material.ambient;

		// 漫反射光
		vec3 norm = normalize(Normal);
		vec3 lightDir = normalize(light.position - FragPos);
		float diff = max(dot(norm, lightDir), 0.0);
		vec3 diffuse = light.diffuse * diff * material.diffuse;

		// 镜面光
		vec3 viewDir = normalize(viewPos - FragPos);
		vec3 reflectDir = reflect(-lightDir, norm);
		float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shininess);
		vec3 specular = light.specular * spec * material.specular;

		frag_colour = vec4(ambient + diffuse + specular, 1.0);
	}
	` + "\x00"

	vertexShaderSource3 = `
	#version 410

	in vec3 aPos;
	in vec3 aNormal;
	in vec2 aTexCoord;

	out vec3 FragPos;
	out vec3 Normal;
	out vec2 TexCoords;

	uniform mat4 model;
	uniform mat4 view;
	uniform mat4 projection;
	uniform mat3 normalModel;

	void main() {
		Normal = normalModel * aNormal;
		FragPos = vec3(model * vec4(aPos, 1.0));
		gl_Position = projection * view * model * vec4(aPos, 1.0);
		TexCoords = vec2(aTexCoord.x, 1.0-aTexCoord.y);
	}
	` + "\x00"

	fragmentShaderSource3 = `
	#version 410
    
	out vec4 frag_colour;

	in vec3 Normal;
	in vec3 FragPos;
	in vec2 TexCoords;

	uniform vec3 viewPos;

	struct Material {
		sampler2D diffuse;
		vec3 specular;
		float shininess;
	}; 

	uniform Material material;

	struct Light {
		vec3 position;
	
		vec3 ambient;
		vec3 diffuse;
		vec3 specular;
	};
	
	uniform Light light;

	void main() {
		// 环境光
		vec3 ambient = light.ambient * vec3(texture(material.diffuse, TexCoords));

		// 漫反射光
		vec3 norm = normalize(Normal);
		vec3 lightDir = normalize(light.position - FragPos);
		float diff = max(dot(norm, lightDir), 0.0);
		vec3 diffuse = light.diffuse * diff * vec3(texture(material.diffuse, TexCoords));

		// 镜面光
		vec3 viewDir = normalize(viewPos - FragPos);
		vec3 reflectDir = reflect(-lightDir, norm);
		float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shininess);
		vec3 specular = light.specular * spec * material.specular;

		frag_colour = vec4(ambient + diffuse + specular, 1.0);
	}
	` + "\x00"

	vertexShaderSource4 = `
	#version 410

	in vec3 aPos;
	in vec3 aNormal;
	in vec2 aTexCoord;

	out vec3 FragPos;
	out vec3 Normal;
	out vec2 TexCoords;

	uniform mat4 model;
	uniform mat4 view;
	uniform mat4 projection;
	uniform mat3 normalModel;

	void main() {
		Normal = normalModel * aNormal;
		FragPos = vec3(model * vec4(aPos, 1.0));
		gl_Position = projection * view * model * vec4(aPos, 1.0);
		TexCoords = vec2(aTexCoord.x, 1.0-aTexCoord.y);
	}
	` + "\x00"

	fragmentShaderSource4 = `
	#version 410
    
	out vec4 frag_colour;

	in vec3 Normal;
	in vec3 FragPos;
	in vec2 TexCoords;

	uniform vec3 viewPos;

	struct Material {
		sampler2D diffuse;
		sampler2D specular;
		float shininess;
	}; 

	uniform Material material;

	struct Light {
		vec3 position;
	
		vec3 ambient;
		vec3 diffuse;
		vec3 specular;
	};
	
	uniform Light light;

	void main() {
		// 环境光
		vec3 ambient = light.ambient * vec3(texture(material.diffuse, TexCoords));

		// 漫反射光
		vec3 norm = normalize(Normal);
		vec3 lightDir = normalize(light.position - FragPos);
		float diff = max(dot(norm, lightDir), 0.0);
		vec3 diffuse = light.diffuse * diff * vec3(texture(material.diffuse, TexCoords));

		// 镜面光
		vec3 viewDir = normalize(viewPos - FragPos);
		vec3 reflectDir = reflect(-lightDir, norm);
		float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shininess);
		vec3 specular = light.specular * spec * vec3(texture(material.specular, TexCoords));

		frag_colour = vec4(ambient + diffuse + specular, 1.0);
	}
	` + "\x00"

	vertexShaderSource5 = `
	#version 410

	in vec3 aPos;
	in vec3 aNormal;
	in vec2 aTexCoord;

	out vec3 FragPos;
	out vec3 Normal;
	out vec2 TexCoords;

	uniform mat4 model;
	uniform mat4 view;
	uniform mat4 projection;
	uniform mat3 normalModel;

	void main() {
		Normal = normalModel * aNormal;
		FragPos = vec3(model * vec4(aPos, 1.0));
		gl_Position = projection * view * model * vec4(aPos, 1.0);
		TexCoords = vec2(aTexCoord.x, 1.0-aTexCoord.y);
	}
	` + "\x00"

	fragmentShaderSource5 = `
	#version 410
    
	out vec4 frag_colour;

	in vec3 Normal;
	in vec3 FragPos;
	in vec2 TexCoords;

	uniform vec3 viewPos;

	struct Material {
		sampler2D diffuse;
		sampler2D specular;
		float shininess;
	}; 

	uniform Material material;

	struct Light {
		vec3 direction;
	
		vec3 ambient;
		vec3 diffuse;
		vec3 specular;
	};
	
	uniform Light light;

	void main() {
		// 环境光
		vec3 ambient = light.ambient * vec3(texture(material.diffuse, TexCoords));

		// 漫反射光
		vec3 norm = normalize(Normal);
		vec3 lightDir = normalize(-light.direction);
		float diff = max(dot(norm, lightDir), 0.0);
		vec3 diffuse = light.diffuse * diff * vec3(texture(material.diffuse, TexCoords));

		// 镜面光
		vec3 viewDir = normalize(viewPos - FragPos);
		vec3 reflectDir = reflect(-lightDir, norm);
		float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shininess);
		vec3 specular = light.specular * spec * vec3(texture(material.specular, TexCoords));

		frag_colour = vec4(ambient + diffuse + specular, 1.0);
	}
	` + "\x00"

	vertexShaderSource6 = `
	#version 410

	in vec3 aPos;
	in vec3 aNormal;
	in vec2 aTexCoord;

	out vec3 FragPos;
	out vec3 Normal;
	out vec2 TexCoords;

	uniform mat4 model;
	uniform mat4 view;
	uniform mat4 projection;
	uniform mat3 normalModel;

	void main() {
		Normal = normalModel * aNormal;
		FragPos = vec3(model * vec4(aPos, 1.0));
		gl_Position = projection * view * model * vec4(aPos, 1.0);
		TexCoords = vec2(aTexCoord.x, 1.0-aTexCoord.y);
	}
	` + "\x00"

	fragmentShaderSource6 = `
	#version 410
    
	out vec4 frag_colour;

	in vec3 Normal;
	in vec3 FragPos;
	in vec2 TexCoords;

	uniform vec3 viewPos;

	struct Material {
		sampler2D diffuse;
		sampler2D specular;
		float shininess;
	}; 

	uniform Material material;

	struct Light {
		vec3 position;
	
		vec3 ambient;
		vec3 diffuse;
		vec3 specular;

		float constant;
		float linear;
		float quadratic;
	};
	
	uniform Light light;

	void main() {
		float distance = length(light.position - FragPos);
		float attenuation = 1.0 / (light.constant + light.linear * distance + light.quadratic * (distance * distance));

		// 环境光
		vec3 ambient = light.ambient * vec3(texture(material.diffuse, TexCoords));
		ambient *= attenuation;

		// 漫反射光
		vec3 norm = normalize(Normal);
		vec3 lightDir = normalize(light.position - FragPos);
		float diff = max(dot(norm, lightDir), 0.0);
		vec3 diffuse = light.diffuse * diff * vec3(texture(material.diffuse, TexCoords));
		diffuse *= attenuation;

		// 镜面光
		vec3 viewDir = normalize(viewPos - FragPos);
		vec3 reflectDir = reflect(-lightDir, norm);
		float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shininess);
		vec3 specular = light.specular * spec * vec3(texture(material.specular, TexCoords));
		specular *= attenuation;

		frag_colour = vec4(ambient + diffuse + specular, 1.0);
	}
	` + "\x00"

	vertexShaderSource7 = `
	#version 410

	in vec3 aPos;
	in vec3 aNormal;
	in vec2 aTexCoord;

	out vec3 FragPos;
	out vec3 Normal;
	out vec2 TexCoords;

	uniform mat4 model;
	uniform mat4 view;
	uniform mat4 projection;
	uniform mat3 normalModel;

	void main() {
		Normal = normalModel * aNormal;
		FragPos = vec3(model * vec4(aPos, 1.0));
		gl_Position = projection * view * model * vec4(aPos, 1.0);
		TexCoords = vec2(aTexCoord.x, 1.0-aTexCoord.y);
	}
	` + "\x00"

	fragmentShaderSource7 = `
	#version 410
    
	out vec4 frag_colour;

	in vec3 Normal;
	in vec3 FragPos;
	in vec2 TexCoords;

	uniform vec3 viewPos;

	struct Material {
		sampler2D diffuse;
		sampler2D specular;
		float shininess;
	}; 

	uniform Material material;

	struct Light {
		vec3 position;
	
		vec3 ambient;
		vec3 diffuse;
		vec3 specular;

		float constant;
		float linear;
		float quadratic;

    	vec3  direction;
    	float cutOff;
	};
	
	uniform Light light;

	void main() {
		vec3 lightDir = normalize(light.position - FragPos);

		float theta = dot(lightDir, normalize(-light.direction));

		if(theta > light.cutOff) {       
			float distance = length(light.position - FragPos);
			float attenuation = 1.0 / (light.constant + light.linear * distance + light.quadratic * (distance * distance));

			// 环境光
			vec3 ambient = light.ambient * vec3(texture(material.diffuse, TexCoords));

			// 漫反射光
			vec3 norm = normalize(Normal);
			float diff = max(dot(norm, lightDir), 0.0);
			vec3 diffuse = light.diffuse * diff * vec3(texture(material.diffuse, TexCoords));
			diffuse *= attenuation;

			// 镜面光
			vec3 viewDir = normalize(viewPos - FragPos);
			vec3 reflectDir = reflect(-lightDir, norm);
			float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shininess);
			vec3 specular = light.specular * spec * vec3(texture(material.specular, TexCoords));
			specular *= attenuation;

			frag_colour = vec4(ambient + diffuse + specular, 1.0);
		} else {
			// 否则，使用环境光，让场景在聚光之外时不至于完全黑暗
			frag_colour = vec4(light.ambient * vec3(texture(material.diffuse, TexCoords)), 1.0);
		}
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

	vertices2 = []float32{
		// positions      // normals      // texture coords
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0,
		0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 0.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 1.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0,
		0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0,

		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 0.0,
		-0.5, 0.5, -0.5, -1.0, 0.0, 0.0, 1.0, 1.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 1.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 1.0,
		-0.5, -0.5, 0.5, -1.0, 0.0, 0.0, 0.0, 0.0,
		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 0.0, 0.0, 1.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 1.0, 1.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 1.0, 0.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0,
		0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 1.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0,
	}
)

// 环境光+漫反射光+镜面光
func Run() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "light")
	defer glfw.Terminate()

	pointNum := int32(len(vertices)) / 5

	program1, _ := util.InitOpenGL(vertexShaderSource, fragmentShaderSource)
	vao1 := util.MakeVaoWithAttrib(program1, vertices, nil, []util.VertAttrib{{Name: "aPos", Size: 3}, {Name: "aNormal", Size: 3}})
	program2, _ := util.MakeProgram(vertexShaderSourceLight, fragmentShaderSourceLight)
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

// 环境光+漫反射光+镜面光
// 材质
func Run2() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "light")
	defer glfw.Terminate()

	pointNum := int32(len(vertices)) / 5

	program1, _ := util.InitOpenGL(vertexShaderSource2, fragmentShaderSource2)
	vao1 := util.MakeVaoWithAttrib(program1, vertices, nil, []util.VertAttrib{{Name: "aPos", Size: 3}, {Name: "aNormal", Size: 3}})
	program2, _ := util.MakeProgram(vertexShaderSourceLight, fragmentShaderSourceLight)
	vao2 := util.MakeVaoWithAttrib(program2, vertices, nil, []util.VertAttrib{{Name: "aPos", Size: 3}, {Name: "aNormal", Size: 3}})

	gl.ClearColor(0.1, 0.1, 0.1, 1.0)
	gl.Enable(gl.DEPTH_TEST)

	lightPos := mgl32.Vec3{2, 0, 0}

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

		gl.Uniform3f(gl.GetUniformLocation(program1, gl.Str("light.ambient\x00")), 0.2, 0.2, 0.2)
		gl.Uniform3f(gl.GetUniformLocation(program1, gl.Str("light.diffuse\x00")), 0.5, 0.5, 0.5)
		gl.Uniform3f(gl.GetUniformLocation(program1, gl.Str("light.specular\x00")), 1, 1, 1)
		gl.Uniform3fv(gl.GetUniformLocation(program1, gl.Str("light.position\x00")), 1, &lightPos[0])

		gl.Uniform3f(gl.GetUniformLocation(program1, gl.Str("material.ambient\x00")), 1, 0.5, 0.31)
		gl.Uniform3f(gl.GetUniformLocation(program1, gl.Str("material.diffuse\x00")), 1, 0.5, 0.31)
		gl.Uniform3f(gl.GetUniformLocation(program1, gl.Str("material.specular\x00")), 0.5, 0.5, 0.5)
		gl.Uniform1f(gl.GetUniformLocation(program1, gl.Str("material.shininess\x00")), 32)

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

// 环境光+漫反射光+镜面光
// 漫反射贴图（光照贴图）
func Run3() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "light")
	defer glfw.Terminate()

	pointNum := int32(len(vertices2)) / 8

	program1, _ := util.InitOpenGL(vertexShaderSource3, fragmentShaderSource3)
	vao1 := util.MakeVaoWithAttrib(program1, vertices2, nil, []util.VertAttrib{{Name: "aPos", Size: 3}, {Name: "aNormal", Size: 3}, {Name: "aTexCoord", Size: 2}})

	program2, _ := util.MakeProgram(vertexShaderSourceLight, fragmentShaderSourceLight)
	vao2 := util.MakeVaoWithAttrib(program2, vertices2, nil, []util.VertAttrib{{Name: "aPos", Size: 3}, {Name: "aNormal", Size: 3}, {Name: "aTexCoord", Size: 2}})

	texture := util.MakeTexture("light/container2.png")

	gl.ClearColor(0.1, 0.1, 0.1, 1.0)
	gl.Enable(gl.DEPTH_TEST)

	lightPos := mgl32.Vec3{2, 0, 0}

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

		gl.Uniform3f(gl.GetUniformLocation(program1, gl.Str("light.ambient\x00")), 0.2, 0.2, 0.2)
		gl.Uniform3f(gl.GetUniformLocation(program1, gl.Str("light.diffuse\x00")), 0.5, 0.5, 0.5)
		gl.Uniform3f(gl.GetUniformLocation(program1, gl.Str("light.specular\x00")), 1, 1, 1)
		gl.Uniform3fv(gl.GetUniformLocation(program1, gl.Str("light.position\x00")), 1, &lightPos[0])

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.Uniform1i(gl.GetUniformLocation(program1, gl.Str("material.diffuse"+"\x00")), 0)
		gl.Uniform3f(gl.GetUniformLocation(program1, gl.Str("material.specular\x00")), 0.5, 0.5, 0.5)
		gl.Uniform1f(gl.GetUniformLocation(program1, gl.Str("material.shininess\x00")), 32)

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

// 环境光+漫反射光+镜面光
// 漫反射光贴图 + 镜面光贴图
func Run4() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "light")
	defer glfw.Terminate()

	pointNum := int32(len(vertices2)) / 8

	program1, _ := util.InitOpenGL(vertexShaderSource4, fragmentShaderSource4)
	vao1 := util.MakeVaoWithAttrib(program1, vertices2, nil, []util.VertAttrib{{Name: "aPos", Size: 3}, {Name: "aNormal", Size: 3}, {Name: "aTexCoord", Size: 2}})

	program2, _ := util.MakeProgram(vertexShaderSourceLight, fragmentShaderSourceLight)
	vao2 := util.MakeVaoWithAttrib(program2, vertices2, nil, []util.VertAttrib{{Name: "aPos", Size: 3}, {Name: "aNormal", Size: 3}, {Name: "aTexCoord", Size: 2}})

	texture0 := util.MakeTexture("light/container2.png")
	texture1 := util.MakeTexture("light/container2_specular.png")

	gl.ClearColor(0.1, 0.1, 0.1, 1.0)
	gl.Enable(gl.DEPTH_TEST)

	lightPos := mgl32.Vec3{2, 0, 0}

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

		gl.Uniform3f(gl.GetUniformLocation(program1, gl.Str("light.ambient\x00")), 0.2, 0.2, 0.2)
		gl.Uniform3f(gl.GetUniformLocation(program1, gl.Str("light.diffuse\x00")), 0.5, 0.5, 0.5)
		gl.Uniform3f(gl.GetUniformLocation(program1, gl.Str("light.specular\x00")), 1, 1, 1)
		gl.Uniform3fv(gl.GetUniformLocation(program1, gl.Str("light.position\x00")), 1, &lightPos[0])

		// 漫反射贴图
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture0)
		gl.Uniform1i(gl.GetUniformLocation(program1, gl.Str("material.diffuse"+"\x00")), 0)

		// 镜面光贴图
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.Uniform1i(gl.GetUniformLocation(program1, gl.Str("material.specular"+"\x00")), 1)

		gl.Uniform1f(gl.GetUniformLocation(program1, gl.Str("material.shininess\x00")), 32)

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

// 平行光/定向光
func Run5() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "light")
	defer glfw.Terminate()

	pointNum := int32(len(vertices2)) / 8

	program1, _ := util.InitOpenGL(vertexShaderSource5, fragmentShaderSource5)
	vao1 := util.MakeVaoWithAttrib(program1, vertices2, nil, []util.VertAttrib{{Name: "aPos", Size: 3}, {Name: "aNormal", Size: 3}, {Name: "aTexCoord", Size: 2}})

	texture0 := util.MakeTexture("light/container2.png")
	texture1 := util.MakeTexture("light/container2_specular.png")

	gl.ClearColor(0.1, 0.1, 0.1, 1.0)
	gl.Enable(gl.DEPTH_TEST)

	// 取反之后为照射方向，也就是说此时是从 Z=1 照射到 z=0
	lightDirection := mgl32.Vec3{0, 0, -1}

	cameraPos := mgl32.Vec3{0, 0, 3}
	cameraFront := mgl32.Vec3{0, 0, -1}
	cameraUp := mgl32.Vec3{0, 1, 0}

	camera := util.NewCamera(cameraPos, cameraFront, cameraUp, width, height)
	camera.SetCursorPosCallback(window)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		lightDirection = mgl32.Vec3{0, float32(math.Sin(glfw.GetTime())), -1}

		// 画箱子
		gl.UseProgram(program1)

		gl.Uniform3f(gl.GetUniformLocation(program1, gl.Str("light.ambient\x00")), 0.2, 0.2, 0.2)
		gl.Uniform3f(gl.GetUniformLocation(program1, gl.Str("light.diffuse\x00")), 0.5, 0.5, 0.5)
		gl.Uniform3f(gl.GetUniformLocation(program1, gl.Str("light.specular\x00")), 1, 1, 1)
		gl.Uniform3fv(gl.GetUniformLocation(program1, gl.Str("light.direction\x00")), 1, &lightDirection[0])

		// 漫反射贴图
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture0)
		gl.Uniform1i(gl.GetUniformLocation(program1, gl.Str("material.diffuse"+"\x00")), 0)

		// 镜面光贴图
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.Uniform1i(gl.GetUniformLocation(program1, gl.Str("material.specular"+"\x00")), 1)

		gl.Uniform1f(gl.GetUniformLocation(program1, gl.Str("material.shininess\x00")), 32)

		gl.Uniform3fv(gl.GetUniformLocation(program1, gl.Str("viewPos\x00")), 1, &camera.CameraPos[0])
		view := camera.LookAt()
		projection := camera.Perspective()

		model1 := mgl32.Ident4()
		normalModel := model1.Inv().Transpose().Mat3()
		gl.UniformMatrix4fv(gl.GetUniformLocation(program1, gl.Str("view\x00")), 1, false, &view[0])
		gl.UniformMatrix4fv(gl.GetUniformLocation(program1, gl.Str("projection\x00")), 1, false, &projection[0])
		gl.UniformMatrix3fv(gl.GetUniformLocation(program1, gl.Str("normalModel\x00")), 1, false, &normalModel[0])
		gl.BindVertexArray(vao1)
		gl.DrawArrays(gl.TRIANGLES, 0, pointNum)

		for k, v := range more {
			model1 = mgl32.HomogRotate3D(mgl32.DegToRad(float32(20*k)), mgl32.Vec3{1, 0.3, 0.5})
			model1 = mgl32.Translate3D(v.X(), v.Y(), v.Z()).Mul4(model1)
			gl.UniformMatrix4fv(gl.GetUniformLocation(program1, gl.Str("model\x00")), 1, false, &model1[0])
			gl.DrawArrays(gl.TRIANGLES, 0, pointNum)
		}

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

// 点光源
func Run6() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "light")
	defer glfw.Terminate()

	pointNum := int32(len(vertices2)) / 8

	program1, _ := util.InitOpenGL(vertexShaderSource6, fragmentShaderSource6)
	vao1 := util.MakeVaoWithAttrib(program1, vertices2, nil, []util.VertAttrib{{Name: "aPos", Size: 3}, {Name: "aNormal", Size: 3}, {Name: "aTexCoord", Size: 2}})

	program2, _ := util.MakeProgram(vertexShaderSourceLight, fragmentShaderSourceLight)
	vao2 := util.MakeVaoWithAttrib(program2, vertices2, nil, []util.VertAttrib{{Name: "aPos", Size: 3}, {Name: "aNormal", Size: 3}, {Name: "aTexCoord", Size: 2}})

	texture0 := util.MakeTexture("light/container2.png")
	texture1 := util.MakeTexture("light/container2_specular.png")

	gl.ClearColor(0.1, 0.1, 0.1, 1.0)
	gl.Enable(gl.DEPTH_TEST)

	lightPos := mgl32.Vec3{2, 0, 0}

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

		gl.Uniform3f(gl.GetUniformLocation(program1, gl.Str("light.ambient\x00")), 0.2, 0.2, 0.2)
		gl.Uniform3f(gl.GetUniformLocation(program1, gl.Str("light.diffuse\x00")), 0.5, 0.5, 0.5)
		gl.Uniform3f(gl.GetUniformLocation(program1, gl.Str("light.specular\x00")), 1, 1, 1)
		gl.Uniform3fv(gl.GetUniformLocation(program1, gl.Str("light.position\x00")), 1, &lightPos[0])

		// 衰减因子
		gl.Uniform1f(gl.GetUniformLocation(program1, gl.Str("light.constant\x00")), 1)
		gl.Uniform1f(gl.GetUniformLocation(program1, gl.Str("light.linear\x00")), 0.09)
		gl.Uniform1f(gl.GetUniformLocation(program1, gl.Str("light.quadratic\x00")), 0.032)

		// 漫反射贴图
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture0)
		gl.Uniform1i(gl.GetUniformLocation(program1, gl.Str("material.diffuse"+"\x00")), 0)

		// 镜面光贴图
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.Uniform1i(gl.GetUniformLocation(program1, gl.Str("material.specular"+"\x00")), 1)

		gl.Uniform1f(gl.GetUniformLocation(program1, gl.Str("material.shininess\x00")), 32)

		gl.Uniform3fv(gl.GetUniformLocation(program1, gl.Str("viewPos\x00")), 1, &camera.CameraPos[0])
		view := camera.LookAt()
		projection := camera.Perspective()

		model1 := mgl32.Ident4()
		normalModel := model1.Inv().Transpose().Mat3()
		gl.UniformMatrix4fv(gl.GetUniformLocation(program1, gl.Str("view\x00")), 1, false, &view[0])
		gl.UniformMatrix4fv(gl.GetUniformLocation(program1, gl.Str("projection\x00")), 1, false, &projection[0])
		gl.UniformMatrix3fv(gl.GetUniformLocation(program1, gl.Str("normalModel\x00")), 1, false, &normalModel[0])
		gl.BindVertexArray(vao1)
		gl.DrawArrays(gl.TRIANGLES, 0, pointNum)

		for k, v := range more {
			model1 = mgl32.HomogRotate3D(mgl32.DegToRad(float32(20*k)), mgl32.Vec3{1, 0.3, 0.5})
			model1 = mgl32.Translate3D(v.X(), v.Y(), v.Z()).Mul4(model1)
			gl.UniformMatrix4fv(gl.GetUniformLocation(program1, gl.Str("model\x00")), 1, false, &model1[0])
			gl.DrawArrays(gl.TRIANGLES, 0, pointNum)
		}

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

// TODO
// 聚光
func Run7() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "light")
	defer glfw.Terminate()

	pointNum := int32(len(vertices2)) / 8

	program1, _ := util.InitOpenGL(vertexShaderSource6, fragmentShaderSource6)
	vao1 := util.MakeVaoWithAttrib(program1, vertices2, nil, []util.VertAttrib{{Name: "aPos", Size: 3}, {Name: "aNormal", Size: 3}, {Name: "aTexCoord", Size: 2}})

	texture0 := util.MakeTexture("light/container2.png")
	texture1 := util.MakeTexture("light/container2_specular.png")

	gl.ClearColor(0.1, 0.1, 0.1, 1.0)
	gl.Enable(gl.DEPTH_TEST)

	cameraPos := mgl32.Vec3{0, 0, 3}
	cameraFront := mgl32.Vec3{0, 0, -1}
	cameraUp := mgl32.Vec3{0, 1, 0}

	camera := util.NewCamera(cameraPos, cameraFront, cameraUp, width, height)
	camera.SetCursorPosCallback(window)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// 画箱子
		gl.UseProgram(program1)

		gl.Uniform3f(gl.GetUniformLocation(program1, gl.Str("light.ambient\x00")), 0.2, 0.2, 0.2)
		gl.Uniform3f(gl.GetUniformLocation(program1, gl.Str("light.diffuse\x00")), 0.5, 0.5, 0.5)
		gl.Uniform3f(gl.GetUniformLocation(program1, gl.Str("light.specular\x00")), 1, 1, 1)

		gl.Uniform3fv(gl.GetUniformLocation(program1, gl.Str("light.position\x00")), 1, &camera.CameraPos[0])
		gl.Uniform3fv(gl.GetUniformLocation(program1, gl.Str("light.direction\x00")), 1, &camera.CameraFront[0])
		gl.Uniform1f(gl.GetUniformLocation(program1, gl.Str("light.cutOff\x00")), float32(math.Cos(float64(mgl32.DegToRad(12.5)))))

		// 衰减因子
		gl.Uniform1f(gl.GetUniformLocation(program1, gl.Str("light.constant\x00")), 1)
		gl.Uniform1f(gl.GetUniformLocation(program1, gl.Str("light.linear\x00")), 0.09)
		gl.Uniform1f(gl.GetUniformLocation(program1, gl.Str("light.quadratic\x00")), 0.032)

		// 漫反射贴图
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture0)
		gl.Uniform1i(gl.GetUniformLocation(program1, gl.Str("material.diffuse"+"\x00")), 0)

		// 镜面光贴图
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.Uniform1i(gl.GetUniformLocation(program1, gl.Str("material.specular"+"\x00")), 1)

		gl.Uniform1f(gl.GetUniformLocation(program1, gl.Str("material.shininess\x00")), 32)

		gl.Uniform3fv(gl.GetUniformLocation(program1, gl.Str("viewPos\x00")), 1, &camera.CameraPos[0])
		view := camera.LookAt()
		projection := camera.Perspective()

		model1 := mgl32.Ident4()
		normalModel := model1.Inv().Transpose().Mat3()
		gl.UniformMatrix4fv(gl.GetUniformLocation(program1, gl.Str("view\x00")), 1, false, &view[0])
		gl.UniformMatrix4fv(gl.GetUniformLocation(program1, gl.Str("projection\x00")), 1, false, &projection[0])
		gl.UniformMatrix3fv(gl.GetUniformLocation(program1, gl.Str("normalModel\x00")), 1, false, &normalModel[0])
		gl.BindVertexArray(vao1)
		gl.DrawArrays(gl.TRIANGLES, 0, pointNum)

		for k, v := range more {
			model1 = mgl32.HomogRotate3D(mgl32.DegToRad(float32(20*k)), mgl32.Vec3{1, 0.3, 0.5})
			model1 = mgl32.Translate3D(v.X(), v.Y(), v.Z()).Mul4(model1)
			gl.UniformMatrix4fv(gl.GetUniformLocation(program1, gl.Str("model\x00")), 1, false, &model1[0])
			gl.DrawArrays(gl.TRIANGLES, 0, pointNum)
		}

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
