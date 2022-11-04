install:
	go install ./cmd/gabi

uninstall:
	rm -f $(GOPATH)/bin/gabi

clean:
	rm -f gabi
