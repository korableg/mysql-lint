GO_FILES=$(shell find . -name '*.go')
NAME="mysql-lint"
LDFLAGS="-s -w"

.PHONY: order_imports
order_imports:
	goimports -v -local github.com/korableg/mysql-lint -w $(GO_FILES)

.PHONY: lint
lint:
	golangci-lint -v run ./...

.PHONY: test
test:
	go test -v --race ./...

.PHONY: build
build:
	go build -o ${NAME} -ldflags ${LDFLAGS} .