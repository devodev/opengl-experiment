package renderer

import (
	"unsafe"

	"github.com/devodev/opengl-experiment/internal/opengl"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	maxQuads    = 10000
	maxVertices = maxQuads * 4
	maxIndices  = maxQuads * 6

	maxTextures = 32

	quadVertexSize = int(unsafe.Sizeof(QuadVertex{}))

	quadVertices = []mgl32.Vec4{
		{-0.5, 0.5, 0.0, 1.0},
		{-0.5, -0.5, 0.0, 1.0},
		{0.5, -0.5, 0.0, 1.0},
		{0.5, 0.5, 0.0, 1.0},
	}
	quadTexCoords = []mgl32.Vec2{
		{0, 1},
		{0, 0},
		{1, 0},
		{1, 1},
	}
	quadLayout = opengl.NewVBOLayout(
		opengl.VBOLayoutElement{Count: 4, Normalized: false, DataType: opengl.GLDataTypeFloat},
		opengl.VBOLayoutElement{Count: 2, Normalized: false, DataType: opengl.GLDataTypeFloat},
		opengl.VBOLayoutElement{Count: 1, Normalized: false, DataType: opengl.GLDataTypeFloat},
	)
)

type TexturedQuad struct {
	Transform mgl32.Mat4
	Texture   opengl.Texture
}

type Quad struct {
	// quad-related rendering primitives
	vao           *opengl.VAO
	vbo           *opengl.VBO
	shaderProgram *opengl.ShaderProgram
	// quad-related batch rendering data
	data *quadData
}

func (q *Quad) Init() error {
	// initialize quad-related rendering primitives
	quadVertexShaderSource := string(append([]byte(quadVertexShader), byte('\x00')))
	quadFragmentShaderSource := string(append([]byte(quadFragmentShader), byte('\x00')))
	shaderProgram, err := opengl.NewShaderProgram(quadVertexShaderSource, quadFragmentShaderSource)
	if err != nil {
		return err
	}

	// create VBO
	vbo, err := opengl.NewVBO(maxVertices * quadVertexSize)
	if err != nil {
		return err
	}
	vbo.SetLayout(quadLayout)

	// create IBO
	ibo := opengl.NewIBO(maxIndices)

	// create VAO and set buffers on it
	vao := opengl.NewVAO()
	vao.AddVBO(vbo)
	vao.SetIBO(ibo)

	// setup textures
	samplers := make([]int32, maxTextures)
	for i := 0; i < maxTextures; i++ {
		samplers[i] = int32(i)
	}
	shaderProgram.SetUniform1iv("tex", int32(len(samplers)), &samplers[0])

	q.vao = vao
	q.vbo = vbo
	q.shaderProgram = shaderProgram
	q.data = newQuadData()

	return nil
}

func (q *Quad) setViewProjectionMatrix(vp mgl32.Mat4) {
	q.shaderProgram.SetUniformMatrix4fv("vp", 1, false, &vp[0])
}

func (q *Quad) Begin(vp mgl32.Mat4) {
	// reset data each frame
	q.data = newQuadData()
	q.setViewProjectionMatrix(vp)
}

func (q *Quad) End() {
	q.vbo.SetData(q.data)
	q.vao.IBO().SetData(q.data)

	for _, t := range q.data.Textures {
		t.Bind()
	}
	q.shaderProgram.Bind()
	q.vao.Bind()

	// actual draw call
	count := q.vao.IBO().Count()
	gl.DrawElements(gl.TRIANGLES, int32(count), gl.UNSIGNED_INT, nil)

	for _, t := range q.data.Textures {
		t.Unbind()
	}
	q.vao.Unbind()
	q.shaderProgram.Unbind()
}

func (q *Quad) AddTextured(quad *TexturedQuad) error {
	return q.data.AddTextured(quad)
}

// QuadVertex .
type QuadVertex struct {
	Position mgl32.Vec4
	TexCoord mgl32.Vec2
	TexIndex float32
}

// quadData .
type quadData struct {
	Textures map[int]opengl.Texture
	Vertices []QuadVertex
	Indices  []uint32
}

func newQuadData() *quadData {
	return &quadData{
		Textures: make(map[int]opengl.Texture),
		Vertices: make([]QuadVertex, 0, maxVertices),
		Indices:  make([]uint32, 0, maxIndices),
	}
}

// AddTextured .
func (d *quadData) AddTextured(quad *TexturedQuad) error {
	if err := d.addTexture(quad.Texture); err != nil {
		return err
	}

	// add indices
	offset := len(d.Vertices)
	d.Indices = append(d.Indices,
		uint32(offset),
		uint32(offset+1),
		uint32(offset+2),
		uint32(offset+2),
		uint32(offset+3),
		uint32(offset),
	)

	// add vertices
	for i := 0; i < len(quadVertices); i++ {
		vertex := QuadVertex{
			Position: quad.Transform.Mul4x1(quadVertices[i]),
			TexCoord: quadTexCoords[i],
			TexIndex: float32(quad.Texture.Index()),
		}
		d.Vertices = append(d.Vertices, vertex)
	}

	return nil
}

func (d *quadData) addTexture(texture opengl.Texture) error {
	// noop if already registered
	if _, ok := d.Textures[texture.Index()]; ok {
		return nil
	}
	d.Textures[texture.Index()] = texture
	return nil
}

// VBOGLPtr implements the VBOData interface.
func (d *quadData) VBOGLPtr() unsafe.Pointer {
	return gl.Ptr(d.Vertices)
}

// VBOSize implements the VBOData interface.
func (d *quadData) VBOSize() int {
	return quadVertexSize * len(d.Vertices)
}

// IBOGLPtr implements the IBOData interface.
func (d *quadData) IBOGLPtr() unsafe.Pointer {
	return gl.Ptr(d.Indices)
}

// IBOCount  implements the IBOData interface.
func (d *quadData) IBOCount() int32 {
	return int32(len(d.Indices))
}
