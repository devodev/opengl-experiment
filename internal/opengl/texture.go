package opengl

import (
	"fmt"
	"image"
	"os"

	"image/draw"
	// need to initialize each image type
	// that could be used in NewTexture
	_ "image/png"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// Texture .
type Texture struct {
	id          uint32
	textureUnit uint32
}

// NewTexture .
func NewTexture(filepath string, textureUnit int) (*Texture, error) {
	rgba, err := rgbaFromFile(filepath)
	if err != nil {
		return nil, err
	}
	if textureUnit < 0 {
		return nil, fmt.Errorf("texture target out of bounds: %d != 0 <= x", textureUnit)
	}

	var id uint32
	gl.GenTextures(1, &id)

	texture := &Texture{
		id:          id,
		textureUnit: uint32(gl.TEXTURE0 + textureUnit),
	}
	texture.Bind()
	defer texture.Unbind()

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)
	gl.GenerateMipmap(gl.TEXTURE_2D)
	return texture, nil
}

// Bind .
func (t *Texture) Bind() {
	gl.ActiveTexture(t.textureUnit)
	gl.BindTexture(gl.TEXTURE_2D, t.id)
}

// Unbind .
func (t *Texture) Unbind() {
	gl.ActiveTexture(t.textureUnit)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// GetTextureUnit .
func (t *Texture) GetTextureUnit() uint32 {
	return t.textureUnit
}

func rgbaFromFile(filepath string) (*image.RGBA, error) {
	reader, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error reading texture file: %s", err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("error decoding texture file: %s", err)
	}
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("error creating texture rgba: unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	return rgba, nil
}
