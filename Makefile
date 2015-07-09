.PHONY: default release

default:
	mkdir -p bin
	go build -o bin/vessel ./src

release:
	mkdir -p release
	CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o release/vessel ./src

test:
	go test -v ./src
