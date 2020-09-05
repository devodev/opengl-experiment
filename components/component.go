package components

import "github.com/go-gl/glfw/v3.3/glfw"

// Component .
type Component interface {
	OnInit(*glfw.Window)
	OnUpdate(*glfw.Window)
}
