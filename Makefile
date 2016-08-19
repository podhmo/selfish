default:
	(cd cmd && go build -o ../bin/selfish)
	make run

run:
	./bin/selfish `cat ./selfish.config` hello.md || echo ok

.PHONY: run default
