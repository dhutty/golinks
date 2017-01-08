.PHONY: dev build clean

all: dev

dev: build
	./golinks -bind 127.0.0.1:8000

build: clean
	go get ./...
	go build -o ./golinks .

clean:
	rm -rf golinks
