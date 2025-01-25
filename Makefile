.DEFAULT_GOAL := run 

.PHONY: build
vet: 
	go vet ./...
build: vet
	go build -o=/tmp/bin/main main.go
run: build
	/tmp/bin/main
