.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	go build -o unique-visitors .

.PHONY: test
test:
	go test -v -parallel=4 ./...