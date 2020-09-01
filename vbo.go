package main

import "github.com/go-gl/gl/v4.6-core/gl"

// VBO .
type VBO struct {
	id       uint32
	offset   int
	elements []*VBOElementFloat
}

// NewVBO .
func NewVBO() *VBO {
	var vboID uint32
	gl.GenBuffers(1, &vboID)
	vbo := &VBO{id: vboID, offset: 0}

	return vbo
}

// AddElement .
func (v *VBO) AddElement(data []float32, count int32, eType uint32, normalized bool) {
	element := &VBOElementFloat{
		count:      count,
		eType:      eType,
		stride:     int32(count * 4), // TODO: replace 4 with lookup using eType
		offset:     v.offset,
		normalized: normalized,
		data:       data,
	}
	v.elements = append(v.elements, element)
	v.offset += 4 * len(data)

	v.Bind()
	defer v.Unbind()

	// TODO: replace 4 with lookup using eType
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(data), gl.Ptr(data), gl.STATIC_DRAW)
}

// Bind .
func (v *VBO) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, v.id)
}

// Unbind .
func (v *VBO) Unbind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

// VBOElementFloat .
type VBOElementFloat struct {
	count      int32
	eType      uint32
	stride     int32
	offset     int
	normalized bool
	data       []float32
}
