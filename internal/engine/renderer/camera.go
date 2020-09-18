package renderer

import (
	"github.com/go-gl/mathgl/mgl32"
)

var (
	defaultCameraPerspectiveFov = mgl32.DegToRad(45.0)
	defaultCameraZoomLevel      = float32(1.0)
	defaultCameraNear           = float32(0.1)
	defaultCameraFar            = float32(10.0)
)

// Camera .
type Camera interface {
	Resize(int, int)
	GetProjectionMatrix() mgl32.Mat4
}

// CameraPerspective .
type CameraPerspective struct {
	width       int
	height      int
	aspectRatio float32
	fov         float32
	near        float32
	far         float32

	projectionMatrix mgl32.Mat4
}

// NewCameraPerspective .
func NewCameraPerspective(width, height int) *CameraPerspective {
	camera := &CameraPerspective{
		fov:  defaultCameraPerspectiveFov,
		near: defaultCameraNear,
		far:  defaultCameraFar,
	}
	camera.Resize(width, height)
	return camera
}

// Resize .
func (c *CameraPerspective) Resize(width, height int) {
	c.width = width
	c.height = height
	c.aspectRatio = float32(c.width) / float32(c.height)
	c.recalculateProjectionMatrix()
}

// GetProjectionMatrix .
func (c *CameraPerspective) GetProjectionMatrix() mgl32.Mat4 {
	return c.projectionMatrix
}

func (c *CameraPerspective) recalculateProjectionMatrix() {
	c.projectionMatrix = mgl32.Perspective(c.fov, c.aspectRatio, c.near, c.far)
}

// CameraOrthographic .
type CameraOrthographic struct {
	width       int
	height      int
	aspectRatio float32
	zoomLevel   float32
	near        float32
	far         float32

	projectionMatrix mgl32.Mat4
}

// NewCameraOrthographic .
func NewCameraOrthographic(width, height int) *CameraOrthographic {
	camera := &CameraOrthographic{
		zoomLevel: defaultCameraZoomLevel,
		near:      defaultCameraNear,
		far:       defaultCameraFar,
	}
	camera.Resize(width, height)
	return camera
}

// Resize .
func (c *CameraOrthographic) Resize(width, height int) {
	c.width = width
	c.height = height
	c.aspectRatio = float32(c.width) / float32(c.height)
	c.recalculateProjectionMatrix()
}

// GetProjectionMatrix .
func (c *CameraOrthographic) GetProjectionMatrix() mgl32.Mat4 {
	return c.projectionMatrix
}

func (c *CameraOrthographic) recalculateProjectionMatrix() {
	c.projectionMatrix = mgl32.Ortho((-c.aspectRatio)*c.zoomLevel, c.aspectRatio*c.zoomLevel, -c.zoomLevel, c.zoomLevel, c.near, c.far)
}
