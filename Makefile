SHELL := /bin/bash

ifneq (,$(wildcard .env))
	include .env
	export
endif

LDFLAGS := \
	-X 'main.DB_URI=$(DB_URI)'

tidy:
	go mod tidy

clean:
	@go clean
	@rm -rf ./.pm
	@rm -rf ./bin
	@rm -rf ./dist
	@rm -rf ./.idea
	@rm -rf node_modules
	@rm -f package-lock.json
	@rm -f package.json
	@rm -f ./build/pm
	@rm -f ./build/pm_bash.sh

build: tidy
	@go build -ldflags "$(LDFLAGS)" -o ./bin/pm ./cmd/pm
	@go build -o ./bin/tui ./cmd/tui
	@go build -o ./bin/web ./cmd/web
	@echo "Build completed successfully. Binaries are located in the ./bin directory."

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

os-build: completions build
	@cp ./bin/pm ./build/pm
	@cp ./bin/pm_bash.sh ./build/pm_bash.sh
	@docker build -t telikz/lazyos ./build
	@echo "Built lazyos image successfully. You can run it using 'make os-run'."

os-run: os-build
	@docker run -d   --name=lazyos   -e PUID=1000   -e PGID=1000   -e TZ=Etc/UTC   -p 3000:3000   -p 3001:3001   --shm-size="1gb"   telikz/lazyos
	@echo "Started lazyos container successfully. You can access the web interface at http://localhost:3000."

os-stop:
	@docker stop lazyos
	@docker rm -v --force lazyos
	@echo "Stopped and removed lazyos container successfully."

os-push:
	@docker push telikz/lazyos

start:
	@go run ./cmd/pm survey start

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
	@mkdir -p ./bin
	@go build -o ./bin/pm ./cmd/pm
	@DEV=True ./bin/pm completion bash > ./bin/pm_bash.sh
	@DEV=True ./bin/pm completion zsh > ./bin/pm_zsh.sh
	@DEV=True ./bin/pm completion fish > ./bin/pm_fish.sh
	@DEV=True ./bin/pm completion powershell > ./bin/pm_powershell.ps1

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


.PHONY: tidy clean build docker-build docker-run docker-push os-build os-run os-stop os-push start cli tui web tw-install dev tw completions install-bash-temp install-zsh-temp install-fish-temp install-powershell-temp install-cli
