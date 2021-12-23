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
#RUN go build -o main

FROM 108356037/kube-and-helm:ubuntu AS base
RUN apt-get update
RUN apt-get install -y ca-certificates

WORKDIR /app
COPY --from=builder /build/main /app/main
COPY configs configs
COPY helmChartFiles helmChartFiles

CMD ["/app/main"]


# FROM 108356037/kube-and-helm:v1

# # RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" &&\ 
# #     install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# # RUN curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 &&\
# #     chmod 700 get_helm.sh &&\
# #     ./get_helm.sh

# ENV GO111MODULE=on \
#     CGO_ENABLED=0 \
#     GOOS=linux \
#     GOARCH=amd64

# WORKDIR /app
# # COPY . .
# # RUN go mod download

# COPY main main
# COPY configs configs
# COPY helmChartFiles helmChartFiles
# #RUN go build -a -ldflags '-w -extldflags "-static"' -o main

# #RUN go build -o main 
# #CMD ["go","run","main.go"]
# CMD ["/app/main"]