package common

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
)

const CNull = "\x00"

func CompileShader(source string, shaderType uint32) (uint32, error) {
	if !strings.HasSuffix(source, CNull) {
		source += CNull
	}

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
		return 0, fmt.Errorf("%s: %s", source, bytes.TrimRight(info, CNull))
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
		return fmt.Errorf("LinkShader: %s", bytes.TrimRight(info, CNull))
	}

	for _, v := range shaders {
		gl.DeleteShader(v)
	}

	return nil
}

type ImageData struct {
	Width  int
	Height int
	Pixels []uint8
}

func LoadImgRGB(path string, rgba ...bool) (*ImageData, error) {
	fr, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(fr)
	_ = fr.Close()

	if err != nil {
		return nil, err
	}

	var (
		i      = 0
		rect   = img.Bounds()
		width  = rect.Dx()
		height = rect.Dy()
		pixels []uint8
		setPix func(r, g, b, a uint32)
	)

	if len(rgba) > 0 && rgba[0] {
		pixels = make([]uint8, width*height*4)
		setPix = func(r, g, b, a uint32) {
			pixels[i] = uint8(r >> 8)
			pixels[i+1] = uint8(g >> 8)
			pixels[i+2] = uint8(b >> 8)
			pixels[i+3] = uint8(a >> 8)
			i += 4
		}
	} else {
		pixels = make([]uint8, width*height*3)
		setPix = func(r, g, b, a uint32) {
			pixels[i] = uint8(r >> 8)
			pixels[i+1] = uint8(g >> 8)
			pixels[i+2] = uint8(b >> 8)
			i += 3
		}
	}

	// Go image: 原点在左上,OpenGL: 原点在左下,需要从下往上遍历
	for y := rect.Max.Y - 1; y >= rect.Min.Y; y-- {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			setPix(img.At(x, y).RGBA())
		}
	}

	return &ImageData{
		Width:  width,
		Height: height,
		Pixels: pixels,
	}, nil
}
