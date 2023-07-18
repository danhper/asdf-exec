all: linux macos


build:
	mkdir -p build

linux: build
	GOOS=linux go build -o build/asdf-exec-linux

macos: build
	GOOS=darwin go build -o build/asdf-exec-darwin
