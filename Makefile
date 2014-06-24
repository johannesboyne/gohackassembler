.PHONY: default

test: default
	./test.sh

default:
	go build main.go
	mv main gohack
