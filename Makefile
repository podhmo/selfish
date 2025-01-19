default:
	make build
	make dry-run

build:
	mkdir -p bin
	go build -o bin/selfish .

test:
	go vet ./...
	go test -cover ./...

dry-run: 
	go run . --client fake _examples/*

.PHONY: dry-run default build test
