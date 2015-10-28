.PHONY: default release deps release release-dc

default:
	mkdir -p bin
	go build -o bin/vessel ./client

release:
	mkdir -p release
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o 	release/vessel-v1.0.0-linux-amd64 ./client

release-dc:
	docker-compose run app make release

deps:
	godep save -r ./...

test:
	go test -v ./client/...
