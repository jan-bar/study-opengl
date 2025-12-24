#include "common.h"

void framebuffer_size_callback(GLFWwindow *window, int width, int height);
void processInput(GLFWwindow *window);

const unsigned int SCR_WIDTH = 800;
const unsigned int SCR_HEIGHT = 600;

const char *vertexShaderSource = "#version 460 core\n"
								 "layout (location = 0) in vec3 aPos;\n"
								 "layout (location = 1) in vec3 aColor;\n"
								 "out vec3 ourColor;\n"
								 "void main()\n"
								 "{\n"
								 "   gl_Position = vec4(aPos, 1.0);\n"
								 "   ourColor = aColor;\n"
								 "}\0";

const char *fragmentShaderSource = "#version 460 core\n"
								   "out vec4 FragColor;\n"
								   "in vec3 ourColor;\n"
								   "void main()\n"
								   "{\n"
								   "   FragColor = vec4(ourColor, 1.0f);\n"
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
	float vertices[] = {
		// positions         // colors
		0.5f, -0.5f, 0.0f, 1.0f, 0.0f, 0.0f,  // bottom right
		-0.5f, -0.5f, 0.0f, 0.0f, 1.0f, 0.0f, // bottom left
		0.0f, 0.5f, 0.0f, 0.0f, 0.0f, 1.0f	  // top
	};

	unsigned int VBO, VAO;
	glGenVertexArrays(1, &VAO);
	glGenBuffers(1, &VBO);
	// 首先绑定顶点数组对象，然后绑定并设置顶点缓冲区，最后配置顶点属性。
	glBindVertexArray(VAO);

	glBindBuffer(GL_ARRAY_BUFFER, VBO);
	glBufferData(GL_ARRAY_BUFFER, sizeof(vertices), vertices, GL_STATIC_DRAW);

	// 位置属性
	glVertexAttribPointer(
		0, // layout (location = 0)
		3,
		GL_FLOAT,
		GL_FALSE,
		6 * sizeof(float), // 每次步长6
		(void *)0);		   // 0偏移的3个元素
	glEnableVertexAttribArray(0);
	// 颜色属性
	glVertexAttribPointer(
		1, // layout (location = 1)
		3,
		GL_FLOAT,
		GL_FALSE,
		6 * sizeof(float),			  // 每次步长6
		(void *)(3 * sizeof(float))); // 3偏移的3个元素
	glEnableVertexAttribArray(1);

	// 之后您可以取消绑定 VAO，这样其他 VAO 调用就不会意外修改此 VAO，但这种情况很少发生。修改其他
	// 无论如何，VAO 都需要调用 glBindVertexArray，因此当不是直接需要时，我们通常不会取消绑定 VAO（或 VBO）。
	// glBindVertexArray(0);

	// 由于我们只有一个着色器，如果我们愿意，我们也可以提前激活着色器一次
	glUseProgram(shaderProgram);

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

		// 渲染三角形
		glBindVertexArray(VAO);
		glDrawArrays(GL_TRIANGLES, 0, 3);

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
