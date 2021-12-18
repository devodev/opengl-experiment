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
	delta := currentTime - f.fpsLastTime
	if delta >= f.printDeltaSeconds {
		fps := float64(f.nbFrames) / delta
		frameTime := (delta * 1000) / float64(f.nbFrames)
		fmt.Printf("%.2f fps (%.2f ms/frame)\n", fps, frameTime)

		f.nbFrames = 0
		f.fpsLastTime += delta
	}
}

// GetDelta .
func (f *FrameCounter) Delta() float64 {
	return f.deltaTime
}
