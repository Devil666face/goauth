GOOS=linux
GOARCH=amd64
LDFLAGS="-w -s"
APP=app

build: ## Build project
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags=$(LDFLAGS) -o $(APP) main.go

run: ## Temp build and run
	go build -o tmp/main main.go
	tmp/main

clear: ## Clear temp dirs
	rm -rf db.sqlite3
	rm -rf static
	rm -rf tmp
	go build -o tmp/main main.go
	tmp/main -migrate
	tmp/main -superuser

air: ## Run dev server
	~/go/bin/air

air-install: ## Install air
	go install github.com/cosmtrek/air@latest

help: ## Prints help for targets with comments
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

gencert: ## Generate server.key and server.crt
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout server.key -out server.crt