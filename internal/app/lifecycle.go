package app

type Lifecycle struct {
	cleanups []func()
}

func NewLifecycle() *Lifecycle {
	return &Lifecycle{
		cleanups: make([]func(), 0),
	}
}

func (l *Lifecycle) Add(fn func()) {
	l.cleanups = append(l.cleanups, fn)
}

func (l *Lifecycle) Close() {
	for i := len(l.cleanups) - 1; i >= 0; i-- {
		l.cleanups[i]()
	}
}
