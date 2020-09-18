package renderer

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	defaultControllerPos       = mgl32.Vec3{0, 0, 2}
	defaultControllerTarget    = mgl32.Vec3{0, 0, -0.5}
	defaultControllerUp        = mgl32.Vec3{0, 1, 0}
	defaultControllerBaseSpeed = float32(2)
)

// CameraController .
type CameraController struct {
	pos       mgl32.Vec3
	target    mgl32.Vec3
	up        mgl32.Vec3
	baseSpeed float32

	viewMatrix mgl32.Mat4

	camera Camera
}

// NewCameraController .
func NewCameraController(camera Camera) *CameraController {
	return &CameraController{
		pos:       defaultControllerPos,
		target:    defaultControllerTarget,
		up:        defaultControllerUp,
		baseSpeed: defaultControllerBaseSpeed,
		camera:    camera,
	}
}

// OnUpdate .
func (c *CameraController) OnUpdate(glfwWindow *glfw.Window, deltaTime float64) {
	if glfwWindow.GetAttrib(glfw.Focused) == glfw.False {
		return
	}

	speed := c.baseSpeed * float32(deltaTime)
	if !(glfwWindow.GetKey(glfw.KeyW) == glfw.Release) {
		c.moveForward(speed)
	}
	if !(glfwWindow.GetKey(glfw.KeyS) == glfw.Release) {
		c.moveBackward(speed)
	}
	if !(glfwWindow.GetKey(glfw.KeyA) == glfw.Release) {
		c.moveLeft(speed)
	}
	if !(glfwWindow.GetKey(glfw.KeyD) == glfw.Release) {
		c.moveRight(speed)
	}
	c.camera.Resize(glfwWindow.GetSize())
	c.recalculateViewMatrix()
}

// GetViewProjectionMatrix .
func (c *CameraController) GetViewProjectionMatrix() mgl32.Mat4 {
	return c.camera.GetProjectionMatrix().Mul4(c.viewMatrix)
}

func (c *CameraController) moveForward(speed float32) {
	c.pos = c.pos.Add(c.target.Mul(speed))
}

func (c *CameraController) moveBackward(speed float32) {
	c.pos = c.pos.Sub(c.target.Mul(speed))
}

func (c *CameraController) moveLeft(speed float32) {
	c.pos = c.pos.Sub(c.target.Normalize().Cross(c.up).Mul(speed))
}

func (c *CameraController) moveRight(speed float32) {
	c.pos = c.pos.Add(c.target.Normalize().Cross(c.up).Mul(speed))
}

func (c *CameraController) recalculateViewMatrix() {
	c.viewMatrix = mgl32.LookAtV(c.pos, c.pos.Add(c.target), c.up)
}

func (c *CameraController) getViewMatrix() mgl32.Mat4 {
	return c.viewMatrix
}
