# Development stage
FROM golang:1.24.3-alpine3.21 AS builder

WORKDIR /main
ENV GOTOOLCHAIN=go1.24.3
ENV GO111MODULE=on

# Install build dependencies
RUN apk update && \
    apk add --no-cache \
    binutils-gold \
    make \
    gcc \
    g++ \
    git \
    openssh \
    tzdata \
    mysql-client && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

# Install Go tools for development
RUN go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2 && \
    go install -v github.com/volatiletech/sqlboiler/v4@v4.15.0 && \
    go install -v github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@v4.15.0 && \
    go install -v github.com/google/wire/cmd/wire@v0.6.0 && \
    go install -v github.com/swaggo/swag/cmd/swag@v1.16.2 && \
    go install -v github.com/cosmtrek/air@v1.41.0 && \
    go install go.uber.org/mock/mockgen@v0.4.0 && \
    go install -v github.com/go-delve/delve/cmd/dlv@v1.25.1



# Set the default command to run air
CMD ["air"]

 