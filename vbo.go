package main

import "github.com/go-gl/gl/v4.6-core/gl"

// VBO .
type VBO struct {
	id uint32
}

// NewVBO .
func NewVBO(vertices []float32) *VBO {
	var vboID uint32
	gl.GenBuffers(1, &vboID)
	vbo := &VBO{id: vboID}

	vbo.Bind()
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	return vbo
}

// Bind .
func (v *VBO) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, v.id)
}

// Unbind .
func (v *VBO) Unbind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
