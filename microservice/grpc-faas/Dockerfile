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

FROM 108356037/faas-cli:v3
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" &&\ 
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
WORKDIR /app

COPY --from=builder /build/main /app/main
COPY ./template /app/template
COPY ./tradebot-dockerfiles /app/tradebot-dockerfiles
COPY ./entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh
#CMD ["/app/main"]

#RUN apk update && apk add curl && curl -sSL https://cli.openfaas.com | sh
#RUN apk add git && faas-cli template pull

# FROM ubuntu:20.04

# RUN apt-get update && apt-get install -y curl && curl -sSL https://cli.openfaas.com | sh
# RUN apt-get install -y git && faas-cli template pull

# FROM 108356037/faas-cli:v2

# ENV DEBIAN_FRONTEND noninteractive

# RUN apt-get install -y \
#     apt-transport-https \
#     ca-certificates \
#     curl \
#     gnupg \
#     lsb-release

# RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# RUN echo \
#   "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
#   $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null


# RUN apt-get update && apt-get install -y docker-ce docker-ce-cli containerd.io