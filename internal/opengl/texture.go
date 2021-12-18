package opengl

import (
	"fmt"
	"image"
	"os"

	// need to initialize each image type
	// that could be used in NewTexture
	_ "image/png"

	"github.com/disintegration/imaging"

	"github.com/go-gl/gl/v4.6-core/gl"
)

var (
	maxTextures  uint32 = 32
	textureCount uint32 = 0
)

type Texture interface {
	BinderUnbinder
	Index() int
	Unit() uint32
}

type texture struct {
	id    uint32
	index uint32
	unit  uint32
}

// Newtexture .
func NewNRGBATexture(filepath string) (*texture, error) {
	rgba, err := rgbaFromFile(filepath)
	if err != nil {
		return nil, err
	}

	if textureCount >= maxTextures {
		return nil, fmt.Errorf("max texture count reached: %d", maxTextures)
	}
	defer func() {
		textureCount += 1
	}()

	// TODO: handle opengl texture registration errors
	var id uint32
	gl.GenTextures(1, &id)

	texture := &texture{
		id:    id,
		index: textureCount,
		unit:  uint32(gl.TEXTURE0 + textureCount),
	}
	texture.setFromNRGBA(rgba)

	return texture, nil
}

// Bind implements the Binder interface.
func (t *texture) Bind() {
	//fmt.Printf("BIND [index: %v, unit: %v, id: %v]\n", t.index, t.unit, t.id)
	gl.ActiveTexture(t.unit)
	gl.BindTexture(gl.TEXTURE_2D, t.id)
}

// Unbind implements the Unbinder interface.
func (t *texture) Unbind() {
	//fmt.Printf("UNBIND [index: %v, unit: %v, id: %v]\n", t.index, t.unit, t.id)
	gl.ActiveTexture(t.unit)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// ID .
func (t *texture) ID() uint32 {
	return t.id
}

// Index .
func (t *texture) Index() int {
	return int(t.index)
}

// Unit .
func (t *texture) Unit() uint32 {
	return t.unit
}

func (t *texture) setFromNRGBA(data *image.NRGBA) {
	t.Bind()

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(data.Rect.Size().X),
		int32(data.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(data.Pix),
	)
	gl.GenerateMipmap(gl.TEXTURE_2D)

	t.Unbind()
}

func rgbaFromFile(filepath string) (*image.NRGBA, error) {
	reader, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error reading texture file: %s", err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("error decoding texture file: %s", err)
	}

	// Replaced manually drawing image.Image into image.RGBA
	// with disintegration/imaging lib, which provide convenience methods
	// for flipping/transposing/etc.
	nrgba := imaging.FlipV(img)
	//
	// rgba := image.NewRGBA(img.Bounds())
	// if rgba.Stride != rgba.Rect.Size().X*4 {
	// 	return nil, fmt.Errorf("error creating texture rgba: unsupported stride")
	// }
	// draw.Draw(rgba, img.Bounds(), img, image.Point{0, 0}, draw.Src)

	return nrgba, nil
}
