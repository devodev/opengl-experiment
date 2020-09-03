package opengl

import (
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
)

var (
	glDataTypes = map[uint32]*GlDataType{
		gl.FLOAT: {
			name:  "FLOAT",
			size:  4,
			value: gl.FLOAT,
		},
	}
)

// GlDataType .
type GlDataType struct {
	name  string
	size  int
	value uint32
}

// GetGlDataType .
func GetGlDataType(dataType uint32) (*GlDataType, error) {
	glDataType, ok := glDataTypes[dataType]
	if !ok {
		return nil, fmt.Errorf("invalid VBO data type of value: %v", dataType)
	}
	return glDataType, nil
}

// VBO .
type VBO struct {
	id       uint32
	stride   int32
	dataType *GlDataType
	data     interface{}
	elements []*VBOElement
}

// NewVBO .
func NewVBO(data interface{}, dataType uint32) (*VBO, error) {
	glDataType, err := GetGlDataType(dataType)
	if err != nil {
		return nil, err
	}

	var vboID uint32
	gl.GenBuffers(1, &vboID)
	vbo := &VBO{id: vboID, dataType: glDataType, data: data}

	vbo.Bind()
	defer vbo.Unbind()

	// hardcoded as float32 for now until I figure out
	// how to handle this correctly.
	realData := data.([]float32)
	gl.BufferData(gl.ARRAY_BUFFER, glDataType.size*len(realData), gl.Ptr(realData), gl.STATIC_DRAW)

	return vbo, nil
}

// AddElement .
func (v *VBO) AddElement(count int32, normalized bool) {
	v.elements = append(v.elements, NewVBOElement(count, normalized))
}

// GetStride .
func (v *VBO) GetStride() int32 {
	stride := int32(0)
	for _, element := range v.elements {
		stride += (int32(v.dataType.size) * element.count)
	}
	return stride
}

// Bind .
func (v *VBO) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, v.id)
}

// Unbind .
func (v *VBO) Unbind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

// VBOElement .
type VBOElement struct {
	count      int32
	normalized bool
}

// NewVBOElement .
func NewVBOElement(count int32, normalized bool) *VBOElement {
	return &VBOElement{
		count:      count,
		normalized: normalized,
	}
}
