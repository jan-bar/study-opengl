package main

import (
	"fmt"
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
layout (location = 2) in float aFace;

flat out int Face;
out vec2 TexCoord;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main()
{
	gl_Position = projection * view * model * vec4(aPos, 1.0f);
	TexCoord = vec2(aTexCoord.x, 1.0 - aTexCoord.y);
	Face = int(aFace);
}`, `
#version 440 core
out vec4 FragColor;

in vec2 TexCoord;

uniform sampler2DArray diceTex;
flat in int Face;

void main()
{
	FragColor = texture(diceTex, vec3(TexCoord, Face));
}`)
	if err != nil {
		return err
	}

	// 每个面有2个三角形,对应6个顶点,2个纹理坐标,1个面索引
	vertices := []float32{
		// positions    // texture coords   // face index
		-0.5, -0.5, -0.5, 0.0, 0.0, 0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0,
		0.5, 0.5, -0.5, 1.0, 1.0, 0,
		0.5, 0.5, -0.5, 1.0, 1.0, 0,
		-0.5, 0.5, -0.5, 0.0, 1.0, 0,
		-0.5, -0.5, -0.5, 0.0, 0.0, 0,

		-0.5, -0.5, 0.5, 0.0, 0.0, 1,
		0.5, -0.5, 0.5, 1.0, 0.0, 1,
		0.5, 0.5, 0.5, 1.0, 1.0, 1,
		0.5, 0.5, 0.5, 1.0, 1.0, 1,
		-0.5, 0.5, 0.5, 0.0, 1.0, 1,
		-0.5, -0.5, 0.5, 0.0, 0.0, 1,

		-0.5, 0.5, 0.5, 1.0, 0.0, 2,
		-0.5, 0.5, -0.5, 1.0, 1.0, 2,
		-0.5, -0.5, -0.5, 0.0, 1.0, 2,
		-0.5, -0.5, -0.5, 0.0, 1.0, 2,
		-0.5, -0.5, 0.5, 0.0, 0.0, 2,
		-0.5, 0.5, 0.5, 1.0, 0.0, 2,

		0.5, 0.5, 0.5, 1.0, 0.0, 3,
		0.5, 0.5, -0.5, 1.0, 1.0, 3,
		0.5, -0.5, -0.5, 0.0, 1.0, 3,
		0.5, -0.5, -0.5, 0.0, 1.0, 3,
		0.5, -0.5, 0.5, 0.0, 0.0, 3,
		0.5, 0.5, 0.5, 1.0, 0.0, 3,

		-0.5, -0.5, -0.5, 0.0, 1.0, 4,
		0.5, -0.5, -0.5, 1.0, 1.0, 4,
		0.5, -0.5, 0.5, 1.0, 0.0, 4,
		0.5, -0.5, 0.5, 1.0, 0.0, 4,
		-0.5, -0.5, 0.5, 0.0, 0.0, 4,
		-0.5, -0.5, -0.5, 0.0, 1.0, 4,

		-0.5, 0.5, -0.5, 0.0, 1.0, 5,
		0.5, 0.5, -0.5, 1.0, 1.0, 5,
		0.5, 0.5, 0.5, 1.0, 0.0, 5,
		0.5, 0.5, 0.5, 1.0, 0.0, 5,
		-0.5, 0.5, 0.5, 0.0, 0.0, 5,
		-0.5, 0.5, -0.5, 0.0, 1.0, 5,
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
		0,
		3,
		gl.FLOAT,
		false,
		6*4,
		0,
	)
	gl.EnableVertexAttribArray(0)
	// 纹理坐标属性
	gl.VertexAttribPointerWithOffset(
		1,
		2,
		gl.FLOAT,
		false,
		6*4,
		3*4,
	)
	gl.EnableVertexAttribArray(1)
	// 每个面的索引
	gl.VertexAttribPointerWithOffset(
		2,
		1,
		gl.FLOAT,
		false,
		6*4,
		5*4,
	)
	gl.EnableVertexAttribArray(2)

	// 创建一个纹理数组
	var texArray uint32
	gl.GenTextures(1, &texArray)
	gl.BindTexture(gl.TEXTURE_2D_ARRAY, texArray)

	gl.TexImage3D(
		gl.TEXTURE_2D_ARRAY,
		0,
		gl.RGB8,
		512, 512, 6, // 指定宽度、高度、层数(6个面)
		0,
		gl.RGB,
		gl.UNSIGNED_BYTE,
		nil,
	)
	for i := 0; i < 6; i++ {
		img, err := common.LoadImgRGB(fmt.Sprintf("resource/dice-%d.png", i+1))
		if err != nil {
			return err
		}
		// 循环上传每个面的纹理数据
		gl.TexSubImage3D(
			gl.TEXTURE_2D_ARRAY,
			0,
			0, 0, int32(i),
			int32(img.Width), int32(img.Height), 1,
			gl.RGB,
			gl.UNSIGNED_BYTE,
			gl.Ptr(img.Pixels),
		)
	}
	// 设置纹理参数
	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	// 绑定纹理数组到纹理单元0
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D_ARRAY, texArray)
	// 将纹理单元0绑定到着色器中的"diceTex"采样器
	sd.Use()
	sd.SetInt("diceTex", 0)

	for !window.ShouldClose() {
		// 处理所有输入：查询GLFW是否按下-释放此帧相关按键并做出相应反应
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}

		// 渲染: 清空屏幕为背景颜色
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // 现在还要清除深度缓冲区！

		// 激活着色器
		sd.Use()

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

		sd.SetMat("view", 4, &view[0])
		sd.SetMat("projection", 4, &projection[0])

		// 渲染容器
		gl.BindVertexArray(vao)
		for i, pos := range cubePositions {
			model := mgl32.Ident4().
				Mul4(
					mgl32.Translate3D(pos[0], pos[1], pos[2]),
				)
			angle := 20.0 * float32(i)
			if i%3 == 0 {
				angle = float32(glfw.GetTime()) * 25
			}
			model = model.Mul4(
				mgl32.HomogRotate3D(
					mgl32.DegToRad(angle),
					mgl32.Vec3{1, 0.3, 0.5},
				),
			)
			sd.SetMat("model", 4, &model[0])

			// 一次性绘制6个面的三角形
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
