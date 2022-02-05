default:
	make build
	make run

build:
	(cd cmd/selfish && go build -o ../../bin/selfish)

run:
	./bin/selfish examples/data/hello.md || echo ok

test:
	go test -cover

vendor-update:
	dep ensure --update

.PHONY: run default build test
