.PHONY: help run install

help:

run:
	mkdir -p bin
	go build -i -o bin/test cmd/test.go && cd bin && ./test
