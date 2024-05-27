# build stage
FROM golang:1.22.2-alpine AS builder

ENV GO111MODULE=on

# install git.
RUN apk update && apk add --no-cache git

RUN mkdir -p /go/src/github.com/prongbang/gokafka-producer
WORKDIR /go/src/github.com/prongbang/gokafka-producer
COPY . .

# Using go mod with go 1.11
RUN go mod vendor

# With go ≥ 1.10
RUN go build -o /go/bin/app cmd/app/main.go

# small image
FROM alpine:3.7

WORKDIR /app
COPY --from=builder /go/bin/app .

ENV TZ=Asia/Bangkok
RUN echo "Asia/Bangkok" > /etc/timezone

# run binary.
ENTRYPOINT ["/app/app", "-env", "production"]