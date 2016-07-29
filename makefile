all: build

build: *.go bin
	go build -o bin/ta

bin: 
	mkdir -p $@
