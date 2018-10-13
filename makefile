.PHONY: help run install

help:
	# http://repo.msys2.org/mingw/x86_64/

run:
	mkdir -p bin
	go build -i -o bin/test cmd/test.go && cd bin && ./test
