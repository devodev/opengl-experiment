package application

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// Camera .
type Camera struct {
	pos   mgl32.Vec3
	front mgl32.Vec3
	up    mgl32.Vec3
	speed float32

	projectionMatrix     mgl32.Mat4
	viewMatrix           mgl32.Mat4
	viewProjectionMatrix mgl32.Mat4
}

// NewCamera .
func NewCamera(width, height int) *Camera {
	posVec3 := mgl32.Vec3{0, 0, 2}
	frontVec3 := mgl32.Vec3{0, 0, -1}
	upVec3 := mgl32.Vec3{0, 1, 0}
	speed := float32(0.05)

	projM := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/float32(height), 0.1, 10.0)
	viewM := mgl32.LookAtV(posVec3, posVec3.Add(frontVec3), upVec3)
	viewProjM := projM.Mul4(viewM)

	return &Camera{
		pos:   posVec3,
		front: frontVec3,
		up:    upVec3,
		speed: speed,

		projectionMatrix:     projM,
		viewMatrix:           viewM,
		viewProjectionMatrix: viewProjM,
	}
}

// OnUpdate .
func (c *Camera) OnUpdate(app *Application, deltaTime float64) {
	c.speed = float32(2 * deltaTime)

	glfwWindow := app.GetWindow().GetGLFWWindow()
	if glfwWindow.GetAttrib(glfw.Focused) == glfw.False {
		return
	}

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

	width, height := app.GetWindow().GetGLFWWindow().GetSize()
	c.projectionMatrix = mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/float32(height), 0.1, 10.0)
	c.viewMatrix = mgl32.LookAtV(c.pos, c.pos.Add(c.front), c.up)
	c.viewProjectionMatrix = c.projectionMatrix.Mul4(c.viewMatrix)
}

// GetViewProjectionMatrix .
func (c *Camera) GetViewProjectionMatrix() mgl32.Mat4 {
	return c.viewProjectionMatrix
}
