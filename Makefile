.PHONY: repl
run:
	docker run -ti -v $$(pwd):/abs -v ~/.abs_history:/root/.abs_history --name abs --rm abs
fmt:
	go fmt ./...
build:
	docker build -t abs .
bench:
	go test ./... -bench=.
test_all: bench test
test:
	# The -vet=off is as YOLO as it gets
	CONTEXT=abs go test ./... -vet=off
test_verbose:
	# this will show successful error [line:col] tests per #38
	CONTEXT='abs' go test ./... -v -vet=off
repl:
	go run main.go
build_simple:
	go build -o builds/abs main.go
release: build_simple
	./builds/abs ./scripts/release.abs
travis: test
