package renderer

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	defaultControllerPos                 = mgl32.Vec3{0, 0, 2}
	defaultControllerTarget              = mgl32.Vec3{0, 0, -0.5}
	defaultControllerUp                  = mgl32.Vec3{0, 1, 0}
	defaultControllerBaseSpeed           = float32(2)
	defaultControllerRotationSensitivity = float32(2)
	defaultControllerYaw                 = float32(-90.0)
	defaultControllerPitch               = float32(0.0)
)

func sin(v float32) float32 {
	return float32(math.Sin(float64(v)))

}
func cos(v float32) float32 {
	return float32(math.Cos(float64(v)))
}

// CameraController .
type CameraController struct {
	pos                   mgl32.Vec3
	target                mgl32.Vec3
	up                    mgl32.Vec3
	baseSpeed             float32
	rotationSensitivity   float32
	yaw                   float32
	pitch                 float32
	mousePosX             float64
	mousePosY             float64
	mouseButton1IsPressed bool

	viewMatrix mgl32.Mat4

	camera Camera
}

// NewCameraController .
func NewCameraController(camera Camera) *CameraController {
	width, height := camera.GetViewPortDimensions()
	return &CameraController{
		pos:                   defaultControllerPos,
		target:                defaultControllerTarget,
		up:                    defaultControllerUp,
		baseSpeed:             defaultControllerBaseSpeed,
		rotationSensitivity:   defaultControllerRotationSensitivity,
		yaw:                   defaultControllerYaw,
		pitch:                 defaultControllerPitch,
		mousePosX:             float64(width / 2),
		mousePosY:             float64(height / 2),
		mouseButton1IsPressed: false,
		camera:                camera,
	}
}

// OnUpdate .
func (c *CameraController) OnUpdate(glfwWindow *glfw.Window, deltaTime float64) {
	if glfwWindow.GetAttrib(glfw.Focused) == glfw.False {
		return
	}

	// speed
	speed := c.baseSpeed * float32(deltaTime)
	// position
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
	// rotation
	windowWidth, windowHeight := glfwWindow.GetSize()
	cursorX, cursorY := glfwWindow.GetCursorPos()
	if cursorX >= 0 &&
		cursorY >= 0 &&
		cursorX <= float64(windowWidth) &&
		cursorY <= float64(windowHeight) {
		if glfwWindow.GetMouseButton(glfw.MouseButton1) == glfw.Press {
			if !c.mouseButton1IsPressed {
				c.mouseButton1IsPressed = true
				c.mousePosX = cursorX
				c.mousePosY = cursorY
			}
			c.rotate(speed, cursorX, cursorY)
		} else {
			c.mouseButton1IsPressed = false
		}
	}

	c.camera.Resize(windowWidth, windowHeight)
	c.recalculateViewMatrix()
}

// GetViewProjectionMatrix .
func (c *CameraController) GetViewProjectionMatrix() mgl32.Mat4 {
	return c.camera.GetProjectionMatrix().Mul4(c.viewMatrix)
}

func (c *CameraController) rotate(speed float32, posX, posY float64) {
	xOffset := posX - c.mousePosX
	yOffset := c.mousePosY - posY
	c.mousePosX = posX
	c.mousePosY = posY

	c.yaw -= float32(xOffset) * c.rotationSensitivity * speed
	c.pitch -= float32(yOffset) * c.rotationSensitivity * speed
	c.recalculateTarget()
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

func (c *CameraController) recalculateTarget() {
	c.target = mgl32.Vec3{
		cos(mgl32.DegToRad(c.yaw)) * cos(mgl32.DegToRad(c.pitch)),
		sin(mgl32.DegToRad(c.pitch)),
		sin(mgl32.DegToRad(c.yaw)) * cos(mgl32.DegToRad(c.pitch)),
	}
}
