package task

import (
	"fmt"
)

var interfaceRegistry = make(map[string]Interface)

func RegisterInterface(name string, iface Interface) {
	if _, exists := interfaceRegistry[name]; exists {
		panic(fmt.Sprintf("interface %q already registered", name))
	}
	interfaceRegistry[name] = iface
}

func GetInterface(name string) (Interface, error) {
	iface, ok := interfaceRegistry[name]
	if !ok {
		return nil, fmt.Errorf("interface %q not found", name)
	}
	return iface, nil
}

func ListInterfaces() []string {
	names := make([]string, 0, len(interfaceRegistry))
	for name := range interfaceRegistry {
		names = append(names, name)
	}
	return names
}

var taskRegistry = make(map[string]func(*App) Tasker)

func RegisterTask(name string, constructor func(*App) Tasker) {
	if _, exists := taskRegistry[name]; exists {
		panic(fmt.Sprintf("task %q already registered", name))
	}
	taskRegistry[name] = constructor
}

func GetTask(name string, app *App) (Tasker, error) {
	constructor, ok := taskRegistry[name]
	if !ok {
		return nil, fmt.Errorf("task %q not found", name)
	}
	return constructor(app), nil
}

func ListTasks() []string {
	names := make([]string, 0, len(taskRegistry))
	for name := range taskRegistry {
		names = append(names, name)
	}
	return names
}
