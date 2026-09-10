.PHONY: dev run build assets check test clean

assets:
	bun run build

run:	assets
	go run ./cmd/firefly

dev:	assets
	@bun run watch & \
		pid=$$!; \
		trap 'kill $$pid 2>/dev/null || true' EXIT INT TERM; \
		go run ./cmd/firefly

build:	assets
	mkdir -p bin
	go build -o bin/firefly ./cmd/firefly

check:
	gofmt -w .
	go vet ./...
	bun run check
	bun run check:ts7

test:
	go test ./...

clean:
	rm -rf bin web/dist
