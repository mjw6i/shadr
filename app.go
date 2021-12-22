package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
)

const shaderSource = `
	#version 460
	layout(local_size_x = 16, local_size_y = 1, local_size_z = 1) in;
	layout(std430, binding = 3) buffer layoutName
	{
		uint data_SSBO[];
	};
	void compare_and_swap(uint a, uint b) {
		if (data_SSBO[a] > data_SSBO[b] + 1) {
			int max = data_SSBO[a] / 20;
			int diff = (data_SSBO[a] - data_SSBO[b]) / 2;
			int t = (diff > max) ? max : diff;
			data_SSBO[a] -= t;
			data_SSBO[b] += t;
		}
	}
	void main() {
		if (gl_LocalInvocationID.x >= 1 ) {
			compare_and_swap(gl_GlobalInvocationID.x, gl_GlobalInvocationID.x - 1);
		}
		if (gl_LocalInvocationID.x < 16 - 1 ) {
			compare_and_swap(gl_GlobalInvocationID.x, gl_GlobalInvocationID.x + 1);
		}
	}
` + "\x00"

const uint32size = 4

func app() {
	frame := []uint32{
		58, 83, 43, 72,
		06, 18, 34, 90,
		54, 89, 16, 44,
		57, 53, 73, 87,
	}
	fmt.Println(frame, sum(frame))

	var ssbo uint32
	gl.GenBuffers(1, &ssbo)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, ssbo)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(frame)*uint32size, gl.Ptr(frame), gl.DYNAMIC_READ) // DYNAMIC_READ is expected usage for optimisation purposes
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 3, ssbo)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0) // unbind (wiki example, code works without defer here)

	for i := 1; i <= 50; i++ {
		gl.DispatchCompute(uint32(len(frame)), 1, 1)
		gl.GetNamedBufferSubData(ssbo, 0, len(frame)*uint32size, gl.Ptr(frame))
		fmt.Println(frame, sum(frame))
	}
}

func sum(arr []uint32) int {
	s := 0
	for _, e := range arr {
		s += int(e)
	}

	return s
}
