package main

import (
	"math"
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
		cameraPos   = mgl32.Vec3{0, 0, 3}
		cameraFront = mgl32.Vec3{0, 0, -1}
		cameraUp    = mgl32.Vec3{0, 1, 0}

		firstMouse = true
		// 偏航角被初始化为 -90.0 度，因为 0.0 的偏航角会导致方向向量指向右侧，因此我们最初向左旋转一点。
		yaw   float32 = -90.0
		pitch float32 = 0.0
		lastX float32 = 800.0 / 2.0
		lastY float32 = 600.0 / 2.0
		fov   float32 = 45.0

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

		// 计算旋转角度
		const sensitivity float32 = 0.1 // 鼠标灵敏度
		xOffset *= sensitivity
		yOffset *= sensitivity

		yaw += xOffset
		pitch += yOffset

		// 限制俯仰角度，防止相机翻转
		if pitch > 89.0 {
			pitch = 89.0
		}
		if pitch < -89.0 {
			pitch = -89.0
		}

		front := mgl32.Vec3{
			float32(math.Cos(float64(mgl32.DegToRad(yaw))) * math.Cos(float64(mgl32.DegToRad(pitch)))),
			float32(math.Sin(float64(mgl32.DegToRad(pitch)))),
			float32(math.Sin(float64(mgl32.DegToRad(yaw))) * math.Cos(float64(mgl32.DegToRad(pitch)))),
		}
		cameraFront = front.Normalize()
	})
	// glfw：每当鼠标滚轮滚动时，都会调用此回调
	window.SetScrollCallback(func(w *glfw.Window, x float64, y float64) {
		fov -= float32(y)
		if fov < 1.0 {
			fov = 1.0
		}
		if fov > 45.0 {
			fov = 45.0
		}
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

	sd, err := common.NewShader(`
#version 440 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec2 aTexCoord;

out vec2 TexCoord;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main()
{
	gl_Position = projection * view * model * vec4(aPos, 1.0f);
	TexCoord = vec2(aTexCoord.x, aTexCoord.y);
}`, `
#version 440 core
out vec4 FragColor;

in vec2 TexCoord;

// texture samplers
uniform sampler2D texture1;
uniform sampler2D texture2;

void main()
{
	// linearly interpolate between both textures (80% container, 20% awesomeface)
	FragColor = mix(texture(texture1, TexCoord), texture(texture2, TexCoord), 0.2);
}`)
	if err != nil {
		return err
	}

	vertices := []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,

		-0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
	}
	cubePositions := []mgl32.Vec3{
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

	var vbo, vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// 位置属性
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 5*4, 0)
	gl.EnableVertexAttribArray(0)
	// 纹理坐标属性
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)
	gl.EnableVertexAttribArray(1)

	var texture1, texture2 uint32
	gl.GenTextures(1, &texture1)
	gl.BindTexture(gl.TEXTURE_2D, texture1)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	img, err := common.LoadImgRGB("resource/container.jpg")
	if err != nil {
		return err
	}
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGB,
		int32(img.Width),
		int32(img.Height),
		0,
		gl.RGB,
		gl.UNSIGNED_BYTE,
		gl.Ptr(img.Pixels),
	)
	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.GenTextures(1, &texture2)
	gl.BindTexture(gl.TEXTURE_2D, texture2)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	img, err = common.LoadImgRGB("resource/awesomeface.png", true)
	if err != nil {
		return err
	}
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(img.Width),
		int32(img.Height),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(img.Pixels),
	)
	gl.GenerateMipmap(gl.TEXTURE_2D)

	sd.Use()
	sd.SetInt("texture1", 0)
	sd.SetInt("texture2", 1)

	for !window.ShouldClose() {
		// 每帧时间逻辑
		currentFrame := float32(glfw.GetTime())
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame

		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}
		cameraSpeed := 2.5 * deltaTime
		if window.GetKey(glfw.KeyW) == glfw.Press {
			cameraPos = cameraPos.Add(cameraFront.Mul(cameraSpeed))
		}
		if window.GetKey(glfw.KeyS) == glfw.Press {
			cameraPos = cameraPos.Sub(cameraFront.Mul(cameraSpeed))
		}
		if window.GetKey(glfw.KeyA) == glfw.Press {
			cameraPos = cameraPos.Sub(cameraFront.Cross(cameraUp).Mul(cameraSpeed))
		}
		if window.GetKey(glfw.KeyD) == glfw.Press {
			cameraPos = cameraPos.Add(cameraFront.Cross(cameraUp).Mul(cameraSpeed))
		}

		// 渲染: 清空屏幕为背景颜色
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// 将纹理绑定到相应的纹理单元上
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)

		// 激活着色器
		sd.Use()
		// 将投影矩阵传递给着色器（请注意，在这种情况下，它可能会更改每一帧）
		projection := mgl32.Ident4().
			Mul4(
				mgl32.Perspective(
					mgl32.DegToRad(fov), // 鼠标滚轮进行缩放
					float32(ScreenWidth)/float32(ScreenHeight),
					0.1,
					100.0,
				),
			)
		sd.SetMat("projection", 4, &projection[0])

		// 相机视图变换
		view := mgl32.Ident4().
			Mul4(
				mgl32.LookAtV(
					cameraPos,
					cameraPos.Add(cameraFront),
					cameraUp,
				),
			)
		sd.SetMat("view", 4, &view[0])

		gl.BindVertexArray(vao)
		for i, v := range cubePositions {
			angle := float32(i * 20)
			model := mgl32.Ident4().
				Mul4(
					mgl32.Translate3D(v[0], v[1], v[2]),
				).
				Mul4(
					mgl32.HomogRotate3D(
						mgl32.DegToRad(angle),
						mgl32.Vec3{1, 0.3, 0.5},
					),
				)
			sd.SetMat("model", 4, &model[0])
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

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
