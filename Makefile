.PHONY: repl
run:
	docker run -tiv $$(pwd):/go/src/abs --name abs --rm abs
fmt:
	go fmt ./...
build:
	docker build -t abs .
test:
	go test ./...
repl:
	go run main.go
