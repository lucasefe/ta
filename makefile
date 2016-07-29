all: bin/ta

bin/ta: *.go bin
	go build -o $@

bin pkg:
	mkdir -p $@

cross-compile: bin/ta
	script/cross-compile ta

clean:
	rm -rf pkg

.PHONY: cross-compile clean
