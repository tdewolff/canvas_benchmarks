SHELL=/usr/bin/env bash

all: clean
	go run prepare.go
	cd ajohnson1; go run run.go; cd ..
	cd ajohnson2; go build; ./ajohnson2; cd ..
	cd tdewolff; go run run.go; cd ..
	cd ioverlay; cargo run; cd ..
	go run draw.go

clean:
	rm -rf data

.PHONY: clean
