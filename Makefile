.DEFAULT_GOAL := run 

.PHONY: build
vet: 
	go vet ./cmd/web/
build: vet
	go build -o=/tmp/bin/main cmd/web/main.go
run: build
	/tmp/bin/main
