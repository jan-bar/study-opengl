#include "common.h"

// 教程地址: https://learnopengl-cn.github.io/01%20Getting%20started/04%20Hello%20Triangle/

void framebuffer_size_callback(GLFWwindow *window, int width, int height);
void processInput(GLFWwindow *window);

const unsigned int SCR_WIDTH = 800;
const unsigned int SCR_HEIGHT = 600;

const char *vertexShaderSource = "#version 460 core\n"
								 "layout (location = 0) in vec3 aPos;\n"
								 "void main()\n"
								 "{\n"
								 "   gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);\n"
								 "}\0";
const char *fragmentShaderSource = "#version 460 core\n"
								   "out vec4 FragColor;\n"
								   "void main()\n"
								   "{\n"
								   "   FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);\n"
								   "}\n\0";

int main()
{
	// glfw: initialize and configure
	// ------------------------------
	glfwInit();
	glfwWindowHint(GLFW_CONTEXT_VERSION_MAJOR, 4);
	glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, 6);
	glfwWindowHint(GLFW_OPENGL_PROFILE, GLFW_OPENGL_CORE_PROFILE);

#ifdef __APPLE__
	glfwWindowHint(GLFW_OPENGL_FORWARD_COMPAT, GL_TRUE);
#endif

	// glfw window creation
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

	// glad: load all OpenGL function pointers
	// ---------------------------------------
	if (!gladLoadGLLoader((GLADloadproc)glfwGetProcAddress))
	{
		std::cout << "Failed to initialize GLAD" << std::endl;
		return -1;
	}

	// 构建并编译我们的着色器程序
	// ------------------------------------
	// 顶点着色器
	unsigned int vertexShader = glCreateShader(GL_VERTEX_SHADER);
	int success = CompileShader(vertexShaderSource, GL_VERTEX_SHADER, &vertexShader);
	if (!success)
	{
		return success;
	}

	// 片段着色器
	unsigned int fragmentShader = glCreateShader(GL_FRAGMENT_SHADER);
	success = CompileShader(fragmentShaderSource, GL_FRAGMENT_SHADER, &fragmentShader);
	if (!success)
	{
		return success;
	}

	// 链接着色器
	unsigned int shaderProgram = glCreateProgram();
	success = LinkShader(shaderProgram, {vertexShader, fragmentShader});
	if (!success)
	{
		return success;
	}

	// 设置顶点数据（和缓冲区）并配置顶点属性
	// X: [-1(左),1(右)]
	// Y: [-1(下),1(上)]
	float vertices[] = {
		// 第一个三角形
		0.5f, 0.5f, 0.0f,  // 右上角
		0.5f, -0.5f, 0.0f, // 右下角
		-0.5f, 0.5f, 0.0f, // 左上角
		// 第二个三角形
		0.5f, -0.5f, 0.0f,	// 右下角
		-0.5f, -0.5f, 0.0f, // 左下角
		-0.5f, 0.5f, 0.0f,	// 左上角
		// 第三个三角形
		0.9f, -0.5f, 0.0f,	// 右下角
		0.5f, -0.5f, 0.0f, // 左下角
		0.8f, 0.5f, 0.0f	// 左上角
	};

	unsigned int VBO, VAO;
	glGenVertexArrays(1, &VAO);
	glGenBuffers(1, &VBO);
	// 首先绑定顶点数组对象，然后绑定并设置顶点缓冲区，最后配置顶点属性
	glBindVertexArray(VAO);

	// 绑定顶点缓冲类型
	glBindBuffer(GL_ARRAY_BUFFER, VBO);
	// 把之前定义的顶点数据复制到缓冲的内存中
	// GL_STATIC_DRAW ：数据不会或几乎不会改变
	// GL_DYNAMIC_DRAW：数据会被改变很多
	// GL_STREAM_DRAW ：数据每次绘制时都会改变
	glBufferData(GL_ARRAY_BUFFER, sizeof(vertices), vertices, GL_STATIC_DRAW);

	glVertexAttribPointer(0, 3, GL_FLOAT, GL_FALSE, 3 * sizeof(float), (void *)0);
	glEnableVertexAttribArray(0);

	// 请注意，这是允许的，对 glVertexAttribPointer 的调用将 VBO 注册为顶点属性的绑定顶点缓冲区对象，因此之后我们可以安全地解除绑定
	glBindBuffer(GL_ARRAY_BUFFER, 0);

	// 之后您可以取消绑定 VAO，这样其他 VAO 调用就不会意外修改此 VAO，但这种情况很少发生。修改其他
	// 无论如何，VAO 都需要调用 glBindVertexArray，因此当不是直接需要时，我们通常不会取消绑定 VAO（或 VBO）
	glBindVertexArray(0);

	// 取消注释此调用以绘制线框多边形
	glPolygonMode(GL_FRONT_AND_BACK, GL_LINE);

	// render loop
	// -----------
	while (!glfwWindowShouldClose(window))
	{
		// input
		// -----
		processInput(window);

		// render
		// ------
		glClearColor(0.2f, 0.3f, 0.3f, 1.0f);
		glClear(GL_COLOR_BUFFER_BIT);

		// 画出我们的第一个三角形
		glUseProgram(shaderProgram);
		glBindVertexArray(VAO); // 由于我们只有一个 VAO，因此无需每次都绑定它，但我们这样做是为了让事情更有条理
		glDrawArrays(GL_TRIANGLES, 0,
					 // 每3个元素算1个顶点, (元素个数 / 3) = 顶点个数
					 std::size(vertices) / 3);
		// glBindVertexArray(0); // 无需每次都解绑

		// glfw: swap buffers and poll IO events (keys pressed/released, mouse moved etc.)
		// -------------------------------------------------------------------------------
		glfwSwapBuffers(window);
		glfwPollEvents();
	}

	// 可选：一旦超出其用途，就取消分配所有资源
	// ------------------------------------------------------------------------
	glDeleteVertexArrays(1, &VAO);
	glDeleteBuffers(1, &VBO);
	glDeleteProgram(shaderProgram);

	// glfw: terminate, clearing all previously allocated GLFW resources.
	// ------------------------------------------------------------------
	glfwTerminate();
	return 0;
}

// process all input: query GLFW whether relevant keys are pressed/released this frame and react accordingly
// ---------------------------------------------------------------------------------------------------------
void processInput(GLFWwindow *window)
{
	if (glfwGetKey(window, GLFW_KEY_ESCAPE) == GLFW_PRESS)
		glfwSetWindowShouldClose(window, true);
}

// glfw: whenever the window size changed (by OS or user resize) this callback function executes
// ---------------------------------------------------------------------------------------------
void framebuffer_size_callback(GLFWwindow *window, int width, int height)
{
	// make sure the viewport matches the new window dimensions; note that width and
	// height will be significantly larger than specified on retina displays.
	glViewport(0, 0, width, height);
}
