### 开发环境与依赖

`github.com/go-gl/gl/v4.1-core/gl`
`github.com/go-gl/glfw/v3.2/glfw`

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

对于opengl，目前windows操作系统都自带有，也不用安装：
Requirements:
* A cgo compiler (typically gcc).
* On Ubuntu/Debian-based systems, the libgl1-mesa-dev package.

opengl 只是一套编程接口，是一种规范，而具体的实现则是由不同的显卡厂商基于不同的操作系统来完成。

golang + opengl api --> 显卡驱动 --> 显示器

所以在开发的时候要注意opengl版本，显卡驱动，显卡的版本，大部分时候应该不会有问题，除非你的机器太古老。

### 参考文档
opengl官方文档：https://www.opengl.org/ --> Documentation --> Current OpenGL Version --> OpenGL 4.1 --> API Core Profile --> https://registry.khronos.org/OpenGL/specs/gl/glspec41.core.pdf

opengl官方教程：http://www.opengl-tutorial.org/cn/

学习opengl的网站
https://learnopengl-cn.github.io/
https://blog.csdn.net/weixin_42050609?type=blog

GLFW文档（最好用）：https://www.glfw.org/docs/latest/

在GLFW官方文档中函数都是`glfw`开头的，比如`glfw.SwapInterval()`对应`glfwSwapInterval()`。
在opengl官方文档中函数则没有`gl`前缀，比如`gl.BindVertexArray()`对应`BindVertexArray()`。

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

- `glfw.WindowHint(target Hint, hint int)`
对应的函数是 glfwWindowHint()，设置窗口和OpenGL上下文的一些属性值。具体有哪些属性可以参考 https://www.glfw.org/docs/latest/window_guide.html#window_hints。

> GLFW_CONTEXT_VERSION_MAJOR and GLFW_CONTEXT_VERSION_MINOR specify the client API version that the created context must be compatible with. The exact behavior of these hints depend on the requested client API.
>
> While there is no way to ask the driver for a context of the highest supported version, GLFW will attempt to provide this when you ask for a version 1.0 context, which is the default for these hints.
>
> Do not confuse these hints with GLFW_VERSION_MAJOR and GLFW_VERSION_MINOR, which provide the API version of the GLFW header.
>
> 也就是设置 opengl 的版本号，此处我们使用 opengl 4.1

> GLFW_OPENGL_PROFILE specifies which OpenGL profile to create the context for. Possible values are one of GLFW_OPENGL_CORE_PROFILE or GLFW_OPENGL_COMPAT_PROFILE, or GLFW_OPENGL_ANY_PROFILE to not request a specific profile. If requesting an OpenGL version below 3.2, GLFW_OPENGL_ANY_PROFILE must be used. If OpenGL ES is requested, this hint is ignored.

关于立即渲染模式(Immediate mode，也就是固定渲染模式)，这个模式下绘制图形很方便。OpenGL的大多数功能都被库隐藏起来，开发者很少有控制OpenGL如何进行计算的自由。而开发者迫切希望能有更多的灵活性。随着时间推移，规范越来越灵活，开发者对绘图细节有了更多的掌控。立即渲染模式确实容易使用和理解，但是效率太低。因此从OpenGL3.2开始，规范文档开始废弃立即渲染模式，并鼓励开发者在OpenGL的核心模式(Core-profile)下进行开发，这个分支的规范完全移除了旧的特性。

> GLFW_OPENGL_FORWARD_COMPAT specifies whether the OpenGL context should be forward-compatible, i.e. one where all functionality deprecated in the requested version of OpenGL is removed. This must only be used if the requested OpenGL version is 3.0 or above. If OpenGL ES is requested, this hint is ignored.

- `glfw.CreateWindow(width, height int, title string, monitor *Monitor, share *Window)`
根据`WindowHint`函数设置的参数来创建一个窗口和与之相关联的OpenGL or OpenGL ES上下文。

关于share参数，[官网解释](https://www.glfw.org/docs/latest/context_guide.html#context_sharing) 。也是一个窗口对象，意味着他两共享同一个OpenGL上下文，比如：
`second_window = glfwCreateWindow(640, 480, "Second Window", NULL, first_window)`
共享的数据包括`textures, vertex and element buffers`等等，但是具体是怎么实现的，则取决于操作系统和图形驱动。

关于monitor参数，[官网解释](https://www.glfw.org/docs/latest/monitor_guide.html#monitor_monitors)，监视器指的就是显示设备，显示设备之所以能够展示图形，是因为数据线一直在一帧一帧的向其输入数据，即便我的电脑现在看到的是静态的画面，其实它一直在重复的渲染。我们的显示器现在显示的是window10的桌面，这是因为电脑向监视器出入的是window10的桌面的图像，那我可不可以只向监视器输入一张图片呢，显然是可以的，这样的话监视器就会全屏展示你的图片，此时你看不到任何windows的元素，但是你需要不停的重复着向监视器输入这张图片，这样才能保持住画面，这就是全屏的含义，跟我们熟知的很多软件右上角的那个放大按钮是不一样的意思。

获取主监视器`glfwGetPrimaryMonitor`，获取已连接的监视器列表`glfwGetMonitors`，具体查看文档。

如果我们指定了监视器，那么显卡将会直接将你的数据给到监视器，这就是全屏的效果。
```go
window, err := glfw.CreateWindow(width, height, name, glfw.GetPrimaryMonitor(), nil)
```
当然，如果width和height跟监视器的宽高不一样，也会出现黑边。为此我们做一个优化。
```go
monitor := glfw.GetPrimaryMonitor()
videoMode := monitor.GetVideoMode()
glfw.WindowHint(glfw.RedBits, videoMode.RedBits)
glfw.WindowHint(glfw.GreenBits, videoMode.GreenBits)
glfw.WindowHint(glfw.BlueBits, videoMode.BlueBits)
glfw.WindowHint(glfw.RefreshRate, videoMode.RefreshRate)
window, err := glfw.CreateWindow(videoMode.Width, videoMode.Height, name, monitor, nil)
```

这让我想起来了平时通过视频播放器观看视频的效果，可以点击全屏，也可以退出全屏，这个也好实现，我们增加两个按键事件。
```go
// 按键K
if window.GetKey(glfw.KeyK) == glfw.Press {
    // 设置显示器，X坐标，Y坐标，图像的宽高，最后是刷新频率
    // 我们这是从左上角开始，宽高为满屏
    window.SetMonitor(glfw.GetPrimaryMonitor(), 0, 0, 1920, 1080, 1)
}
// 按键M
if window.GetKey(glfw.KeyM) == glfw.Press {
    // 缩回来后的位置和宽高
    window.SetMonitor(nil, 100, 100, 500, 500, 1)
}
```
为了适应不同的监视器，我们来优化一下
```go
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
```

如果不指定监视器，那么就是`windowed mode`，在win10上它的效果是创建一个黑色窗口，你的数据渲染到这个窗口里，此时使用的监视器还是主监视器(primary monitor)。

[窗口位置](https://www.glfw.org/docs/latest/window_guide.html#window_pos)可以通过函数`glfwSetWindowPos(window, 100, 100)`，以左上角为原点。获取窗口位置`glfwGetWindowPos`。

关于窗口居中，显然先要得到监视器的宽高，我们创建一个非全屏的窗口，默认使用的是主监视器，另外还需要知道显示模式`video mode`，可以理解为分辨率相关的，比如我的分辨率是1920*1080.
```go
sw := glfw.GetPrimaryMonitor().GetVideoMode().Width // 1920
sh := glfw.GetPrimaryMonitor().GetVideoMode().Height // 1080
window.SetPos((sw-width)/2, (sh-height)/2)
```
其实在`CreateWindow`的时候窗口就展示了出来，而`SetPos`就会看到窗口移动到中心了，我希望一出来就是中心的，没有移动的过程可以先隐藏，再显示。
```go
glfw.WindowHint(glfw.Visible, glfw.False)

window, err := glfw.CreateWindow(width, height, name, nil, nil)
if err != nil {
    panic(err)
}

sw := glfw.GetPrimaryMonitor().GetVideoMode().Width
sh := glfw.GetPrimaryMonitor().GetVideoMode().Height
window.SetPos((sw-width)/2, (sh-height)/2)

window.Show()
```

- window.MakeContextCurrent()
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

在使用 opengl api 之前需要先设置上下文，可以将当前窗口的上下文绑定为当前上下文。因为新创建的上下文不能被使用，需要绑定为Current Context才能使用。任意时刻，一个线程只能有一个Current Context与之绑定，同时，任意时刻，一个Current Context只能有一个线程与之绑定。

OpenGL自身是一个巨大的状态机(State Machine)：一系列的变量描述OpenGL此刻应当如何运行。OpenGL的状态通常被称为OpenGL上下文(Context)。我们通常使用设置选项和操作缓冲的方式去更改OpenGL状态。最后，我们使用当前OpenGL上下文来渲染。

window，OpenGL上下文，线程之间的关系：
1、每一个window都有一个OpenGL上下文。
2、多个window可以共享一个OpenGL上下文。
3、一个线程同时只能绑定一个OpenGL上下文，作为Current Context。

此处是禁止了窗口缩放，如果允许缩放的话，就会看到视口大小(viewport)和窗口大小不一样了，因此需要监听缩放事件来让OpenGL改变视口大小输出
```go
sizeCallback := func(w *glfw.Window, width int, height int) {
    gl.Viewport(0, 0, int32(width), int32(height))
}
window.SetSizeCallback(sizeCallback)
```

#### 初始化opengl
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

#### 构建VBO，VAO，EBO
VAO和VBO都是用来存储顶点信息的，并把这些信息送入顶点着色器。一个VAO可以对应多个VBO。

在OpenGL程序中，同时只会有一个VAO被绑定到opengl，当然你可以操作完一个之后再绑定到另一个VAO。

VBO是顶点缓冲对象(Vertex Buffer Objects, VBO)，包含了顶点的3d坐标和颜色。但它们是按同类数组存储的，存放在一片显存空间中，程序并不知道这些数字哪个代表3d坐标，哪个代表颜色。

VAO是顶点数组对象(Vertex Array Object, VAO)，用来表示这些数字的第几位分别代表顶点的什么属性。比如这些数字的第1-3位代表3d的xyz坐标，第4-7位代表rbg颜色和透明度。

EBO是元素缓冲对象(Element Buffer Object，EBO)，EBO是一个缓冲区，就像一个VBO一样，它存储 OpenGL 用来决定要绘制哪些顶点的索引,设置顶点的绘制顺序。EBO由VAO进行绑定
![image-20230524113220299](https://videoactivity.bookan.com.cn/ac_202305260902270522126301.png)

如图所示：VAO1中下标为0的指针attribute pointer[0]对应VBO1中的pos[0],表示VOB1数组中下标为0位置代表坐标pos。

VAO2中下标为0的指针attribute pointer[0]对应VBO2中的pos[0]，VAO2中下标为1的指针attribute pointer[1]对应VBO2中的col[0]，表示VOB2数组中下标为1位置代表颜色col。

我们用VBO来存储数据，而用VAO来告诉计算机这些数据分别有什么属性、起什么作用。

VBO是 CPU 和 GPU 之间传递信息的桥梁，我们把数据存入VBO(这一步在CPU执行)，然后VBO会自动把数据送入GPU。但是，对GPU来说，VBO中存的就只是一堆数字而已，要怎么解释它们呢？这就要用到VAO了。

VBO是在显卡存储空间中开辟出的一块内存缓存区，用于存储顶点的各类属性信息，如顶点坐标，顶点法向量，顶点颜色数据等。在渲染时，可以直接从VBO中取出顶点的各类属性数据，由于VBO在显存而不是在内存中，不需要从CPU传输数据，处理效率更高。

顶点坐标为`[]float32`类型，依次为`X, Y, Z`，`窗口中心点为原点，向右为X正，上为Y正`，取值范围`[-1,1]`，
```go
triangle = []float32{
    0, 0.5, 0,
    -0.5, -0.5, 0,
    0.5, -0.5, 0,
}
```

```go
func makeVao(points []float32) uint32 {
	var vbo uint32

    // 在显卡中开辟一块空间，创建顶点缓存对象，个数为1，变量vbo会被赋予一个ID值。
	gl.GenBuffers(1, &vbo)

    // 将 vbo 赋值给 gl.ARRAY_BUFFER，要知道这个对象会被赋予不同的vbo，因此其值是变化的
    // 可选类型：GL_ARRAY_BUFFER, GL_ELEMENT_ARRAY_BUFFER, GL_PIXEL_PACK_BUFFER, GL_PIXEL_UNPACK_BUFFER
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

    // 将内存中的数据传递到显卡中的gl.ARRAY_BUFFER对象上，其实是把数据传递到绑定在其上面的vbo对象上。
    // 4*len(points) 代表总的字节数，因为是32位的
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
    // 创建顶点数组对象，个数为1，变量vao会被赋予一个ID值。
	gl.GenVertexArrays(1, &vao)

    // 后面的两个函数都是要操作具体的vao的，因此需要先将vao绑定到opengl上。
    // 解绑：gl.BindVertexArray(0)，opengl中很多的解绑操作都是传入0
	gl.BindVertexArray(vao)

    // 使vao去引用到gl.ARRAY_BUFFER上面的vbo，这一步完成之后vao就建立了对特定vbo的引用，后面即使gl.ARRAY_BUFFER
    // 的值发生了变化也不影响vao的使用
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
    // 设置 vertex attribute 的状态enabled，默认是disabled
，	gl.EnableVertexAttribArray(0)

	return vao
}
```

- `VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, pointer unsafe.Pointer)`
每一个顶点有会有多个属性，比如常用的有：
1、位置属性，也就是坐标，包括X, Y, Z三个值。
2、颜色属性，如果是RGB，那就是三个值，如果是RGBA，那就是四个值。
3、纹理坐标，S和T，是两个值。
4、其他自定义属性。

我们在定义顶点数据的时候会把这些数据统统注入一个数组中即VBO，比如：
```go
vertices = []float32{
    0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0,
    0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0,
    -0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0,
}
```
那么这堆数据到底有几个点，这取决于你定义了什么属性。比如我可以说它有三个点，每个点有三个属性：1、位置坐标（三个值），2、颜色（三个值），3、纹理坐标（两个值）。

我也可以说它有八个点，给个点有一个属性：1、位置坐标（三个值）。

所以，`VertexAttribPointer`函数的作用是规定了该如何来拆分和使用顶点数据VBO。参数说明如下：

- 1、index：既然属性可能有多个，那就给它标个号吧，从0开始，我们在顶点着色器中使用layout(location = 0)可以把顶点属性的位置值设置为0。
- 2、size：这个属性有几个值，与着色器中的`vec`几对应。
- 3、xtype：属性的值是什么类型，比如`gl.FLOAT`。
- 4、normalized：是否希望数据被归一化，如果为TRUE，那么所有的unsigned数据都会被转换成[0,1]之间，对于signed数据都会被转换成[-1,1]之间，一般选FALSE.
- 5、stride：步长，单位是字节，计算方式是每一个像素点的所有属性总共占用多少字节，比如如果只有位置属性，那就是3个float32，即12个字节，也就是填12。当然，我们也可以填0，让OpenGL自己来算。
- 6、pointer：这个比较难懂，其值为此属性的偏移量，单位字节，还是以`vertices`数据为例，假设它有三个属性：1、位置坐标（三个值），2、颜色（三个值），3、纹理坐标（两个值）。第一个属性偏移量为0，第二个属性偏移量为3个float32为12，第三个属性偏移量为6个float32为24。

为了使用方便，我们使用`VertexAttribPointerWithOffset`来代替`VertexAttribPointer`。比如上面的三个属性可以这么设置

```go
// Position attribute
gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 8*4, 0)
gl.EnableVertexAttribArray(0)
// Color attribute
gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 8*4, 3*4)
gl.EnableVertexAttribArray(1)
// TexCoord attribute
gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, 8*4, 6*4)
gl.EnableVertexAttribArray(2)
```



每一个点虽然设置了几个属性，但是默认情况下，出于性能考虑，这些属性都是不生效的disabled，意味着数据在着色器端是不可见的，哪怕数据已经上传到GPU，所以要手动调用`EnableVertexAttribArray`来逐个让它们生效。其参数index跟`VertexAttribPointer`一样。

我们能声明的顶点属性是有上限的，它一般由硬件来决定。OpenGL确保至少有16个包含4分量的顶点属性可用，但是有些硬件或许允许更多的顶点属性，你可以查询 GL_MAX_VERTEX_ATTRIBS 来获取具体的上限。通常情况下它至少会返回16个，大部分情况下是够用了。

那么，glEnableVertexAttribArray应该在glVertexAttribPointer之前还是之后调用？答案是都可以，只要在绘图调用（glDraw*系列函数）前调用即可。

至此，就完成了当前VAO对当前VBO的引用。

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
着色器中经常看到 vec2,vec3,vec4 等类型，代表有几个分量或者几个数，此处vp是坐标，有x,y,x三个分量。
关键字 in 表明这是输入参数，out为输出参数。

输入参数来自VAO中的一个顶点的全部属性，此处只有坐标属性，因此只有一个变量，如果有多个属性，需要多个in变量。

如果顶点属性中有颜色和纹理属性，那么需要定义out变量，然后out变量会被传给片元着色器的in变量。

#### 片元着色器
片元着色器的作用是计算出每一个像素点最终的颜色，通常片段着色器会包含3D场景的一些额外的数据，如光线，阴影等。

用 RGBA 形式的值通过 vec4 来定义我们图形的颜色。四个分量的值都是`[0, 1]`。

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

程序没有输入变量，因此是固定颜色并输出给下游处理。

#### 画图
OpenGL中所有的图形都是通过分解成三角形的方式进行绘制。

OpenGL的坐标系为右手定则，坐标归一化为[-1,1]。
![右手定则img](https://videoactivity.bookan.com.cn/ac_1_1685591255_546.jpg)

使用VAO作为数据，在opengl program中画图，其中画图和颜色填充交给可编程管线（顶点着色器和片元着色器）来完成，最后呈现在window上。
```go
func draw(vao uint32, window *glfw.Window, prog uint32) {
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

	glfw.PollEvents()

	window.SwapBuffers()
}
```
GLFW与OpenGL之间通过关联上下文来共享内存数据。

- `glfw.PollEvents()`
可以理解为一个消费程序，如果我们设置了一些事件或者回调函数，一旦被触发了就会加入到事件队列中，此函数就会消费并执行对应的回调函数，如果队列为空，此函数就会立即返回，因此此函数要放在循环体中。如果没有调用此函数，会发现窗口提示`无响应`。在一些平台上，窗口的移动，缩放等操作会导致pollevent程序阻塞住，如果有必要的话可以使用`glfwSetWindowRefreshCallback`来重新画，对应的函数为`window.SetRefreshCallback()`

- `gl.ClearColor(1.0, 0.0, 0.0, 1.0)`
给OpenGL上下文中的一个特定对象赋值，此处假设为ColorObj，代表一个RGBA颜色值[0,1]。如果再次调用此函数那么ColorObj的值就会被覆盖，否则其值一直存在。
如果ColorObj被用到了，意味着这一帧里面的所有像素点的RGBA值都是一样的，一个纯色的屏幕，也可以理解为清屏色。

- `gl.ClearDepth(1.0)`
给OpenGL上下文中的一个特定对象赋值，此处假设为DepthObj，代表一个深度值[0,1]。如果再次调用此函数那么DepthObj的值就会被覆盖，否则其值一直存在。
如果DepthObj被用到了，意味着这一帧里面的所有像素点的深度值都是一样的，一个平面。

所谓深度值，就是在3D空间中每个点在Z轴方向上的值，这样就知道哪个物体在上面哪个物体在下面，Z轴是相对于窗口屏幕而言的。

如果你使用了GL_LESS（默认），`gl.DepthFunc(gl.LESS)`，那么Z轴垂直于屏幕向里，也就是深度值小的离人眼近，深度值大的远，也就会被挡住。如果使用的是gl.GREATER那么结果会相反。

默认情况下，深度值为0，也就是Z轴的0点在窗口屏幕上，设置缓冲区深度值的作用在于，将XOY平面做上下平移，看到的东西会不一样。

为了启用深度缓冲区进行深度测试，只需要调用：glEnable（GL_DEPTH_TEST）；另外，即使深度缓冲区未被启用，如果深度缓冲区被创建，OpenGL也会把所有写入到颜色缓冲区的颜色片段对应的深度值写入到深度缓冲区中。但是，如果我们希望在进行深度测试时临时禁止把值写入到深度缓冲区，我们可以使用函数：void glDepthMask（GLboolean mask）；把GL_FALSE作为参数，经禁止写入深度值，但并不禁止用已经写入到深度缓冲区的值进行深度测试。把GL_TRUE作为参数，可以重新启用深度缓冲区的写入。同时，这也是默认的设置。

- `gl.ClearStencil(1.0)`
模板缓冲区可以为屏幕上的每个像素点保存一个无符号整数值。这个值的具体意义视程序的具体应用而定。

在渲染的过程中，可以用这个值与一个预先设定的参考值相比较，根据比较的结果来决定是否更新相应的像素点的颜色值。这个比较的过程被称为模板测试。模板测试发生在透明度测试（alpha test）之后，深度测试（depth test）之前。如果模板测试通过，则相应的像素点更新，否则不更新。就像使用纸板和喷漆一样精确的混图一样，当启动模板测试时，通过模板测试的片段像素点会被替换到颜色缓冲区中，从而显示出来，未通过的则不会保存到颜色缓冲区中，从而达到了过滤的功能。

- `gl.Clear(gl.COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT)`
清除缓冲区数据，包括颜色缓冲(GL_COLOR_BUFFER_BIT)，深度缓冲(GL_DEPTH_BUFFER_BIT)，模板缓冲(GL_STENCIL_BUFFER_BIT)。可以同时传入多个，使用或运算符连接，如果底层显卡支持同时清除就会执行同时清除，如果不支持就会逐个清除。

所谓缓冲区就是暂存数据的，OpenGL在渲染图像的时候会将参数数据放在缓冲区中，其中像素点的颜色值放在颜色缓冲区，像素点的深度值放在深度缓冲区，然后把数据传入显卡来进行绘制，绘制好之后将这一帧的图形放到窗口的`back buffer`处。

如果要展示的是一个静态图片，那显然就只有一帧了，如何才能让图片一直展示在屏幕上呢？通常的方法是，把它看成是有很多帧的数据，但是每一帧都是一样的，因此在for循环中重复画着同一张图片即可。而且静态图片也不会有深度数据。

缓冲区的数据，如果你不主动调用`gl.Clear()`的话，数据会一直存在的，除非被覆盖，因此，合理的方式是在画每一帧之前都进行`gl.Clear()`操作。

`gl.Clear()`在执行的时候先会去查找OpenGL是否有设置ColorObj和DepthObj的值，如果有，就会以它的值来初始化缓冲区，如果没有，就会以默认的黑色和0深度来初始化缓冲区。

但是，调用`gl.Clear()`会触发一次显卡的渲染也就是画图，并且会把这一帧放到`back buffer`处。示例：
```go
var count int
glfw.SwapInterval(10)
func draw(vao uint32, window *glfw.Window, prog uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	count++
	if count%10 == 0 {
		gl.UseProgram(prog)

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3)) // 三角形
		log.Println(count)
	}
	if count >= 100000 {
		count = 0
	}

	glfw.PollEvents()
	window.SwapBuffers()
    log.Println("ok")
}
```
可以看到闪烁的三角形，并且是每个9次才出一个三角形，说明纯黑的屏幕也是在一直输出的。

于是，我们知道，画每一帧的步骤都是：
1、初始化OpenGL缓冲区，对应Clear操作，此处会触发一次画图。
2、调用OpenGL的Draw相关函数进行画图。
3、OpenGL把画好的图放到窗口的`back buffer`处。
4、调用GLFW的SwapBuffers展示`back buffer`上的画面。

- `window.SwapBuffers()`
对应的函数为`glfwSwapBuffers(GLFWwindow *window)`，GLFW的窗口对象有两个缓冲区`front buffer` 和 `back buffer`，当使用OpenGL或者OpenGL ES来渲染窗口的时候，此函数的作用是交换两个缓冲区，具体过程是，`front buffer` 存储这当前帧画面，`back buffer`存储着下一帧画面（如果有的话），切换的过程就是改变了指针的指向而已，于是前就变成了后，后变成了前，如此循环，此窗口必须有OpenGL或者OpenGL ES的上下文，否则会报错。

另外这里还有一个交换频率的设置`glfw.SwapInterval(10)`，也可以理解为刷新频率，正常设置为1即可，默认是0，也就是不停地刷，如果设置了n，就表示会间隔n帧才刷新一次，当然你也可以使用sleep函数自己控制频率。

GLFW的[双缓冲区](https://www.glfw.org/docs/latest/window_guide.html#buffer_swap)的好处是，提高展示效率，避免了边画边展示的尴尬，因为渲染并不是一瞬间就完成的，也是需要时间的，因此提前渲染好，然后直接切换。

这里需要注意的是，虽然GLFW在切换两个缓存，但是并不会清除它们，而是由显卡画新的帧画面来覆盖，如果没有新的画面进来，那就变成了两个缓存画面的交替展现了。

```go
KeyPressAction(window)

glfw.SwapInterval(100)

for !window.ShouldClose() {
    glfw.PollEvents()
    window.SwapBuffers()
    log.Println("ok")
}

func KeyPressAction(window *glfw.Window) {
	keyCallback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if window.GetKey(glfw.KeyR) == glfw.Press {
			log.Println("R")
			gl.ClearColor(1.0, 0.0, 0.0, 1.0)
			gl.Clear(gl.COLOR_BUFFER_BIT)
		}
	}

	window.SetKeyCallback(keyCallback)
}
```

在黑的时候按R之后没有闪屏，在红的时候按R之后依旧闪屏  ?????????

- `gl.PolygonMode(face uint32, mode uint32)`
在调用`gl.DrawArrays`类似函数来画图之前，我们还可以设置多边形模式，其参数`face`是固定的`gl.FRONT_AND_BACK`，而`mode`的取值有三个`gl.POINT, gl.LINE, gl.FILL`，对应的效果分别是画三个点，画三角形的线，画三角形并填充颜色。默认的`mode`是`gl.FILL`。

上面的例子是画一个三角形，如果要画一个正方形，只需要再增加三个顶点即可
```go
square = []float32{
    -0.5, 0.5, 0,
    -0.5, -0.5, 0,
    0.5, -0.5, 0,
    -0.5, 0.5, 0,
    0.5, 0.5, 0,
    0.5, -0.5, 0,
}
```

如果我只想使用四个点能不能画出正方形呢，肯定是可以的，这个时候就要使用到`EBO`对象来指明顶点的索引。
```go
square2 = []float32{
    -0.5, 0.5, 0,
    -0.5, -0.5, 0,
    0.5, -0.5, 0,
    0.5, 0.5, 0,
}

// 索引数据
indexs = []uint32{
    0, 1, 2, // 使用第0,1,2三个顶点来绘制第一个三角形
    0, 2, 3, // 使用第0,2,3三个顶点来绘制第二个三角形
}
```
同时要修改`makVao`方法
```go
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

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}
```
这里要注意了，虽然我们只定义了四个顶点，但是在Draw的时候我们是绘制了6个顶点的，只是重复使用了其中两个顶点，此时我们使用`DrawElements`来绘制。
```go
gl.BindVertexArray(vao)
gl.DrawElements(gl.TRIANGLES, int32(len(indexs)), gl.UNSIGNED_INT, gl.Ptr(indexs))
```


-------------------------------------------------------------------------------------------------------

#### 纹理
**贴图目标**

`GL_TEXTURE_1D、GL_TEXTURE_2D、GL_TEXTURE_3D`

**纹理坐标**

纹理坐标是二维的，表示对应在纹理图片上的哪个点。归一化后[0,1]，左下角为[0,0]，右上角为[1,1]。

我们为每一个顶点设置纹理坐标的属性，就是要将纹理图片贴到我们要绘制的图形上去。

![image-20230602114512537](D:\dev\php\magook\trunk\server\md\img\image-20230602114512537.png)

**纹理环绕**

如果我们的多边形比纹理图片大，或者尺寸比例不一样会出现什么情况呢？我们可以对OpenGL进行设置，以决定当纹理坐标不位于这一区间时应采取的操作。

| 环绕方式           | 描述                                                         |
| ------------------ | ------------------------------------------------------------ |
| GL_REPEAT          | 对纹理的默认行为。重复纹理图像。                             |
| GL_MIRRORED_REPEAT | 和GL_REPEAT一样，但每次重复图片是镜像放置的。                |
| GL_CLAMP_TO_EDGE   | 纹理坐标会被约束在0到1之间，超出的部分会重复纹理坐标的边缘，产生一种边缘被拉伸的效果。 |
| GL_CLAMP_TO_BORDER | 超出的坐标为用户指定的边缘颜色。                             |

我们可以指定两种操作：GL_CLAMP和GL_REPEAT。对于GL_CLAMP,超出纹理坐标的区域会使用纹理图像的边界颜色来代替，如图所示。

![image-20230602110225384](D:\dev\php\magook\trunk\server\md\img\image-20230602110225384.png)

而GL_REPEAT方式则是对纹理坐标进行重置而得到重复的图像。观察图，你就能很容易地发现这一点。

![image-20230602110245840](D:\dev\php\magook\trunk\server\md\img\image-20230602110245840.png)

 ```go
 glTexParameteri(GL_TEXTURE_2D,GL_TEXTURE_WRAP_S,*WrapMode*);//在s方向上的缠绕方式
 glTexParameteri(GL_TEXTURE_2D,GL_TEXTURE_WRAP_T,*WrapMode*);//在t方向上的缠绕方式
 ```

**纹理过滤**

纹理坐标不依赖于分辨率(Resolution)，它可以是任意浮点值，也就是说你定义的纹理坐标，可能不是一个纹理像素点的中心，所以OpenGL需要知道怎样将纹理像素(Texture Pixel，也叫Texel，译注1)映射到纹理坐标。当你有一个很大的物体但是纹理的分辨率很低的时候这就变得很重要了。你可能已经猜到了，OpenGL也有对于纹理过滤(Texture Filtering)的选项。纹理过滤有很多个选项，但是现在我们只讨论最重要的两种：GL_NEAREST和GL_LINEAR。

Texture Pixel也叫Texel，你可以想象你打开一张.jpg格式图片，不断放大你会发现它是由无数像素点组成的，这个点就是纹理像素；注意不要和纹理坐标搞混，纹理坐标是你给模型顶点设置的那个数组，OpenGL以这个顶点的纹理坐标数据去查找纹理图像上的像素，然后进行采样提取纹理像素的颜色。

GL_NEAREST（也叫邻近过滤，Nearest Neighbor Filtering）是OpenGL默认的纹理过滤方式。当设置为GL_NEAREST的时候，OpenGL会选择中心点最接近纹理坐标的那个像素。下图中你可以看到四个像素，加号代表纹理坐标。左上角那个纹理像素的中心距离纹理坐标最近，所以它会被选择为样本颜色：

![image-20230602115805146](D:\dev\php\magook\trunk\server\md\img\image-20230602115805146.png)

GL_LINEAR（也叫线性过滤，(Bi)linear Filtering）它会基于纹理坐标附近的纹理像素，计算出一个插值，近似出这些纹理像素之间的颜色。一个纹理像素的中心距离纹理坐标越近，那么这个纹理像素的颜色对最终的样本颜色的贡献越大。下图中你可以看到返回的颜色是邻近像素的混合色：

![image-20230602115829741](D:\dev\php\magook\trunk\server\md\img\image-20230602115829741.png)

GL_NEAREST产生了颗粒状的图案，我们能够清晰看到组成纹理的像素，而GL_LINEAR能够产生更平滑的图案，很难看出单个的纹理像素。GL_LINEAR可以产生更真实的输出，但有些开发者更喜欢8-bit风格，所以他们会用GL_NEAREST选项。

当进行放大(Magnify)和缩小(Minify)操作的时候可以设置纹理过滤的选项，比如你可以在纹理被缩小的时候使用邻近过滤，被放大时使用线性过滤。我们需要使用glTexParameter*函数为放大和缩小指定过滤方式。这段代码看起来会和纹理环绕方式的设置很相似：

```go
glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_NEAREST);
glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR);
```

可使用的纹理滤镜

| 滤镜                      | 描述                                 |
| ------------------------- | ------------------------------------ |
| GL_NEAREST                | 取最邻近像素                         |
| GL_LINEAR                 | 线性内部插值                         |
| GL_NEAREST_MIPMAP_NEAREST | 最近多贴图等级的最邻近像素           |
| GL_NEAREST_MIPMAP_LINEAR  | 在最近多贴图等级的内部线性插值       |
| GL_LINEAR_MIPMAP_NEAREST  | 在最近多贴图等级的外部线性插值       |
| GL_LINEAR_MIPMAP_LINEAR   | 在最近多贴图等级的外部和内部线性插值 |

### 多级渐远纹理

想象一下，假设我们有一个包含着上千物体的大房间，每个物体上都有纹理。有些物体会很远，但其纹理会拥有与近处物体同样高的分辨率。由于远处的物体可能只产生很少的片段，OpenGL从高分辨率纹理中为这些片段获取正确的颜色值就很困难，因为它需要对一个跨过纹理很大部分的片段只拾取一个纹理颜色。在小物体上这会产生不真实的感觉，更不用说对它们使用高分辨率纹理浪费内存的问题了。

![image-20230602110354834](D:\dev\php\magook\trunk\server\md\img\image-20230602110354834.png)

OpenGL使用一种叫做多级渐远纹理(Mipmap)的概念来解决这个问题，它简单来说就是一系列的纹理图像，后一个纹理图像是前一个的二分之一。多级渐远纹理背后的理念很简单：距观察者的距离超过一定的阈值，OpenGL会使用不同的多级渐远纹理，即最适合物体的距离的那个。由于距离远，解析度不高也不会被用户注意到。同时，多级渐远纹理另一加分之处是它的性能非常好。让我们看一下多级渐远纹理是什么样子的：

![image-20230602141503773](D:\dev\php\magook\trunk\server\md\img\image-20230602141503773.png)

手工为每个纹理图像创建一系列多级渐远纹理很麻烦，幸好OpenGL有一个glGenerateMipmaps函数，在创建完一个纹理后调用它OpenGL就会承担接下来的所有工作了。后面的教程中你会看到该如何使用它。

在渲染中切换多级渐远纹理级别(Level)时，OpenGL在两个不同级别的多级渐远纹理层之间会产生不真实的生硬边界。就像普通的纹理过滤一样，切换多级渐远纹理级别时你也可以在两个不同多级渐远纹理级别之间使用NEAREST和LINEAR过滤。为了指定不同多级渐远纹理级别之间的过滤方式，你可以使用下面四个选项中的一个代替原有的过滤方式：

| 过滤方式                  | 描述                                                         |
| ------------------------- | ------------------------------------------------------------ |
| GL_NEAREST_MIPMAP_NEAREST | 使用最邻近的多级渐远纹理来匹配像素大小，并使用邻近插值进行纹理采样 |
| GL_LINEAR_MIPMAP_NEAREST  | 使用最邻近的多级渐远纹理级别，并使用线性插值进行采样         |
| GL_NEAREST_MIPMAP_LINEAR  | 在两个最匹配像素大小的多级渐远纹理之间进行线性插值，使用邻近插值进行采样 |
| GL_LINEAR_MIPMAP_LINEAR   | 在两个邻近的多级渐远纹理之间使用线性插值，并使用线性插值进行采样 |

就像纹理过滤一样，我们可以使用glTexParameteri将过滤方式设置为前面四种提到的方法之一：

```go
glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_LINEAR_MIPMAP_LINEAR);
glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR);
```

一个常见的错误是，将放大过滤的选项设置为多级渐远纹理过滤选项之一。这样没有任何效果，因为多级渐远纹理主要是使用在纹理被缩小的情况下的：纹理放大不会使用多级渐远纹理，为放大过滤设置多级渐远纹理的选项会产生一个GL_INVALID_ENUM错误代码。

**贴图模式**

取值`GL_MODULATE、GL_DECAL、GL_BLEND`。默认情况下，贴图模式是`GL_MODULATE`，在这种模式下，OpenGL会根据当前的光照系统调整物体的色彩和明暗。第二种模式是GL_DECAL，在这种模式下所有的光照效果都是无效的，OpenGL将仅依据纹理贴图来绘制物体的表面。最后是GL_BLEND，这种模式允许我们使用混合纹理。在这种模式下，我们可以把当前纹理同一个颜色混合而得到一个新的纹理。我们可以调用glTexEnvi函数来设置当前贴图模式：

```go
glTexEnvi(GL_TEXTURE_ENV, GL_TEXTURE_ENV_MODE, *TextureMode*);
```

**创建二维纹理图像**

- `TexImage2D(target uint32, level int32, internalformat int32, width int32, height int32, border int32, format uint32, xtype uint32, pixels unsafe.Pointer)`
  - target ：纹理目标，`GL_TEXTURE_1D、GL_TEXTURE_2D、GL_TEXTURE_3D`
  - level ：指定多级渐远纹理的级别，如果你希望单独手动设置每个多级渐远纹理的级别的话。这里我们填0，也就是基本级别。
  - internalformat：告诉OpenGL我们希望把纹理储存为何种格式，即哪种颜色模型。
  - width，height：纹理的宽度和高度。
  - border：总是被设为`0`（历史遗留的问题）。
  - format：原图的颜色模型。
  - xtype：像素点数据的数据类型。
  - pixels：像素点数据数组的指针。

示例

```go
// 引入 
// _ "image/jpeg"
// _  "image/png"
func MakeTexture(filepath string) uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	imgFile2, _ := os.Open(filepath)
	defer imgFile2.Close()
	img2, _, _ := image.Decode(imgFile2)
	rgba2 := image.NewRGBA(img2.Bounds())
	draw.Draw(rgba2, rgba2.Bounds(), img2, image.Point{0, 0}, draw.Src)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(rgba2.Rect.Size().X), int32(rgba2.Rect.Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba2.Pix))
    
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return texture
}
```

**使用纹理**

纹理是贴在多边形上面的，因此先需要定义一个多边形的顶点数组，每个顶点要定义纹理坐标属性。

```go
width  = 800
height = 600

vertices = []float32{
    // Positions   // Colors      // Texture Coords
    0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, // Top Right
    0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, // Bottom Right
    -0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // Bottom Left
    -0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, // Top Left
}

indices = []uint32{
    0, 1, 3, // First Triangle
    1, 2, 3, // Second Triangle
}
```

步长和偏移信息如下

![image-20230602150240118](D:\dev\php\magook\trunk\server\md\img\image-20230602150240118.png)

于是构建VAO

```go
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
```

顶点着色器，解析顶点属性作为输入参数，并把颜色属性和纹理坐标属性作为输出变量传递给片元着色器。

```go
#version 410

in vec3 vPosition;
in vec3 vColor;
in vec2 vTexCoord;

out vec3 fColor;
out vec2 fTexCoord;

void main() {
    gl_Position = vec4(vPosition, 1.0);
    fColor = vColor;
    fTexCoord = vTexCoord;
}
```

片元着色器，虽然得到了纹理坐标，但是它还需要纹理图片才能得到颜色值，GLSL有一个供纹理对象使用的内建数据类型，叫做采样器(Sampler)，它以纹理类型作为后缀，比如sampler1D、sampler3D，或在我们的例子中的sampler2D。我们可以简单声明一个uniform sampler2D把一个纹理添加到片段着色器中，稍后我们会把纹理赋值给这个uniform。

```go
#version 410

in vec3 fColor;
in vec2 fTexCoord;

out vec4 frag_colour;

uniform sampler2D ourTexture;

void main() {
    frag_colour = texture(ourTexture, fTexCoord);
}
```

我们使用GLSL内建的texture函数来采样纹理的颜色，它第一个参数是纹理采样器，第二个参数是对应的纹理坐标。texture函数会使用之前设置的纹理参数对相应的颜色值进行采样。这个片段着色器的输出就是纹理的（插值）纹理坐标上的（过滤后的）颜色。

```go
func Run() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "texture2d")
	defer glfw.Terminate()

	program := util.InitOpenGL(vertexShaderSource, fragmentShaderSource)
	vao := util.MakeVaoWithEboAndAttrib(vertices, indices)
	pointNum := int32(len(indices))
	texture1 := util.MakeTexture("demo4/container.jpg")

	for !window.ShouldClose() {
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.UseProgram(program)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)

		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, pointNum, gl.UNSIGNED_INT, gl.Ptr(indices))

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
```

![image-20230602154206628](D:\dev\php\magook\trunk\server\md\img\image-20230602154206628.png)

上面的例子中只有一个纹理，且片元着色器中只定义了一个纹理变量，所以我们不需要指定对应关系，如果我们定义了多个纹理变量呢？这意味着，两个纹理需要以一定的比例线性插值的展示出来，也就是混合在一起。我们先修改片元着色器

```go
#version 410

in vec3 fColor;
in vec2 fTexCoord;

out vec4 frag_colour;

uniform sampler2D ourTexture1;
uniform sampler2D ourTexture2;

void main() {
    frag_colour = mix(texture(ourTexture1, fTexCoord), texture(ourTexture2, fTexCoord), 0.2);
}
```

GLSL内建的mix函数需要接受两个值作为参数，并对它们根据第三个参数进行线性插值。。如果第三个值是`0.0`，它会返回第一个输入；如果是`1.0`，会返回第二个输入值。`0.2`会返回`80%`的第一个输入颜色和`20%`的第二个输入颜色，即返回两个纹理的混合色。

然后我们需要指定哪个uniform采样器对应那个纹理。

```go
func Run() {
	runtime.LockOSThread()
	window := util.InitGlfw(width, height, "texture2d")
	defer glfw.Terminate()

	program := util.InitOpenGL(vertexShaderSource, fragmentShaderSource)
	vao := util.MakeVaoWithEboAndAttrib(vertices, indices)
	pointNum := int32(len(indices))
	texture1 := util.MakeTexture("demo4/container.jpg")
	texture2 := util.MakeTexture("demo4/awesomeface.png")

	for !window.ShouldClose() {
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.UseProgram(program)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture1"+"\x00")), 0)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture2"+"\x00")), 1)

		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, pointNum, gl.UNSIGNED_INT, gl.Ptr(indices))

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
```

![image-20230602160414394](D:\dev\php\magook\trunk\server\md\img\image-20230602160414394.png)

可以看到，这个笑脸是上下颠倒了，其实箱子也是颠倒了，这是因为图片的像素扫描是从左上角为原点的，我们使用任何程序读取图片都是这样的，这没有问题，但是OpenGL要求y轴`0`坐标是在图片的底部的，所以我们可以修改一下顶点着色器，使用`1-y`来翻转纹理坐标的y。

```go
#version 410

in vec3 vPosition;
in vec3 vColor;
in vec2 vTexCoord;

out vec3 fColor;
out vec2 fTexCoord;

void main() {
    gl_Position = vec4(vPosition, 1.0);
    fColor = vColor;
	fTexCoord = vec2(vTexCoord.x, 1.0-vTexCoord.y);
}
```

![image-20230602161541893](D:\dev\php\magook\trunk\server\md\img\image-20230602161541893.png)





------



#### 摄像机

------




#### 键盘事件
如果设置了键盘事件，那么在`窗口被聚焦`的时候，按下键盘会触发此事件。

```go
window := util.InitGlfw(width, height, "Conway's Game of Life")
// scancode是一个系统平台相关的键位扫描码信息
// action参数表示这个按键是被按下还是释放，按下的时候会触发action=1，如果不放会一直触发action=2，放开的时候会触发action=0事件
// mods表示是否有Ctrl、Shift、Alt、Super四个按钮的操作，1-shift,2-ctrl,4-alt，8-win
keyCallback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
    log.Printf("key:%d, scancode:%d, action:%d, mods:%v, name:%s\n", key, scancode, action, mods, glfw.GetKeyName(key, scancode))
    // 如果按下了ESC键就关闭窗口
    if key == glfw.KeyEscape && action == glfw.Press {
        window.SetShouldClose(true)
    }
}
// 或者 glfw.GetCurrentContext().SetKeyCallback(keyCallback)
window.SetKeyCallback(keyCallback)
```
取消键盘事件
```go
window.SetKeyCallback(nil)
```
字符输入事件，鼠标聚焦到窗口，然后打开输入法输入
```go
charCallback := func(w *glfw.Window, char rune) {
    log.Printf("char:%s", string(char))
}
window.SetCharCallback(charCallback)

```
```bash
2023/05/29 09:30:48 char:我
2023/05/29 09:30:48 char:们
2023/05/29 09:31:02 char:a
2023/05/29 09:31:02 char:s
2023/05/29 09:31:02 char:d
```

------



#### 鼠标事件

鼠标点击事件沿用了键盘事件，只是将按键变成了左键，右键，滚轮
```go
// 左键：button=0，按下action=1，松开action=0，没有按住事件
// 右键：button=1，按下action=1，松开action=0，没有按住事件
// 滚轮：button=2，按下action=1，松开action=0，没有按住事件
mouseCallback := func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
    log.Printf("button:%d, action:%d, mod:%d\n", button, action, mod)
}
window.SetMouseButtonCallback(mouseCallback)

```

鼠标坐标移动事件
```go
cursorPosCallback := func(w *glfw.Window, xpos float64, ypos float64) {
    // 窗口左上角为 (0, 0)
    log.Printf("x:%f, y:%f", xpos, ypos)
}
window.SetCursorPosCallback(cursorPosCallback)
```

按下左键并拖到鼠标的效果

```go
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
```

滚轮事件？？？？

鼠标滚轮或者触摸板，鼠标滚轮只有yoff，表示垂直滚动了多少，触摸板有xoff和yoff
```go
scrollCallback := func(w *glfw.Window, xoff float64, yoff float64) {
    log.Printf("xoff:%f, yoff:%f", x2, y2)
}
window.SetScrollCallback(scrollCallback)
```

将对象拖拽到窗口放下事件，可以是多选文件，`names`为这些文件的绝对地址。
```go
dropCallback := func(w *glfw.Window, names []string) {
    // names:[D:\dev\php\magook\trunk\server\go-graphic\demo5\square.png]
    log.Printf("names:%v", names)
}
window.SetDropCallback(dropCallback)
```

鼠标聚焦到窗口事件
```go
// CursorEnterCallback is the cursor boundary crossing callback.
type CursorEnterCallback func(w *Window, entered bool)

// SetCursorEnterCallback the cursor boundary crossing callback which is called
// when the cursor enters or leaves the client area of the window.
func (w *Window) SetCursorEnterCallback(cbfun CursorEnterCallback) (previous CursorEnterCallback)
```

操纵杆，控制杆事件
```go
// JoystickCallback is the joystick configuration callback.
type JoystickCallback func(joy, event int)

// SetJoystickCallback sets the joystick configuration callback, or removes the
// currently set callback. This is called when a joystick is connected to or
// disconnected from the system.
func SetJoystickCallback(cbfun JoystickCallback) (previous JoystickCallback)

// JoystickPresent reports whether the specified joystick is present.
func JoystickPresent(joy Joystick) bool

// GetJoystickAxes returns a slice of axis values.
func GetJoystickAxes(joy Joystick) []float32

// GetJoystickButtons returns a slice of button values.
func GetJoystickButtons(joy Joystick) []byte

// GetJoystickName returns the name, encoded as UTF-8, of the specified joystick.
func GetJoystickName(joy Joystick) string
```

------

