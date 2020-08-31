package main

import "github.com/go-gl/gl/v4.6-core/gl"

// IBO .
type IBO struct {
	id    uint32
	count int32
}

// NewIBO .
func NewIBO(indices []uint32) *IBO {
	var iboID uint32
	gl.GenBuffers(1, &iboID)
	ibo := &IBO{id: iboID, count: int32(len(indices))}

	ibo.Bind()
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

	return ibo
}

// Bind .
func (v *IBO) Bind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, v.id)
}

// Unbind .
func (v *IBO) Unbind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}
