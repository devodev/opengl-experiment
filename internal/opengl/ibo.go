package opengl

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// IBOData .
type IBOData interface {
	IBOGLPtr() unsafe.Pointer
	IBOCount() int32
}

// IBO .
type IBO struct {
	id    uint32
	count int32
}

// NewIBO .
func NewIBO(count int) *IBO {
	var iboID uint32
	gl.GenBuffers(1, &iboID)

	ibo := &IBO{id: iboID, count: int32(count)}

	ibo.Bind()
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*count, nil, gl.DYNAMIC_DRAW)
	ibo.Unbind()

	return ibo
}

// SetData .
func (v *IBO) SetData(data IBOData) {
	v.count = data.IBOCount()

	v.Bind()
	gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, 0, 4*int(v.count), data.IBOGLPtr())
	v.Unbind()
}

// Bind .
func (v *IBO) Bind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, v.id)
}

// Unbind .
func (v *IBO) Unbind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

// Count .
func (v *IBO) Count() int32 {
	return v.count
}
