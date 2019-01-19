.PHONY: repl
run:
	docker run -tiv $$(pwd):/go/src/github.com/abs-lang/abs --name abs --rm abs
fmt:
	go fmt ./...
build:
	docker build -t abs .
test:
	# The -vet=off is as YOLO as it gets
	CONTEXT='abs' go test ./... -vet=off
test_verbose:
	CONTEXT='abs' go test ./... -v -vet=off
repl:
	go run main.go
build_simple:
	go build -o builds/abs main.go
release: build_simple
	./builds/abs ./scripts/release.abs
install:
	go get github.com/c-bata/go-prompt
	go get -v github.com/mattn/go-colorable
	go get -v github.com/mattn/go-tty
	go get -d -v ./...
travis: install test
