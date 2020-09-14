package window

// Option .
type Option func(*Window) error

// WithDimensionsOption .
func WithDimensionsOption(width, height int) Option {
	return func(w *Window) error {
		w.width = width
		w.height = height
		return nil
	}
}

// WithTitleOption .
func WithTitleOption(title string) Option {
	return func(w *Window) error {
		w.title = title
		return nil
	}
}
