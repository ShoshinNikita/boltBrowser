default: build run

BUILD_PATH=./bin/boltBrowser

build:
	go build -o $(BUILD_PATH) cmd/boltBrowser/main.go

run:
	$(BUILD_PATH)

test:
	go test ./...

# lint runs golangci-lint - https://github.com/golangci/golangci-lint
#
# Use go cache to speed up execution: https://github.com/golangci/golangci-lint/issues/1004
#
lint:
	@ echo "Run golangci-lint..."
	@ docker run --rm -it --network=none \
		-v $(shell go env GOCACHE):/cache/go \
		-e GOCACHE=/cache/go \
		-e GOLANGCI_LINT_CACHE=/cache/go \
		-v $(shell go env GOPATH)/pkg:/go/pkg \
		-v $(shell pwd):/app \
		-w /app \
		golangci/golangci-lint:v1.41-alpine golangci-lint run --config .golangci.yml

check: build lint test
