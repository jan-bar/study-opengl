package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	runtime.LockOSThread()

	err := HelloWindow()
	if err != nil {
		panic(err)
	}
}

func HelloWindow() error {
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

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Printf("OpenGL version: %s\n", version)

	for !window.ShouldClose() {
		// 处理所有输入：查询GLFW是否按下-释放此帧相关按键并做出相应反应
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}

		// 渲染: 清空屏幕为背景颜色
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// glfw: 交换缓冲区和轮询 IO 事件（按键按下、释放、鼠标移动等）
		window.SwapBuffers()
		glfw.PollEvents()
	}

	return nil
}
