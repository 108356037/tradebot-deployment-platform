# --- 因為 base container 找不到 configs，所以暫時捨棄 builder 方法 ---- 
FROM golang:1.16 AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download 

COPY . .
RUN go build -ldflags '-w -extldflags "-static"' -o main

FROM ubuntu:20.04 AS base
RUN apt-get update
#RUN apt-get install -y vim curl
WORKDIR /app
COPY --from=builder /build/main /app/main
CMD ["/app/main"]