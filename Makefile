.PHONY: build test lint clean install

build:
	go build -o ppp .

test:
	go test ./... -count=1 -race

lint:
	go vet ./...
	gofmt -l .
	@test -z "$$(gofmt -l .)" || (echo "Run gofmt" && exit 1)

clean:
	rm -f ppp

install: build
	cp ppp /usr/local/bin/ppp
