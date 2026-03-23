.PHONY: build test lint clean install

build:
	go build -o pp .

test:
	go test ./... -count=1 -race

lint:
	go vet ./...
	gofmt -l .
	@test -z "$$(gofmt -l .)" || (echo "Run gofmt" && exit 1)

clean:
	rm -f pp

install: build
	cp pp /usr/local/bin/pp
