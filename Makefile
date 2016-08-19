default:
	make build
	make run

build:
	(cd cmd/selfish && go build -o ../../bin/selfish)

run:
	./bin/selfish data/hello.md || echo ok

.PHONY: run default build
