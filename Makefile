all: 
	go fmt *.go
	go build
	go test

.PHONY: all