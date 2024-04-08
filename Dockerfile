FROM golang:1.22.2-alpine3.19 AS build

MAINTAINER Dmitry Titov <dim@titovcode.com>

ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ARG VERSION="latest"
ARG LDFLAGS

RUN go build -o ./bin/app -ldflags "$LDFLAGS" .

FROM golang:1.22.2-alpine3.19

COPY --from=build /app/bin/app /app

ENTRYPOINT ["/app"]