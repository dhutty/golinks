.PHONY: dev build clean

all: dev

dev: build
	./search -bind 0.0.0.0:8000

build: clean
	go build -o ./search .

clean:
	rm -rf search
