SHELL := /bin/bash

tidy:
	go mod tidy

clean:
	@go clean
	@rm -rf ./bin
	@rm -rf ./.pm
	@rm -rf node_modules
	@rm package-lock.json
	@rm package.json

build: tidy
	go build -o ./bin/pm ./cmd/pm
	go build -o ./bin/tui ./cmd/tui
	go build -o ./bin/web ./cmd/web
	go build -o ./bin/survey ./cmd/survey

enable-multiplatform-build:
	@docker buildx create --name multiplatform --use 2>/dev/null || docker buildx use multiplatform

docker-build: enable-multiplatform-build
	@docker buildx build \
  --platform linux/amd64,linux/arm64,linux/arm/v7 \
  -t telikz/lazypm:latest \
  --push .


docker-run:
	@docker run -it -p 8080:8080 -v "$PWD:/data" telikz/lazypm:latest start

docker-push:
	@docker push telikz/lazypm:latest

build-os: build
	@cp ./bin/survey ./build/survey
	@docker build -t lazyos ./build

start:
	@go run ./cmd/survey start

cli: tidy
	go run ./cmd/pm

tui: tidy
	go run ./cmd/tui

web: tidy
	go run ./cmd/web

tw-install:
	@npm install tailwindcss@latest @tailwindcss/cli@latest @tailwindcss/typography daisyui@latest

# Run both dev and tw in parallel for generating templates and compiling Tailwind CSS on file changes
dev: tidy
	@go tool templ generate -watch -cmd "go run ./cmd/web"

tw: tw-install
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