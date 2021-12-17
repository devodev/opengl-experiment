package main

import (
	"runtime"

	"github.com/devodev/opengl-experiment/internal/engine"
	"github.com/devodev/opengl-experiment/internal/engine/application"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	logger := engine.NewLogger()
	application.SetLogger(logger)

	layer := &SquareTextureLayer{}
	application.AddLayer(layer)

	if err := application.Run(); err != nil {
		logger.Errorf("error running application: %s", err)
		return
	}
}
