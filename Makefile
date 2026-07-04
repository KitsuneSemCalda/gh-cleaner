BINARY := gh-cleaner
CMD_DIR := ./cmd/gh-cleaner
INSTALL_DIR := $(HOME)/.local/bin

.PHONY: all build run test clean vet fmt install uninstall

all: build

build:
	go build -o $(BINARY) $(CMD_DIR)

run: build
	./$(BINARY)

test:
	go test ./...

clean:
	go clean
	rm -f $(BINARY)

install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY) $(INSTALL_DIR)/$(BINARY)

uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY)

vet:
	go vet ./...

fmt:
	go fmt ./...
