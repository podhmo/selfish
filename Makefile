default:
	go build -o a.out
	make run

run:
	./a.out `cat ./selfish.config` hello.md || echo ok

.PHONY: run default
