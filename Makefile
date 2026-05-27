.PHONY: build install uninstall run clean test

BINARY_NAME=tidysnap
INSTALL_PATH=/usr/local/bin

build:
	go build -ldflags "-s -w" -o bin/$(BINARY_NAME) ./cmd/main.go

run:
	go run ./cmd/main.go

install: build
	cp bin/$(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Installed to $(INSTALL_PATH)/$(BINARY_NAME)"
	@echo "Run '$(BINARY_NAME)' to start"

uninstall:
	$(INSTALL_PATH)/$(BINARY_NAME) --uninstall 2>/dev/null || true
	rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Uninstalled"

clean:
	rm -rf bin/

test:
	go test ./...
