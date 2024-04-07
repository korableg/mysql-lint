GO_FILES=$(shell find . -name '*.go')
VERSION?="dev"
NAME="mysql-lint"
LDFLAGS="-s -w -X github.com/korableg/mysql-lint/cmd.Version=${VERSION}"

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

.PHONY: docker_release
docker_release:
	docker build \
		--build-arg VERSION=${VERSION} \
		--build-arg LDFLAGS=${LDFLAGS} \
 		-t korableg/${NAME}:${VERSION} \
 		-t korableg/${NAME}:latest \
 		--push \
 		.