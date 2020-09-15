package main

import (
	"runtime"

	"github.com/devodev/opengl-experimentation/internal/engine"
	"github.com/devodev/opengl-experimentation/internal/engine/application"
)

func init() {
	// All calls to GLFW/OpenGL must happen on the main thread.
	// This locks the calling goroutine(main here) to
	// the current OS Thread(main here).
	runtime.LockOSThread()
}

func main() {
	logger := engine.NewLogger()
	app, err := application.New(application.WithLoggerOption(logger))
	if err != nil {
		logger.Errorf("error creating application: %s", err)
		return
	}
	defer func() {
		if err := app.Close(); err != nil {
			if err != application.ErrAlreadyClosed {
				logger.Errorf("error closing application: %s", err)
			}
		}
	}()

	squareTextureLayer, err := NewSquareTexture(app)
	if err != nil {
		logger.Errorf("error creating layer: %s", err)
		return
	}
	app.AddLayer(squareTextureLayer)

	if err := app.Run(); err != nil {
		logger.Errorf("error running application: %s", err)
		return
	}
}
