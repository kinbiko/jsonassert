LINTER_VERSION := v2.13.2

.PHONY: check
check: lint test

.PHONY: get-deps
get-deps:
	go get -v -t -d ./...

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(LINTER_VERSION) run ./...

.PHONY: test
test:
	go test -race -count=1 ./...

.PHONY: coverage
coverage:
	go test -race -v -coverprofile=profile.cov -covermode=atomic ./...
