package main

import (
	"runtime"

	"github.com/devodev/opengl-experimentation/components"
)

func init() {
	// All calls to GLFW/OpenGL must happen on the main thread.
	// This locks the calling goroutine(main here) to
	// the current OS Thread(main here).
	runtime.LockOSThread()
}

func main() {
	logger := NewLogger()
	app, err := NewApplication(WithLoggerOption(logger))
	if err != nil {
		logger.Errorf("error creating application: %s", err)
		return
	}
	defer func() {
		if err := app.Close(); err != nil {
			logger.Errorf("error closing application: %s", err)
		}
	}()

	// square, err := components.NewSquare()
	squareTexture, err := components.NewSquareTexture()
	if err != nil {
		logger.Errorf("error creating application: %s", err)
		return
	}
	app.AddComponent(squareTexture)

	if err := app.Run(); err != nil {
		logger.Errorf("error running application: %s", err)
		return
	}
}
