package main

import (
	"runtime"

	"opengl/common"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// https://learnopengl-cn.github.io/01%20Getting%20started/08%20Coordinate%20Systems/

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

	// 设置顶点数据（和缓冲区）并配置顶点属性
	vertices := []float32{
		// positions          // texture coords
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

	// 位置属性
	gl.VertexAttribPointerWithOffset(
		// layout (location = 0)
		0,
		// 顶点属性的大小, vec3
		3,
		// 数据的类型
		gl.FLOAT,
		false,
		5*4,
		0,
	)
	gl.EnableVertexAttribArray(0)
	// 纹理坐标属性
	gl.VertexAttribPointerWithOffset(
		// layout (location = 1)
		1,
		// 顶点属性的大小, vec3
		2,
		// 数据的类型
		gl.FLOAT,
		false,
		5*4,
		3*4,
	)
	gl.EnableVertexAttribArray(1)

	// 加载并创建纹理
	var texture1, texture2 uint32
	gl.GenTextures(1, &texture1)
	// 所有即将进行的 GL_TEXTURE_2D 操作现在都会对该纹理对象产生影响
	gl.BindTexture(gl.TEXTURE_2D, texture1)
	// 设置纹理包裹参数
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	// 设置纹理过滤参数
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	// 加载图像、创建纹理并生成 mipmap
	img, err := common.LoadImgRGB("resource/container.jpg")
	if err != nil {
		return err
	}
	// 创建纹理
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
	// 生成 mipmap
	gl.GenerateMipmap(gl.TEXTURE_2D)

	// 加载并创建第二个纹理
	gl.GenTextures(1, &texture2)
	// 所有即将进行的 GL_TEXTURE_2D 操作现在都会对该纹理对象产生影响
	gl.BindTexture(gl.TEXTURE_2D, texture2)
	// 设置纹理包裹参数
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	// 设置纹理过滤参数
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	// 加载图像、创建纹理并生成 mipmap
	img, err = common.LoadImgRGB("resource/awesomeface.png", true)
	if err != nil {
		return err
	}
	// 创建纹理
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
	// 生成 mipmap
	gl.GenerateMipmap(gl.TEXTURE_2D)

	// 告诉 opengl 每个采样器属于哪个纹理单元（只需执行一次）
	sd.Use() // 在设置制服之前不要忘记激活 - 使用着色器！
	sd.SetInt("texture1", 0)
	sd.SetInt("texture2", 1)

	// 创建纹理
	for !window.ShouldClose() {
		// 处理所有输入：查询GLFW是否按下-释放此帧相关按键并做出相应反应
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}

		// 渲染: 清空屏幕为背景颜色
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // 现在还要清除深度缓冲区！

		// 将纹理绑定到相应的纹理单元上
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)

		// 激活着色器
		sd.Use()

		model := mgl32.Ident4().
			Mul4(
				mgl32.HomogRotate3D(
					float32(glfw.GetTime()),
					mgl32.Vec3{0.5, 1, 0},
				),
			)
		view := mgl32.Ident4().
			Mul4(
				mgl32.Translate3D(0, 0, -3),
			)
		projection := mgl32.Ident4().
			Mul4(
				mgl32.Perspective(
					mgl32.DegToRad(45),
					float32(ScreenWidth)/float32(ScreenHeight),
					0.1,
					100.0,
				),
			)

		sd.SetMat("model", 4, &model[0])
		sd.SetMat("view", 4, &view[0])
		sd.SetMat("projection", 4, &projection[0])

		// 渲染容器
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, 36) // 36 个顶点

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
