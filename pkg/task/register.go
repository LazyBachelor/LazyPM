package task

import (
	"fmt"

	"github.com/LazyBachelor/LazyPM/internal/service"
)

var registry = make(map[string]func(*service.Services) Tasker)

func Register(name string, constructor func(*service.Services) Tasker) {
	if _, exists := registry[name]; exists {
		panic(fmt.Sprintf("task %q already registered", name))
	}
	registry[name] = constructor
}

func Get(name string, svc *service.Services) (Tasker, error) {
	constructor, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("task %q not found", name)
	}
	return constructor(svc), nil
}

func List() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}
