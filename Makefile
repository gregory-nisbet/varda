SRCDIRS := cmd pkg

build:
	go fmt ./...
	go test ./...
	go build -o=./varda ./cmd/varda/...

gccgo:
	go test -compiler=gccgo ./...
	go build -compiler=gccgo -o=./varda.gccgo ./cmd/varda/...

install:
	install -T ./varda /usr/local/bin/varda

cscope:
	find $(SRCDIRS) -name '*.go' > cscope.files
	cscope -b -q -k

.PHONY: tags
tags:
	ctags-universal -R .
