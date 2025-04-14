.PHONY: repl docs
run:
	docker run -ti -v $$(pwd):/abs --net host -v ~/.abs_history:/root/.abs_history --name abs --rm abs
fmt:
	go fmt ./...
build:
	docker build -t abs .
bench:
	CONTEXT=abs go test `go list -buildvcs=false ./... | grep -v "/js"` -bench=.
test_all: bench test
test:
	# we don't want to test the JS package
	CONTEXT=abs go test `go list -buildvcs=false ./... | grep -v "/js"`
test_verbose:
	# this will show successful error [line:col] tests per #38
	CONTEXT=abs go test `go list -buildvcs=false ./... | grep -v "/js"` -v
repl:
	go run main.go
build_simple:
	go build -o builds/abs main.go
release: build_simple
	./builds/abs ./scripts/release.abs
docs:
	cd docs && npm i && NODE_OPTIONS=--openssl-legacy-provider npm run dev
build_docs:
	cd docs && npm i && NODE_OPTIONS=--openssl-legacy-provider npm run build
wasm:
	GOOS=js GOARCH=wasm go build -o docs/abs.wasm js/js.go
tapes: build_simple
	docker build -t abs-tapes docs/vhs
	docker run -ti -v $$(pwd)/builds/abs:/usr/bin/abs -v $$(pwd):/abs abs-tapes
