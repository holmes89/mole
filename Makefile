.PHONY: build test vet clean

build:
	go build -o bin/mole .

test:
	go test -timeout 120s ./...

vet:
	go vet ./...

clean:
	rm -rf bin/
