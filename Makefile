.PHONY: default release deps release release-dc

default:
	mkdir -p bin
	go build -o bin/vessel ./client

release:
	mkdir -p release
	CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o release/vessel ./client

release-dc:
	docker-compose run app make release

deps:
	godep save -r ./...

test:
	go test -v ./client/...
