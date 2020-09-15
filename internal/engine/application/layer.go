package application

// Layer .
type Layer interface {
	OnInit(*Application)
	OnUpdate(*Application, float64)
	OnRender(*Application, float64)
}
