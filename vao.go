package main

import "github.com/go-gl/gl/v4.6-core/gl"

// VAO .
type VAO struct {
	id uint32
}

// NewVAO .
func NewVAO() *VAO {
	var vaoID uint32
	gl.GenVertexArrays(1, &vaoID)
	return &VAO{id: vaoID}
}

// AddVBO .
func (v *VAO) AddVBO(vbo *VBO) {
	v.Bind()
	defer v.Unbind()

	vbo.Bind()
	defer vbo.Unbind()

	for idx, element := range vbo.elements {
		// VertexAttribPointer index refers to `layout (location = 0) ` in the vertex shader
		// stride can be set to 0 when the values are tightly packed
		gl.VertexAttribPointer(uint32(idx), element.count, element.eType, element.normalized, element.stride, gl.PtrOffset(element.offset))
		gl.EnableVertexAttribArray(uint32(idx))
	}
}

// Bind .
func (v *VAO) Bind() {
	gl.BindVertexArray(v.id)
}

// Unbind .
func (v *VAO) Unbind() {
	gl.BindVertexArray(0)
}
