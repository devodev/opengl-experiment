package application

import (
	"github.com/devodev/opengl-experimentation/internal/engine"
	"github.com/devodev/opengl-experimentation/internal/engine/window"
)

// Option .
type Option func(*Application) error

// WithLoggerOption .
func WithLoggerOption(logger *engine.SimpleLogger) Option {
	return func(a *Application) error {
		a.logger = logger
		return nil
	}
}

// WithWindowOption .
func WithWindowOption(window *window.Window) Option {
	return func(a *Application) error {
		a.window = window
		return nil
	}
}
