.DEFAULT_GOAL := run 

.PHONY: build
tidy:
	go mod tidy
vet: tidy
	go vet ./...
build: vet
	go build -o=/tmp/bin/main cmd/web/main.go
run: build
	/tmp/bin/main
