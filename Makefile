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
	@go tool templ generate -watch -cmd "go run ./cmd/web"

tw:
	@npx --yes @tailwindcss/cli -i ./pkg/web/input.css -o ./pkg/web/assets/css/styles.css --watch

watch:
	@make -j2 dev tw

.PHONY: tidy clean build cli tui web dev tw