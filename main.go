package main

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	err := glfw.Init() // No event handling, so I might not need thread locking
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Visible, glfw.False) // Compute only
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(100, 100, "", nil, nil) // GLFW is required for GL
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	program, err := newProgram(shaderSource)
	if err != nil {
		panic(err)
	}

	var size int32
	gl.GetIntegeri_v(gl.MAX_COMPUTE_WORK_GROUP_SIZE, 0, &size)
	fmt.Println("Max local work group size x: ", size)
	gl.GetIntegeri_v(gl.MAX_COMPUTE_WORK_GROUP_SIZE, 1, &size)
	fmt.Println("Max local work group size y: ", size)
	gl.GetIntegeri_v(gl.MAX_COMPUTE_WORK_GROUP_SIZE, 2, &size)
	fmt.Println("Max local work group size z: ", size)
	gl.GetIntegerv(gl.MAX_COMPUTE_SHARED_MEMORY_SIZE, &size)
	fmt.Println("Shared memory: ", size)

	gl.UseProgram(program)

	app()
}

func newProgram(shaderSource string) (uint32, error) {
	shader, err := compileShader(shaderSource, gl.COMPUTE_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, shader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(shader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
