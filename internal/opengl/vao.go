package opengl

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

	offset := 0
	for idx, element := range vbo.elements {
		gl.VertexAttribPointer(uint32(idx), element.count, vbo.dataType.value, element.normalized, vbo.GetStride(), gl.PtrOffset(offset))
		gl.EnableVertexAttribArray(uint32(idx))
		offset += (vbo.dataType.size * int(element.count))
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
