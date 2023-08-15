run:
	go build -o tmp/main app/cmd/main.go
	tmp/main
clear:
	rm -rf static
	rm -rf tmp
