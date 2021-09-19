run:
	go run api/main.go
build:
	go build -v -o bin/ api/main.go
lint:
	golangci-lint run