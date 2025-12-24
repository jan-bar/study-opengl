package main

import (
	"math"
	"runtime"

	"opengl/common"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	runtime.LockOSThread()

	err := HelloTriangle()
	if err != nil {
		panic(err)
	}
}

func HelloTriangle() error {
	err := glfw.Init()
	if err != nil {
		return err
	}

	defer glfw.Terminate() // glfw: 终止，清除所有先前分配的 GLFW 资源

	// glfw: 初始化配置,设置gl版本
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 4)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	// Mac OS X 需要如下配置
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(800, 600, "LearnOpenGL", nil, nil)
	if err != nil {
		return err
	}

	window.MakeContextCurrent()
	window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		// 确保视口与新窗口尺寸匹配；请注意宽度和
		// 高度将明显大于视网膜显示器上指定的高度
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	// 加载所有 OpenGL 方法
	err = gl.Init()
	if err != nil {
		return err
	}

	sd, err := common.NewShader(`
#version 440 core
layout (location = 0) in vec3 aPos;
void main()
{
	gl_Position = vec4(aPos, 1.0);
}`, `
#version 440 core
out vec4 FragColor;
uniform vec4 ourColor;
void main()
{
	FragColor = ourColor;
}`)
	if err != nil {
		return err
	}

	vertices := []float32{
		0.5, -0.5, 0.0, // bottom right
		-0.5, -0.5, 0.0, // bottom left
		0.0, 0.5, 0.0, // top
	}

	// 创建顶点数组对象(VAO)
	var vbo, vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	// 首先绑定顶点数组对象，然后绑定并设置顶点缓冲区，最后配置顶点属性
	gl.BindVertexArray(vao)

	// 绑定顶点缓冲类型
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	// 把之前定义的顶点数据复制到缓冲的内存中
	gl.BufferData(gl.ARRAY_BUFFER,
		// float32 每个元素占用 4 字节
		len(vertices)*4,
		gl.Ptr(vertices),
		// GL_STATIC_DRAW ：数据不会或几乎不会改变
		// GL_DYNAMIC_DRAW：数据会被改变很多
		// GL_STREAM_DRAW ：数据每次绘制时都会改变
		gl.STATIC_DRAW,
	)

	// 设置顶点属性指针
	gl.VertexAttribPointer(
		// 要配置的顶点属性,layout(location = 0)
		0,
		// 顶点属性的大小, vec3
		3,
		// 数据的类型
		gl.FLOAT,
		// 是否希望数据被标准化,
		// 如果设置为 true，所有数据都会被映射到 0（或 -1 用于符号数据）到 1 之间
		false,
		// 步长设置为3 * sizeof(float)
		3*4,
		nil,
	)
	gl.EnableVertexAttribArray(0)

	// 之后您可以取消绑定 VAO，这样其他 VAO 调用就不会意外修改此 VAO，但这种情况很少发生。修改其他
	// 无论如何，VAO 都需要调用 glBindVertexArray，因此当不是直接需要时，我们通常不会取消绑定 VAO（或 VBO）
	// gl.BindVertexArray(0);

	// 绑定 VAO（它已经绑定，但只是为了演示）：因为我们只有一个 VAO，所以我们可以
	// 只需在渲染相应三角形之前预先绑定它即可；这是另一种方法。
	gl.BindVertexArray(vao)

	for !window.ShouldClose() {
		// 处理所有输入：查询GLFW是否按下-释放此帧相关按键并做出相应反应
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}

		// 渲染: 清空屏幕为背景颜色
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// 确保在调用 glUniform 之前激活着色器
		sd.Use()

		timeVal := glfw.GetTime()
		greenVal := float32(math.Sin(timeVal)/2.0 + 0.5)
		sd.SetFloat("ourColor", 0.0, greenVal, 0.0, 1.0)

		// 渲染三角形
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		// glfw: 交换缓冲区和轮询 IO 事件（按键按下、释放、鼠标移动等）
		window.SwapBuffers()
		glfw.PollEvents()
	}

	// 可选：一旦超出其用途，就取消分配所有资源
	// ------------------------------------------------------------------------
	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)
	sd.Del()

	return nil
}
