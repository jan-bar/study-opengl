#include <glad/glad.h>
#include <GLFW/glfw3.h>

#include <iostream>

// 教程地址: https://learnopengl-cn.github.io/01%20Getting%20started/03%20Hello%20Window/

void framebuffer_size_callback(GLFWwindow *window, int width, int height);
void processInput(GLFWwindow *window);

// 设置窗口的宽高
const unsigned int SCR_WIDTH = 800;
const unsigned int SCR_HEIGHT = 600;

int main()
{
	// glfw: 初始化配置,设置gl版本
	// ------------------------------
	glfwInit();
	glfwWindowHint(GLFW_CONTEXT_VERSION_MAJOR, 4);
	glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, 6);
	glfwWindowHint(GLFW_OPENGL_PROFILE, GLFW_OPENGL_CORE_PROFILE);

#ifdef __APPLE__
	glfwWindowHint(GLFW_OPENGL_FORWARD_COMPAT, GL_TRUE);
#endif

	// glfw 创建一个窗口
	// --------------------
	GLFWwindow *window = glfwCreateWindow(SCR_WIDTH, SCR_HEIGHT, "LearnOpenGL", NULL, NULL);
	if (window == NULL)
	{
		std::cout << "Failed to create GLFW window" << std::endl;
		glfwTerminate();
		return -1;
	}
	glfwMakeContextCurrent(window);
	glfwSetFramebufferSizeCallback(window, framebuffer_size_callback);

	// glad: 加载所有OpenGL方法
	// ---------------------------------------
	if (!gladLoadGLLoader((GLADloadproc)glfwGetProcAddress))
	{
		std::cout << "Failed to initialize GLAD" << std::endl;
		return -1;
	}

	// 打印加载的 OpenGL 版本
	std::cout << "OpenGL version: " << glGetString(GL_VERSION) << std::endl;

	// 循环渲染
	// -----------
	while (!glfwWindowShouldClose(window))
	{
		// 键盘输入
		// -----
		processInput(window);

		// 渲染
		// ------
		glClearColor(0.2f, 0.3f, 0.3f, 1.0f);
		glClear(GL_COLOR_BUFFER_BIT);

		// glfw: 交换缓冲区和轮询 IO 事件（按键按下、释放、鼠标移动等）
		// -------------------------------------------------------------------------------
		glfwSwapBuffers(window);
		glfwPollEvents();
	}

	// glfw: 终止，清除所有先前分配的 GLFW 资源。
	// ------------------------------------------------------------------
	glfwTerminate();
	return 0;
}

// 处理所有输入：查询GLFW是否按下-释放此帧相关按键并做出相应反应
// ---------------------------------------------------------------------------------------------------------
void processInput(GLFWwindow *window)
{
	if (glfwGetKey(window, GLFW_KEY_ESCAPE) == GLFW_PRESS)
		glfwSetWindowShouldClose(window, true);
}

// glfw：每当窗口大小发生变化（由操作系统或用户调整大小）时，都会执行此回调函数
// ---------------------------------------------------------------------------------------------
void framebuffer_size_callback(GLFWwindow *window, int width, int height)
{
	// 确保视口与新窗口尺寸匹配；请注意宽度和
	// 高度将明显大于视网膜显示器上指定的高度。
	glViewport(0, 0, width, height);
}
