package opengl

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

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
	defer ibo.Unbind()

	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*count, nil, gl.DYNAMIC_DRAW)

	return ibo
}

// SetData .
func (v *IBO) SetData(data IBOData) {
	v.Bind()
	defer v.Unbind()

	v.count = data.GetIBOCount()
	gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, 0, 4*int(v.count), data.GetIBOGLPtr())
}

// Bind .
func (v *IBO) Bind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, v.id)
}

// Unbind .
func (v *IBO) Unbind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

// GetCount .
func (v *IBO) GetCount() int32 {
	return v.count
}

// IBOData .
type IBOData interface {
	GetIBOGLPtr() unsafe.Pointer
	GetIBOCount() int32
}
