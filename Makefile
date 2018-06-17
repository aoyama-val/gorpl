all: rpl

test: rpl
	cd test && ./testall

rpl: rpl.go
	go build
