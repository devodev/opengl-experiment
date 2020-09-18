package renderer

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// Camera .
type Camera interface {
	OnUpdate(*glfw.Window, float64)
	GetViewProjectionMatrix() mgl32.Mat4
}

// CameraPerspective .
type CameraPerspective struct {
	pos   mgl32.Vec3
	front mgl32.Vec3
	up    mgl32.Vec3
	speed float32

	projectionMatrix     mgl32.Mat4
	viewMatrix           mgl32.Mat4
	viewProjectionMatrix mgl32.Mat4
}

// NewCameraPerspective .
func NewCameraPerspective(width, height int) *CameraPerspective {
	camera := &CameraPerspective{
		pos:   mgl32.Vec3{0, 0, 2},
		front: mgl32.Vec3{0, 0, -0.5},
		up:    mgl32.Vec3{0, 1, 0},
		speed: float32(0.05),
	}
	camera.recalculateViewProjectionMatrix(width, height)
	return camera
}

// OnUpdate .
func (c *CameraPerspective) OnUpdate(glfwWindow *glfw.Window, deltaTime float64) {
	c.speed = float32(2 * deltaTime)
	if glfwWindow.GetAttrib(glfw.Focused) != glfw.False {
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
	}
	c.recalculateViewProjectionMatrix(glfwWindow.GetSize())
}

func (c *CameraPerspective) recalculateViewProjectionMatrix(width, height int) {
	c.projectionMatrix = mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/float32(height), 0.1, 10.0)
	c.viewMatrix = mgl32.LookAtV(c.pos, c.pos.Add(c.front), c.up)
	c.viewProjectionMatrix = c.projectionMatrix.Mul4(c.viewMatrix)
}

// GetViewProjectionMatrix .
func (c *CameraPerspective) GetViewProjectionMatrix() mgl32.Mat4 {
	return c.viewProjectionMatrix
}
