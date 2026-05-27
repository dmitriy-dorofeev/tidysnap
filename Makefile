.PHONY: build install uninstall run clean test check install-tools

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

check:
	@echo "==> Checking formatting..."
	@test -z "$$(gofmt -l .)" || (echo "Найдены неотформатированные файлы:" && gofmt -l . && exit 1)
	@echo "==> Running go vet..."
	go vet ./...
	@echo "==> Running staticcheck..."
	staticcheck ./...
	@echo "==> Running govulncheck..."
	govulncheck ./...
	@echo "==> Running gosec..."
	gosec -exclude=G404 ./...
	@echo "==> All checks passed!"

install-tools:
	go install honnef.co/go/tools/cmd/staticcheck@2026.1
	go install golang.org/x/vuln/cmd/govulncheck@v1.3.0
	go install github.com/securego/gosec/v2/cmd/gosec@v2.22.11
