@echo off
if exist "%~1" (
	echo Build %1% ...
	g++ -g -std=c++17 -I./glad/include -I./glfw/include -I./include -L./glfw/lib-mingw-w64 glad/src/glad.c %1% -lglfw3dll -o main.exe
) else (
	echo Skip Build ...
)
.\main.exe
