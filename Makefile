LINTER_VERSION := v1.61.0

.PHONY: check
check: lint test

.PHONY: get-deps
get-deps:
	go get -v -t -d ./...

.PHONY: lint
lint: ./bin/linter
	./bin/linter run ./...

.PHONY: test
test:
	go test -race -count=1 ./...

.PHONY: coverage
coverage:
	go test -race -v -coverprofile=profile.cov -covermode=atomic ./...

bin/linter: Makefile
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin $(LINTER_VERSION)
	mv ./bin/golangci-lint ./bin/linter
