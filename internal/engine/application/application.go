package application

import (
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/devodev/opengl-experiment/internal/engine"
	"github.com/devodev/opengl-experiment/internal/engine/renderer"
	"github.com/devodev/opengl-experiment/internal/engine/window"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Errors
var (
	ErrAlreadyClosed = errors.New("application already closed")
)

var (
	app *application
)

func init() {
	if err := initApp(); err != nil {
		panic(err)
	}
}

func initApp() error {
	window, err := window.New()
	if err != nil {
		return err
	}
	renderer, err := renderer.New()
	if err != nil {
		return err
	}
	app = &application{
		window:       window,
		renderer:     renderer,
		logger:       engine.NewLogger(),
		frameCounter: NewFrameCounter(),
	}
	return nil
}

type application struct {
	running bool

	profiling bool

	window       *window.Window
	renderer     *renderer.Renderer
	logger       *engine.SimpleLogger
	frameCounter *FrameCounter

	layers []Layer
}

func (a *application) init() error {
	if err := app.window.Init(); err != nil {
		return fmt.Errorf("error initializing window: %v", err)
	}
	app.logger.Printf("GLFW version: %s", glfw.GetVersionString())

	if err := app.renderer.Init(); err != nil {
		return fmt.Errorf("error initializing renderer: %v", err)
	}
	app.logger.Printf("OpenGL version: %s", gl.GoStr(gl.GetString(gl.VERSION)))
	return nil
}

func (a *application) run() error {
	a.setupProfiling()

	a.running = true

	// init components
	for _, layer := range a.layers {
		if err := layer.OnInit(); err != nil {
			return err
		}
	}

	a.frameCounter.Init(glfw.GetTime())

	// main loop
	for a.running {
		a.frameCounter.OnUpdate(glfw.GetTime())
		deltaTime := a.frameCounter.Delta()

		a.onUpdate()
		a.renderer.Clear()
		for _, layer := range a.layers {
			layer.OnUpdate(deltaTime)
			layer.OnRender(deltaTime)
		}
		a.window.GetGLFWWindow().SwapBuffers()
	}
	return nil
}

func (a *application) setupProfiling() {
	if !a.profiling {
		return
	}

	// starts a http server that will serve pprof endpoints
	// pprof endpoints are registered when importing its package (_ "net/http/pprof")
	go func() {
		addr := "localhost:6060"
		a.logger.Infof("pprof server listening on http://%s", addr)
		a.logger.Println(http.ListenAndServe(addr, nil))
	}()
}

func (a *application) onUpdate() {
	if a.window.ShouldClose() {
		a.running = false
		return
	}
	// mainthread only
	glfw.PollEvents()
	//
}

// Run .
func Run() error {
	if err := app.init(); err != nil {
		return err
	}
	defer app.window.Close()

	return app.run()
}

// Close .
func Close() {
	app.running = false
}

// AddLayer .
func AddLayer(l Layer) {
	app.layers = append(app.layers, l)
}

// GetWindow .
func GetWindow() *window.Window {
	return app.window
}

// GetRenderer .
func GetRenderer() *renderer.Renderer {
	return app.renderer
}

// SetLogger .
func SetLogger(logger *engine.SimpleLogger) {
	app.logger = logger
}

// SetWindow .
func SetWindow(window *window.Window) {
	app.window = window
}

// SetWindow .
func SetWindowSize(width, height int) {
	app.window.SetSize(width, height)
}

func EnableProfiling() {
	app.profiling = true
}
