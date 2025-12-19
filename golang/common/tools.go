package common

import (
	"bytes"
	"fmt"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
)

func CompileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	shaderSource, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, shaderSource, nil)
	free()
	gl.CompileShader(shader)

	var tmp int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &tmp)
	if tmp == gl.FALSE {
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &tmp)
		info := make([]byte, tmp)
		gl.GetShaderInfoLog(shader, tmp, nil, unsafe.SliceData(info))
		return 0, fmt.Errorf("%s: %s", source, bytes.TrimRight(info, "\x00"))
	}

	return shader, nil
}

func LinkShader(program uint32, shaders ...uint32) error {
	for _, v := range shaders {
		gl.AttachShader(program, v)
	}
	gl.LinkProgram(program)

	var tmp int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &tmp)
	if tmp == gl.FALSE {
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &tmp)
		info := make([]byte, tmp)
		gl.GetProgramInfoLog(program, tmp, nil, unsafe.SliceData(info))
		return fmt.Errorf("LinkShader: %s", bytes.TrimRight(info, "\x00"))
	}

	for _, v := range shaders {
		gl.DeleteShader(v)
	}

	return nil
}
