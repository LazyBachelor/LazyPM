tidy:
	go mod tidy

clean:
	go clean
	
cli:
	go run ./cmd/cli

tui: 
	go run ./cmd/tui

web:
	go run ./cmd/web

dev:
	templ generate -watch -cmd "go run ./cmd/web"