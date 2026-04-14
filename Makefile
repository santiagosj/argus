.PHONY: help build demo run clean test install docker docker-run

help:
	@echo "Argus - Cognitive Security Framework"
	@echo ""
	@echo "Available targets:"
	@echo "  make install      - Install dependencies and build"
	@echo "  make build        - Build Argus binary"
	@echo "  make build-linux  - Build for Linux"
	@echo "  make build-macos  - Build for macOS"
	@echo "  make build-windows- Build for Windows"
	@echo "  make demo         - Run demo workflow"
	@echo "  make run          - Run Argus (interactive)"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make test         - Run tests"
	@echo "  make docker       - Build Docker image"
	@echo "  make docker-run   - Run in Docker"

install:
	@echo "Installing Argus..."
	go mod download
	$(MAKE) build
	@echo "✓ Installation complete"
	@echo "Run: ./argus demo"

build:
	@echo "Building Argus..."
	CGO_ENABLED=1 go build -o argus ./cmd/argus
	@echo "✓ Build complete: ./argus"

build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o argus-linux ./cmd/argus

build-macos:
	GOOS=darwin GOARCH=amd64 go build -o argus-macos ./cmd/argus

build-windows:
	GOOS=windows GOARCH=amd64 go build -o argus.exe ./cmd/argus

demo: build
	./argus demo

run: build
	./argus

clean:
	rm -f argus argus-linux argus-macos argus.exe
	rm -f argus_memory.db argus_memory.db-shm argus_memory.db-wal
	rm -f argus_audit.jsonl

test:
	go test -v ./...

docker:
	docker build -t argus:latest .

docker-run: docker
	docker run -p 8080:8080 argus:latest demo

.DEFAULT_GOAL := help
