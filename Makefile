run:
	go build -o tmp/main main.go
	tmp/main
clear:
	rm -rf static
	rm -rf tmp
air:
	~/go/bin/air
air-install:
	go install github.com/cosmtrek/air@latest
