.PHONY: repl
run:
	docker run -tiv $$(pwd):/go/src/abs --name abs --rm abs
fmt:
	go fmt ./...
build:
	docker build -t abs .
test:
	go test ./... -vet=off
repl:
	go run main.go
build_simple:
	go build -o builds/abs main.go
release:
	./scripts/release.sh
