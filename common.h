#ifndef __common_h_
#define __common_h_

#include <glad/glad.h>
#include <GLFW/glfw3.h>

#include <iostream>
#include <vector>

inline int CompileShader(const char *source, GLenum shaderType, unsigned int *shader)
{
	*shader = glCreateShader(shaderType);
	glShaderSource(*shader, 1, &source, NULL);
	glCompileShader(*shader);

	int success;
	glGetShaderiv(*shader, GL_COMPILE_STATUS, &success);
	if (!success)
	{
		int length;
		glGetShaderiv(*shader, GL_INFO_LOG_LENGTH, &length);
		std::vector<char> infoLog(length);
		glGetShaderInfoLog(*shader, length, NULL, infoLog.data());
		std::cerr << source << std::endl
				  << infoLog.data() << std::endl;
	}
	return success;
}

inline int LinkShader(unsigned int program, std::initializer_list<unsigned int> shaders)
{
	for (unsigned int v : shaders)
	{
		glAttachShader(program, v);
	}
	glLinkProgram(program);

	int success;
	glGetProgramiv(program, GL_LINK_STATUS, &success);
	if (!success)
	{
		int length;
		glGetProgramiv(program, GL_INFO_LOG_LENGTH, &length);
		std::vector<char> infoLog(length);
		glGetProgramInfoLog(program, length, NULL, infoLog.data());
		std::cerr << "LinkShader:" << std::endl
				  << infoLog.data() << std::endl;
	}

	for (unsigned int v : shaders)
	{
		glDeleteShader(v);
	}
	return success;
}

#endif
