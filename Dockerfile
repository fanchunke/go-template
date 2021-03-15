ARG GO_VERSION=1.15
FROM golang:${GO_VERSION}-alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY internal ./internal
COPY main.go .

# Build the application
RUN go build -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

# Build a small image
FROM alpine:latest

RUN mkdir -p /home/works/program/logs

WORKDIR /home/works/program

COPY --from=builder /dist/main .
COPY configs ./configs

EXPOSE 8000

# Command to run
CMD ["./main"]