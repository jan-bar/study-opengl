package common

import (
	"os"

	"github.com/go-gl/gl/v4.4-core/gl"
)

type Shader struct {
	ID uint32
}

func NewShader(vertex, fragment string) (*Shader, error) {
	vertexByte, err := os.ReadFile(vertex)
	if err == nil {
		vertex = string(vertexByte)
	}

	fragmentByte, err := os.ReadFile(fragment)
	if err == nil {
		fragment = string(fragmentByte)
	}

	vertexShader, err := CompileShader(vertex, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}

	fragmentShader, err := CompileShader(fragment, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}

	ID := gl.CreateProgram()

	err = LinkShader(ID, vertexShader, fragmentShader)
	if err != nil {
		return nil, err
	}

	return &Shader{ID: ID}, nil
}

func (s *Shader) Use() {
	gl.UseProgram(s.ID)
}

func (s *Shader) Del() {
	gl.DeleteProgram(s.ID)
}

func (s *Shader) GetUniformLocation(name string) int32 {
	r := gl.GetUniformLocation(s.ID, gl.Str(name+CNull))
	if r < 0 {
		panic("uniform location not found")
	}

	return r
}

func (s *Shader) SetInt(name string, v ...int32) {
	loc := s.GetUniformLocation(name)

	switch len(v) {
	case 1:
		gl.Uniform1i(loc, v[0])
	case 2:
		gl.Uniform2i(loc, v[0], v[1])
	case 3:
		gl.Uniform3i(loc, v[0], v[1], v[2])
	case 4:
		gl.Uniform4i(loc, v[0], v[1], v[2], v[3])
	default:
		panic("unexpected v length")
	}
}

func (s *Shader) SetFloat(name string, v ...float32) {
	loc := s.GetUniformLocation(name)

	switch len(v) {
	case 1:
		gl.Uniform1f(loc, v[0])
	case 2:
		gl.Uniform2f(loc, v[0], v[1])
	case 3:
		gl.Uniform3f(loc, v[0], v[1], v[2])
	case 4:
		gl.Uniform4f(loc, v[0], v[1], v[2], v[3])
	default:
		panic("unexpected v length")
	}
}

func (s *Shader) SetMat(name string, num int, v *float32) {
	loc := s.GetUniformLocation(name)

	switch num {
	case 2:
		gl.UniformMatrix2fv(loc, 1, false, v)
	case 3:
		gl.UniformMatrix3fv(loc, 1, false, v)
	case 4:
		gl.UniformMatrix4fv(loc, 1, false, v)
	default:
		panic("unexpected num")
	}
}
