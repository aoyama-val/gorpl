all: gorpl

test: gorpl
	cd test && ./testall

gorpl: gorpl.go
	go build
