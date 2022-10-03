# syntax = docker/dockerfile:experimental
# Build Container
FROM golang:1.18 as builder

ENV GO111MODULE on
ENV GOPRIVATE=github.com/latonaio
WORKDIR /go/src/github.com/latonaio

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o data-platform-authenticator ./cmd/server

# Runtime Container
FROM alpine:3.14
RUN apk add --no-cache libc6-compat
ENV SERVICE=data-platform-authenticator \
    APP_DIR="${AION_HOME}/${POSITION}/${SERVICE}"

ARG PRIVATE_KEY
ENV PRIVATE_KEY=$PRIVATE_KEY

ARG PUBLIC_KEY
ENV PUBLIC_KEY=$PUBLIC_KEY

WORKDIR ${AION_HOME}

COPY --from=builder /go/src/github.com/latonaio/data-platform-authenticator .

CMD ["./data-platform-authenticator"]
