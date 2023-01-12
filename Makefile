build:
	go build -o gabi ./cmd/gabi || exit 1

install:
	go install ./cmd/gabi

uninstall:
	rm -f $(GOPATH)/bin/gabi

clean:
	rm -f gabi
