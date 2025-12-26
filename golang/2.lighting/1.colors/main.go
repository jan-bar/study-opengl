package main

import (
	"runtime"

	"opengl/common"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// https://learnopengl-cn.github.io/01%20Getting%20started/09%20Camera/

func main() {
	runtime.LockOSThread()

	err := HelloTriangle()
	if err != nil {
		panic(err)
	}
}

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

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

	window, err := glfw.CreateWindow(ScreenWidth, ScreenHeight, "LearnOpenGL", nil, nil)
	if err != nil {
		return err
	}

	window.MakeContextCurrent()
	window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		// 确保视口与新窗口尺寸匹配；请注意宽度和
		// 高度将明显大于视网膜显示器上指定的高度
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	var (
		camera = common.NewCamera(
			common.WithPosition(mgl32.Vec3{0, 0, 3}),
		)

		firstMouse         = true
		lastX      float32 = ScreenWidth / 2.0
		lastY      float32 = ScreenHeight / 2.0

		deltaTime float32 = 0
		lastFrame float32 = 0
	)
	// glfw：每当鼠标移动时，都会调用此回调
	window.SetCursorPosCallback(func(w *glfw.Window, x float64, y float64) {
		xPos := float32(x)
		yPos := float32(y)

		if firstMouse {
			lastX = xPos
			lastY = yPos
			firstMouse = false
		}

		// 计算当前光标位置与上次位置的偏移量
		xOffset := xPos - lastX
		yOffset := lastY - yPos // 注意：y 轴是从下到上的
		lastX = xPos
		lastY = yPos

		camera.ProcessMouseMovement(xOffset, yOffset)
	})
	// glfw：每当鼠标滚轮滚动时，都会调用此回调
	window.SetScrollCallback(func(w *glfw.Window, x float64, y float64) {
		camera.ProcessMouseScroll(float32(y))
	})
	// 告诉 GLFW 捕获我们的鼠标
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	// 加载所有 OpenGL 方法
	err = gl.Init()
	if err != nil {
		return err
	}

	// 配置全局 opengl 状态
	gl.Enable(gl.DEPTH_TEST) // 启用深度测试

	lightingShader, err := common.NewShader(`
#version 440 core
layout (location = 0) in vec3 aPos;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main()
{
	gl_Position = projection * view * model * vec4(aPos, 1.0);
}`, `
#version 440 core
out vec4 FragColor;
  
uniform vec3 objectColor;
uniform vec3 lightColor;

void main()
{
    FragColor = vec4(lightColor * objectColor, 1.0);
}`)
	if err != nil {
		return err
	}

	lightCubeShader, err := common.NewShader(`
#version 440 core
layout (location = 0) in vec3 aPos;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main()
{
	gl_Position = projection * view * model * vec4(aPos, 1.0);
}`, `
#version 440 core
out vec4 FragColor;

void main()
{
    FragColor = vec4(1.0); // set all 4 vector values to 1.0
}`)
	if err != nil {
		return err
	}

	vertices := []float32{
		-0.5, -0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, 0.5, -0.5,
		0.5, 0.5, -0.5,
		-0.5, 0.5, -0.5,
		-0.5, -0.5, -0.5,

		-0.5, -0.5, 0.5,
		0.5, -0.5, 0.5,
		0.5, 0.5, 0.5,
		0.5, 0.5, 0.5,
		-0.5, 0.5, 0.5,
		-0.5, -0.5, 0.5,

		-0.5, 0.5, 0.5,
		-0.5, 0.5, -0.5,
		-0.5, -0.5, -0.5,
		-0.5, -0.5, -0.5,
		-0.5, -0.5, 0.5,
		-0.5, 0.5, 0.5,

		0.5, 0.5, 0.5,
		0.5, 0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, -0.5, 0.5,
		0.5, 0.5, 0.5,

		-0.5, -0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, -0.5, 0.5,
		0.5, -0.5, 0.5,
		-0.5, -0.5, 0.5,
		-0.5, -0.5, -0.5,

		-0.5, 0.5, -0.5,
		0.5, 0.5, -0.5,
		0.5, 0.5, 0.5,
		0.5, 0.5, 0.5,
		-0.5, 0.5, 0.5,
		-0.5, 0.5, -0.5,
	}

	var vbo, cubeVao uint32
	gl.GenVertexArrays(1, &cubeVao)
	gl.GenBuffers(1, &vbo)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindVertexArray(cubeVao)

	// 位置属性
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4, 0)
	gl.EnableVertexAttribArray(0)

	var lightCubeVao uint32
	gl.GenVertexArrays(1, &lightCubeVao)
	gl.BindVertexArray(lightCubeVao)

	// 我们只需要绑定到VBO（将其与glVertexAttribPointer链接），不需要填充它； VBO 的数据已经包含我们需要的所有内容（它已经绑定，但我们出于教育目的再次绑定）
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	// 位置属性
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4, 0)
	gl.EnableVertexAttribArray(0)

	lightPos := mgl32.Vec3{1.2, 1.0, 2.0}

	for !window.ShouldClose() {
		// 每帧时间逻辑
		currentFrame := float32(glfw.GetTime())
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame

		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}
		if window.GetKey(glfw.KeyW) == glfw.Press {
			camera.ProcessKeyboard(common.ForWard, deltaTime)
		}
		if window.GetKey(glfw.KeyS) == glfw.Press {
			camera.ProcessKeyboard(common.BackWard, deltaTime)
		}
		if window.GetKey(glfw.KeyA) == glfw.Press {
			camera.ProcessKeyboard(common.Left, deltaTime)
		}
		if window.GetKey(glfw.KeyD) == glfw.Press {
			camera.ProcessKeyboard(common.Right, deltaTime)
		}

		// 渲染: 清空屏幕为背景颜色
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		lightingShader.Use()
		lightingShader.SetFloat("objectColor", 1.0, 0.5, 0.31)
		lightingShader.SetFloat("lightColor", 1.0, 1.0, 1.0)

		projection := mgl32.Ident4().
			Mul4(
				mgl32.Perspective(
					mgl32.DegToRad(camera.Zoom), // 鼠标滚轮进行缩放
					float32(ScreenWidth)/float32(ScreenHeight),
					0.1,
					100.0,
				),
			)
		lightingShader.SetMat("projection", 4, &projection[0])

		view := camera.GetViewMatrix()
		lightingShader.SetMat("view", 4, &view[0])

		model := mgl32.Ident4()
		lightingShader.SetMat("model", 4, &model[0])

		gl.BindVertexArray(cubeVao)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		lightCubeShader.Use()
		lightCubeShader.SetMat("projection", 4, &projection[0])
		lightCubeShader.SetMat("view", 4, &view[0])

		model = mgl32.Ident4().
			Mul4(
				mgl32.Translate3D(lightPos[0], lightPos[1], lightPos[2]),
			).
			Mul4(
				mgl32.Scale3D(0.2, 0.2, 0.2),
			)
		lightCubeShader.SetMat("model", 4, &model[0])

		gl.BindVertexArray(lightCubeVao)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		// glfw: 交换缓冲区和轮询 IO 事件（按键按下、释放、鼠标移动等）
		window.SwapBuffers()
		glfw.PollEvents()
	}

	// 可选：一旦超出其用途，就取消分配所有资源
	// ------------------------------------------------------------------------
	gl.DeleteVertexArrays(1, &cubeVao)
	gl.DeleteVertexArrays(1, &lightCubeVao)
	gl.DeleteBuffers(1, &vbo)

	return nil
}
