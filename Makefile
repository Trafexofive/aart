.PHONY: build run clean install

build:
	go build -o aart ./cmd/aart

run: build
	./aart

clean:
	rm -f aart

install:
	go install ./cmd/aart

test:
	go test ./...

fmt:
	go fmt ./...

deps:
	go mod tidy
	go mod download
