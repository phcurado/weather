.PHONY: fmt vet lint test build check tidy

fmt:
	gofmt -w .

vet:
	go vet ./...

lint:
	golangci-lint run

test:
	go test -race ./...

build:
	go build ./...

tidy:
	go mod tidy

check: fmt vet lint test
