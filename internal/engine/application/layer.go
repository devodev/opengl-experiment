package application

// Layer .
type Layer interface {
	OnInit() error
	OnUpdate(float64)
	OnRender(float64)
}
