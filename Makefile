.PHONY: all build-cliente build-servidor clean distclean

BINARY_CLIENTE=cliente
BINARY_SERVIDOR=servidor

all: build-cliente build-servidor

go.mod:
	test -f go.mod || go mod init fppd-jogo-t2
	go get -u github.com/nsf/termbox-go

build-cliente: go.mod
	go build -o $(BINARY_CLIENTE) ./src/cliente

build-servidor: go.mod
	go build -o $(BINARY_SERVIDOR) ./src/servidor

clean:
	rm -f $(BINARY_CLIENTE) $(BINARY_SERVIDOR)

distclean: clean
	rm -f go.mod go.sum
