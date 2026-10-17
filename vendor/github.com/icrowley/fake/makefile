.DEFAULT_GOAL = all

name     := fake
package  := github.com/icrowley/$(name)

.PHONY: all
all:

.PHONY: test
test:
	go test -v ./...

.PHONY: bench
bench:: dependencies
	go test        \
           -bench=. -v \
           $(shell glide novendor)

.PHONY: lint
lint:
	golangci-lint --color=always --timeout=120s run ./...

.PHONY: check
check: lint test
