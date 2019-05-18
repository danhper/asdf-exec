all: linux macos


build:
	mkdir -p build

linux: build
	GOOS=linux GOARCH=amd64 go build -o build/asdf-exec-linux-x64

macos: build
	GOOS=darwin GOARCH=amd64 go build -o build/asdf-exec-darwin-x64
