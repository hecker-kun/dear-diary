run:
	go run api/main.go
build:
	go build -v -o bin/ api/main.go
lint:
	golangci-lint run

git:
	git add .
	git rm -r --cached .idea bin config/config_dev.go
	git commit -m "$a"
	git push -u origin main