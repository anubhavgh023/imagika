.DEFAULT_GOAL := run 

.PHONY: build
build:
	go build -o=/tmp/bin/main main.go
run: build
	/tmp/bin/main
