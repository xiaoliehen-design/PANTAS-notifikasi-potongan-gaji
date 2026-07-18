.PHONY: fmt test vet build check

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

test:
	go test ./...

vet:
	go vet ./...

build:
	CGO_ENABLED=0 go build -buildvcs=false -trimpath -o bin/pantas ./cmd/server

check: vet test build
