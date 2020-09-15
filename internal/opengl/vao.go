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
	defer v.Unbind()

	vbo.Bind()
	defer vbo.Unbind()

	layout := vbo.GetLayout()

	offset := 0
	for idx, element := range layout.elements {
		gl.VertexAttribPointer(uint32(idx), element.Count, element.DataType.value, element.Normalized, layout.GetStride(), gl.PtrOffset(offset))
		gl.EnableVertexAttribArray(uint32(idx))
		offset += (element.DataType.size * int(element.Count))
	}
}

// SetIBO .
func (v *VAO) SetIBO(ibo *IBO) {
	v.Bind()
	ibo.Bind()
	v.Unbind()
	ibo.Unbind()

	v.ibo = ibo
}

// GetIBO .
func (v *VAO) GetIBO() *IBO {
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
