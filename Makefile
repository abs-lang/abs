.PHONY: repl docs
run:
	docker run -ti -v $$(pwd):/abs -v ~/.abs_history:/root/.abs_history --name abs --rm abs
fmt:
	go fmt ./...
build:
	docker build -t abs .
bench:
	CONTEXT=abs go test `go list ./... | grep -v "/js"` -bench=.
test_all: bench test
test:
	# we don't want to test the JS package
	CONTEXT=abs go test `go list ./... | grep -v "/js"`
test_verbose:
	# this will show successful error [line:col] tests per #38
	CONTEXT=abs go test `go list ./... | grep -v "/js"` -v
repl:
	go run main.go
build_simple:
	go build -o builds/abs main.go
release: build_simple
	./builds/abs ./scripts/release.abs
travis: test
docs:
	cd docs && bundle exec jekyll serve
wasm:
	GOOS=js GOARCH=wasm go build -o docs/abs.wasm js/js.go
