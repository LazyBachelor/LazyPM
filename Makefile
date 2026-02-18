SHELL := /bin/bash

tidy:
	go mod tidy

clean:
	go clean

build:
	go build -o ./bin/pm ./cmd/pm
	go build -o ./bin/tui ./cmd/tui
	go build -o ./bin/web ./cmd/web

docker-build:
	@docker build -t survey .

docker-run:
	@docker run -it -p 8080:8080 survey:latest

cli:
	go run ./cmd/pm

tui: 
	go run ./cmd/tui

web:
	go run ./cmd/web

# Run both dev and tw in parallel for generating templates and compiling Tailwind CSS on file changes
dev:
	@go tool templ generate -watch -cmd "go run ./cmd/web"
tw:
	@npm install tailwindcss@latest @tailwindcss/cli@latest @tailwindcss/typography daisyui@latest
	@npx --yes @tailwindcss/cli -i ./pkg/web/input.css -o ./pkg/web/assets/css/styles.css --watch --minify


completions:
	@go build -o ./bin/pm ./cmd/pm
	@mkdir -p ./bin
	@./bin/pm completion bash > ./bin/pm_bash.sh
	@./bin/pm completion zsh > ./bin/pm_zsh.sh
	@./bin/pm completion fish > ./bin/pm_fish.sh
	@./bin/pm completion powershell > ./bin/pm_powershell.ps1

install-bash-temp: completions
	@go install ./cmd/pm
	@source ./bin/pm_bash.sh

install-zsh-temp: completions
	@go install ./cmd/pm
	@source ./bin/pm_zsh.sh

install-fish-temp: completions
	@go install ./cmd/pm
	@source ./bin/pm_fish.sh

install-powershell-temp: completions
	@go install ./cmd/pm
	@source ./bin/pm_powershell.ps1


install-cli: completions
	@go install ./cmd/pm
	@sudo cp ./bin/pm_bash.sh /etc/bash_completion.d/pm


.PHONY: tidy clean build cli tui web dev tw completions install-bash-temp install-zsh-temp install-fish-temp install-powershell-temp install-cli