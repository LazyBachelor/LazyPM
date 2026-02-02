tidy:
	go mod tidy

clean:
	go clean

build:
	go build -o ./bin/cli ./cmd/cli
	go build -o ./bin/tui ./cmd/tui
	go build -o ./bin/web ./cmd/web

cli:
	go run ./cmd/cli

tui: 
	go run ./cmd/tui

web:
	go run ./cmd/web

dev:
	go tool templ generate -watch -cmd "go run ./cmd/web"