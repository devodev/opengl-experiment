package opengl

import "github.com/go-gl/gl/v4.6-core/gl"

// VAO .
type VAO struct {
	id  uint32
	ibo *IBO
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
	vbo.Bind()

	layout := vbo.Layout()

	offset := 0
	for idx, element := range layout.elements {
		gl.VertexAttribPointerWithOffset(uint32(idx), element.Count, element.DataType.value, element.Normalized, layout.Stride(), uintptr(offset))
		gl.EnableVertexAttribArray(uint32(idx))
		offset += (element.DataType.size * int(element.Count))
	}

	v.Unbind()
	vbo.Unbind()
}

// SetIBO .
func (v *VAO) SetIBO(ibo *IBO) {
	v.Bind()
	ibo.Bind()
	v.Unbind()
	ibo.Unbind()

	v.ibo = ibo
}

// IBO .
func (v *VAO) IBO() *IBO {
	return v.ibo
}

// Bind .
func (v *VAO) Bind() {
	gl.BindVertexArray(v.id)
}

// Unbind .
func (v *VAO) Unbind() {
	gl.BindVertexArray(0)
}
