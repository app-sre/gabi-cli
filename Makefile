OUTPUT_DIR :=_output

# Constants
GOPATH := $(shell go env GOPATH)

clean:
	rm gabi

build: clean
	go build -o gabi ./cmd/gabi || exit 1

install:
	go build -o $(GOPATH)/bin/gabi ./cmd/gabi || exit 1
	gabi version

uninstall:
	rm -f $(GOPATH)/bin/gabi



