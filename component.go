package main

import "github.com/go-gl/gl/v4.6-core/gl"

// Component .
type Component struct {
	points []float32
	vao    uint32
}

// NewComponent .
func NewComponent(points []float32) *Component {
	return &Component{points: points, vao: makeVao(points)}
}

// Draw .
func (c *Component) Draw() {
	gl.BindVertexArray(c.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(c.points)/3))

	// Unbinding is optional if we always bind a VAO before a draw call
	// Also, would like to benchmark this
	// It is still safer to unbind so that if someone tries
	// to draw without binding a VAO prior, it fails right away
	gl.BindVertexArray(0)
}
