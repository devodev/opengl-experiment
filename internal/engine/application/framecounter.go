package application

import (
	"fmt"
)

var (
	defaultPrintDeltaSeconds = 1.0
)

// FrameCounter .
type FrameCounter struct {
	printDeltaSeconds float64

	deltaTime   float64
	lastTime    float64
	fpsLastTime float64
	nbFrames    int
}

// NewFrameCounter .
func NewFrameCounter() *FrameCounter {
	return &FrameCounter{
		printDeltaSeconds: defaultPrintDeltaSeconds,
	}
}

// Init .
func (f *FrameCounter) Init(currentTime float64) {
	f.lastTime = currentTime
	f.fpsLastTime = currentTime
}

// OnUpdate .
func (f *FrameCounter) OnUpdate(currentTime float64) {
	f.deltaTime = currentTime - f.lastTime
	f.lastTime = currentTime

	f.nbFrames++
	if currentTime-f.fpsLastTime >= f.printDeltaSeconds {
		fmt.Printf("%f ms/frame\n", (f.printDeltaSeconds*1000)/float64(f.nbFrames))
		f.nbFrames = 0
		f.fpsLastTime += 1.0
	}
}

// GetDelta .
func (f *FrameCounter) GetDelta() float64 {
	return f.deltaTime
}
