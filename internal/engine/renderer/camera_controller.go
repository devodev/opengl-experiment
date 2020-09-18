package renderer

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	defaultControllerPos   = mgl32.Vec3{0, 0, 2}
	defaultControllerFront = mgl32.Vec3{0, 0, -0.5}
	defaultControllerUp    = mgl32.Vec3{0, 1, 0}
	defaultControllerSpeed = float32(0.05)
)

// CameraController .
type CameraController struct {
	pos   mgl32.Vec3
	front mgl32.Vec3
	up    mgl32.Vec3
	speed float32

	viewMatrix mgl32.Mat4

	camera Camera
}

// NewCameraController .
func NewCameraController(camera Camera) *CameraController {
	return &CameraController{
		pos:    defaultControllerPos,
		front:  defaultControllerFront,
		up:     defaultControllerUp,
		speed:  defaultControllerSpeed,
		camera: camera,
	}
}

// OnUpdate .
func (c *CameraController) OnUpdate(glfwWindow *glfw.Window, deltaTime float64) {
	if glfwWindow.GetAttrib(glfw.Focused) == glfw.False {
		return
	}

	c.speed = float32(2 * deltaTime)
	if !(glfwWindow.GetKey(glfw.KeyW) == glfw.Release) {
		c.pos = c.pos.Add(c.front.Mul(c.speed))
	}
	if !(glfwWindow.GetKey(glfw.KeyS) == glfw.Release) {
		c.pos = c.pos.Sub(c.front.Mul(c.speed))
	}
	if !(glfwWindow.GetKey(glfw.KeyA) == glfw.Release) {
		c.pos = c.pos.Sub(c.front.Normalize().Cross(c.up).Mul(c.speed))
	}
	if !(glfwWindow.GetKey(glfw.KeyD) == glfw.Release) {
		c.pos = c.pos.Add(c.front.Normalize().Cross(c.up).Mul(c.speed))
	}
	c.camera.Resize(glfwWindow.GetSize())
	c.recalculateViewMatrix()
}

// GetViewProjectionMatrix .
func (c *CameraController) GetViewProjectionMatrix() mgl32.Mat4 {
	return c.camera.GetProjectionMatrix().Mul4(c.viewMatrix)
}

func (c *CameraController) recalculateViewMatrix() {
	c.viewMatrix = mgl32.LookAtV(c.pos, c.pos.Add(c.front), c.up)
}

func (c *CameraController) getViewMatrix() mgl32.Mat4 {
	return c.viewMatrix
}
