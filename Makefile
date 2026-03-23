.PHONY: build test lint clean install uninstall demo

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
	cp ppp.1 /usr/local/share/man/man1/ppp.1

uninstall:
	rm -f /usr/local/bin/ppp
	rm -f /usr/local/share/man/man1/ppp.1

demo: build
	@for tape in demo/*.tape; do echo "Recording $$tape..."; vhs $$tape; done
