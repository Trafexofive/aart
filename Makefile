.PHONY: build run clean install

build:
	go build -o aart ./cmd/aart

run: build
	./aart

clean:
	rm -f aart

install: 
	go install ./cmd/aart
	cp aart /usr/local/bin/aart

uninstall:
	rm -f /usr/local/bin/aart

reinstall: uninstall install

test:
	go test ./...

fmt:
	go fmt ./...

deps:
	go mod tidy
	go mod download
