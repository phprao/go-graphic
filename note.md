### 开发环境与依赖

github.com/go-gl/gl/v4.1-core/gl
github.com/go-gl/glfw/v3.2/glfw

OpenGL只提供了绘图功能，创建窗口是需要自己完成的。这就需要学习相应操作系统的创建窗口方法，比较复杂，并且每个操作系统都不同。为简化创建窗口的过程，可以使用专门的窗口库，例如GLUT、GLFW等。由于GLUT已经是90年代的东西了（不过后来还有freeglut），而GLFW是新的，因此建议使用GLFW。

GLFW 是配合 OpenGL 使用的轻量级工具程序库，缩写自 Graphics Library Framework（图形库框架）。GLFW 的主要功能是创建并管理窗口和 OpenGL 上下文，同时还提供了处理手柄、键盘、鼠标输入，以及事件处理的功能。

gl 和 glfw 都使用了 cgo ，因此需要安装 gcc，也就是 MinGW，安装参考：https://blog.csdn.net/raoxiaoya/article/details/130820906

UI控件库：Nuclear是个C语言开发的UI控件库，有golang的绑定，可以帮我们处理简单的界面显示问题。github.com/golang-ui/nuklear

关于glad（还没使用过）：
glad的API包括：窗口操作、窗口初始化、窗口大小、位置调整等；回调函数；响应刷新消息、键盘消息、鼠标消息、定时器函数等；创建复杂三维体；菜单函数；程序运行函数等
安装glad：https://glad.dav1d.de/

glfw需要安装吗？
golang使用cgo来调用clang代码，既可以调用以编译好的动态库和静态库，也可以将clang代码放在golang项目中，这样能直接调用，显然包`github.com/go-gl/glfw/v3.2/glfw`使用的是后者，这样就不需要你的电脑上去安装glfw程序了，在`README.md`中也有说明

* GLFW C library source is included and built automatically as part of the Go package. But you need to make sure you have dependencies of GLFW:
    * On macOS, you need Xcode or Command Line Tools for Xcode (xcode-select --install) for required headers and libraries.
    * On Ubuntu/Debian-like Linux distributions, you need libgl1-mesa-dev and xorg-dev packages.
    * On CentOS/Fedora-like Linux distributions, you need libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel libXxf86vm-devel packages.
    * On FreeBSD, you need the package pkgconf. To build for X, you also need the package xorg; and to build for Wayland, you need the package wayland.
    * On NetBSD, to build for X, you need the X11 sets installed. These are included in all graphical installs, and can be added to the system with sysinst(8) on non-graphical systems. Wayland support is incomplete, due to missing wscons support in upstream GLFW. To attempt to build for Wayland, you need to install the wayland libepoll-shim packages and set the environment variable PKG_CONFIG_PATH=/usr/pkg/libdata/pkgconfig.
    * On OpenBSD, you need the X11 sets. These are installed by default, and can be added from the ramdisk kernel at any time.
    * See here for full details.
* Go 1.4+ is required on Windows (otherwise you must use MinGW v4.8.1 exactly, see Go issue 8811).

但是也提到，由于glfw还依赖了其他东西，因此在特定系统上，还需要安装这些依赖，但是在windows系统上则不需要额外安装什么。

对于opengl，目前操作系统都自带有，也不用安装：
Requirements:
* A cgo compiler (typically gcc).
* On Ubuntu/Debian-based systems, the libgl1-mesa-dev package.

### 参考文档
opengl官方文档：https://www.opengl.org/ --> Documentation --> Current OpenGL Version --> OpenGL 4.1 --> API Core Profile --> https://registry.khronos.org/OpenGL/specs/gl/glspec41.core.pdf

opengl官方教程：http://www.opengl-tutorial.org/cn/

学习opengl的网站
https://learnopengl-cn.github.io/

GLFW文档：https://www.glfw.org/docs/latest/

### 实践

#### 初始化glfw
glfw是用来操作窗口的，因此需要先初始化。
glfw相关函数需要在主线程 main thread 中运行，因此需要在主goroutine中调用`runtime.LockOSThread()`，然后再调用glfw。

```go
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(width, height, "Conway's Game of Life", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	return window
}
```
##### glfw.WindowHint()
对应的函数是 glfwWindowHint()，设置窗口的一些属性值。具体有哪些属性可以参考 https://www.glfw.org/docs/latest/window_guide.html#window_hints。

> GLFW_CONTEXT_VERSION_MAJOR and GLFW_CONTEXT_VERSION_MINOR specify the client API version that the created context must be compatible with. The exact behavior of these hints depend on the requested client API.
>
> While there is no way to ask the driver for a context of the highest supported version, GLFW will attempt to provide this when you ask for a version 1.0 context, which is the default for these hints.
>
> Do not confuse these hints with GLFW_VERSION_MAJOR and GLFW_VERSION_MINOR, which provide the API version of the GLFW header.
>
> 也就是设置 opengl 的版本号，此处我们使用 opengl 4.1

> GLFW_OPENGL_PROFILE specifies which OpenGL profile to create the context for. Possible values are one of GLFW_OPENGL_CORE_PROFILE or GLFW_OPENGL_COMPAT_PROFILE, or GLFW_OPENGL_ANY_PROFILE to not request a specific profile. If requesting an OpenGL version below 3.2, GLFW_OPENGL_ANY_PROFILE must be used. If OpenGL ES is requested, this hint is ignored.

> GLFW_OPENGL_FORWARD_COMPAT specifies whether the OpenGL context should be forward-compatible, i.e. one where all functionality deprecated in the requested version of OpenGL is removed. This must only be used if the requested OpenGL version is 3.0 or above. If OpenGL ES is requested, this hint is ignored.

#####  window.MakeContextCurrent()
我们在glfw官网搜`glfwMakeContextCurrent`，解释如下：

`void glfwMakeContextCurrent(GLFWwindow *window)`
> This function makes the OpenGL or OpenGL ES context of the specified window current on the calling thread. A context must only be made current on a single thread at a time and each thread can have only a single current context at a time.
> 
> When moving a context between threads, you must make it non-current on the old thread before making it current on the new one.
> 
> By default, making a context non-current implicitly forces a pipeline flush. On machines that support GL_KHR_context_flush_control, you can control whether a context performs this flush by setting the GLFW_CONTEXT_RELEASE_BEHAVIOR hint.
> 
> The specified window must have an OpenGL or OpenGL ES context. Specifying a window without a context will generate a GLFW_NO_WINDOW_CONTEXT error.

> Before you can use the OpenGL API, you must have a current OpenGL context.
> The context will remain current until you make another context current or until the window owning the current context is destroyed.

在使用 opengl api 之前需要先设置上下文，可以将当前窗口绑定为当前上下文。

#### 初始化opengl
opengl
opengl相关函数需要在主线程 main thread 中运行。
```go
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()

	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)

	gl.LinkProgram(prog)
	return prog
}
```

#### 构建VBO，VAO
VAO和VBO都是用来存储顶点信息的，并把这些信息送入顶点着色器。一个VAO可以对应多个VBO。

在OpenGL程序中，同时只会有一个VAO被绑定到opengl，当然你可以操作完一个之后再绑定到另一个VAO。

VBO是顶点缓冲对象(Vertex Buffer Objects, VBO)，包含了顶点的3d坐标和颜色。但它们是按同类数组存储的，存放在一片显存空间中，程序并不知道这些数字哪个代表3d坐标，哪个代表颜色。

VAO是顶点数组对象(Vertex Array Object, VAO)，用来表示这些数字的第几位分别代表顶点的什么属性。比如这些数字的第1-3位代表3d的xyz坐标，第4-7位代表rbg颜色和透明度。

EBO是元素缓冲对象(Element Buffer Object，EBO)，EBO是一个缓冲区，就像一个VBO一样，它存储 OpenGL 用来决定要绘制哪些顶点的索引,设置顶点的绘制顺序。EBO由VAO进行绑定
![image-20230524113220299](D:\dev\php\magook\trunk\server\md\img\image-20230524113220299.png)

如图所示：VAO1中下标为0的指针attribute pointer[0]对应VBO1中的pos[0],表示VOB1数组中下标为0位置代表坐标pos。

VAO2中下标为0的指针attribute pointer[0]对应VBO2中的pos[0]，VAO2中下标为1的指针attribute pointer[1]对应VBO2中的col[0]，表示VOB2数组中下标为1位置代表颜色col。

我们用VBO来存储数据，而用VAO来告诉计算机这些数据分别有什么属性、起什么作用。

VBO是 CPU 和 GPU 之间传递信息的桥梁，我们把数据存入VBO(这一步在CPU执行)，然后VBO会自动把数据送入GPU。但是，对GPU来说，VBO中存的就只是一堆数字而已，要怎么解释它们呢？这就要用到VAO了。

VBO是在显卡存储空间中开辟出的一块内存缓存区，用于存储顶点的各类属性信息，如顶点坐标，顶点法向量，顶点颜色数据等。在渲染时，可以直接从VBO中取出顶点的各类属性数据，由于VBO在显存而不是在内存中，不需要从CPU传输数据，处理效率更高。

```go
func makeVao(points []float32) uint32 {
	var vbo uint32

    // 在显卡中开辟一块空间，创建顶点缓存对象，个数为1
	gl.GenBuffers(1, &vbo)

    // 创建的VBO可用来保存不同类型的顶点数据，因此需要指定VBO的类型
    // 可选类型：GL_ARRAY_BUFFER, GL_ELEMENT_ARRAY_BUFFER, GL_PIXEL_PACK_BUFFER, GL_PIXEL_UNPACK_BUFFER
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

    // 将内存中的数据传递到显卡中
    // 4*len(points) 代表总的字节数，因为是32位的
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
    // 创建顶点数组对象，个数为1
	gl.GenVertexArrays(1, &vao)

    // 后面的两个函数都是需要先绑定VAO的
	gl.BindVertexArray(vao)
    // 设置 vertex attribute 的状态，默认是disabled
    // 参数index，目前是第0个
，	gl.EnableVertexAttribArray(0)
    // 参数index，目前是第0个
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	return vao
}
```

顶点坐标为`[]float32`类型，依次为`X, Y, Z`，`窗口中心点为原点，向右为X正，上为Y正`，取值范围`[-1,1]`，
```go
triangle = []float32{
    0, 0.5, 0,
    -0.5, -0.5, 0,
    0.5, -0.5, 0,
}
```

当我们特别谈论到顶点着色器的时候，每个输入变量也叫顶点属性(Vertex Attribute)。我们能声明的顶点属性是有上限的，它一般由硬件来决定。OpenGL确保至少有16个包含4分量的顶点属性可用，但是有些硬件或许允许更多的顶点属性，你可以查询 GL_MAX_VERTEX_ATTRIBS 来获取具体的上限。通常情况下它至少会返回16个，大部分情况下是够用了。

默认情况下，出于性能考虑，所有顶点着色器的属性（Attribute）变量都是关闭的，意味着数据在着色器端是不可见的，哪怕数据已经上传到GPU，由 glEnableVertexAttribArray 启用指定属性，才可在顶点着色器中访问逐个顶点的属性数据。glVertexAttribPointer或VBO只是建立CPU和GPU之间的逻辑连接，从而实现了CPU数据上传至GPU。但是，数据在GPU端是否可见，即，着色器能否读取到数据，由是否启用了对应的属性决定，这就是 glEnableVertexAttribArray 的功能，允许顶点着色器读取GPU（服务器端）数据。

那么，glEnableVertexAttribArray应该在glVertexAttribPointer之前还是之后调用？答案是都可以，只要在绘图调用（glDraw*系列函数）前调用即可。

不用glEnableVertexAttribArray对应属性时绘制内容为清除缓冲区颜色。在取消glEnableVertexAttribArray(0);注释后我的代码就可以得到正常的绘制结果。

如果有多个VBO，可以依次设置，VBO的索引是根据生成的先后顺序依次递增的
```go
gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
gl.EnableVertexAttribArray(0)

gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, nil)
gl.EnableVertexAttribArray(1)

gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 0, nil)
gl.EnableVertexAttribArray(2)
```

我们发现函数 `gl.EnableVertexAttribArray(0)` 和 `gl.VertexAttribPointer()` 的参数并没有指定是操作哪个VAO，这是因为在前面调用了`gl.BindVertexArray(vao)`函数已经绑定到了VAO，关于这一点后面还会提到。

#### 顶点着色器
着色器是opengl内部的小程序，是由 GLSL (OpenGL Shader Language) 语言编写。

顶点着色器包含对一些顶点属性（数据）的基本处理。
```bash
vertexShaderSource = `
    #version 410
    in vec3 vp;
    void main() {
        gl_Position = vec4(vp, 1.0);
    }
` + "\x00"
```


#### 片元着色器
片元着色器的作用是计算出每一个像素点最终的颜色，通常片段着色器会包含3D场景的一些额外的数据，如光线，阴影等。

用 RGBA 形式的值通过 vec4 来定义我们图形的颜色。颜色值`[0, 1]`。

```bash
fragmentShaderSource = `
    #version 410
    out vec4 frag_colour;
    void main() {
        frag_colour = vec4(1, 1, 1, 1);
    }
` + "\x00"
```

同样需要注意的是这两个程序都是运行在 `#version 410` 版本下，如果你用的是 `OpenGL 2.1`，那你也可以改成 `#version 120`。

#### 画图
OpenGL中所有的图形都是通过分解成三角形的方式进行绘制。

使用VAO作为数据，在opengl program中画图，其中画图和颜色填充交给可编程管线（顶点着色器和片元着色器）来完成，最后呈现在window上。
```go
func draw(vao uint32, window *glfw.Window, prog uint32) {
    // 来清除上一帧在窗口中绘制的东西
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

    // 使用这个程序
	gl.UseProgram(prog)

    // 绑定VAO，可鞥会感到奇怪，明明在 makeVao 中已经调用 gl.BindVertexArray(vao) 进行了绑定，为什么这里还要再绑定一次呢？
    // 因为绑定的操作是为了后续的操作服务的，并且有可能在中途又绑定了别的VAO，所以最好是在每次调用跟VAO有关的函数之前绑定一次。
	gl.BindVertexArray(vao)

    // 绘制的类型mode：
    //   1、gl.TRIANGLES：每三个顶之间绘制三角形，之间不连接
    //   2、gl.TRIANGLE_FAN：以V0V1V2,V0V2V3,V0V3V4，……的形式绘制三角形
    //   3、gl.TRIANGLE_STRIP：以V0V1V2,V1V2V3,V2V3V4……的形式绘制三角形
    // first：一般从第一个数据开始
    // count：除以3为顶点个数
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3)) // 三角形

    // 事件循环：去检查是否有鼠标或者键盘事件
	glfw.PollEvents()

    // 交换缓冲区
    // 因为 GLFW（像其他图形库一样）使用双缓冲，也就是说你绘制的所有东西实际上是绘制到一个不可见的画布上，当
    // 你准备好进行展示的时候就把绘制的这些东西放到可见的画布中
	window.SwapBuffers()
}
```

#### 函数说明




