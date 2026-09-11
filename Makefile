.PHONY: check
check: lint test

.PHONY: get-deps
get-deps:
	go get -v -t -d ./...

.PHONY: lint
lint:
	go tool -modfile=tools.mod golangci-lint run

.PHONY: test
test:
	go test -race -count=1 ./...

.PHONY: coverage
coverage:
	go test -race -v -coverprofile=profile.cov -covermode=atomic ./...
