# syntax=docker/dockerfile:1.19

############################
# Etapa de build ARM64
############################
FROM golang:1.25.4 AS builder
WORKDIR /src

ENV GOPROXY=https://proxy.golang.org,direct

# 1) Deps (capa estable + cache)
COPY go.mod go.sum ./
RUN go mod download

# 2) Codi
COPY . .

# 3) Build (cache de compilació)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /out/pinger-exporter ./cmd/pinger-exporter

# --- runtime ---
FROM alpine:3.22
COPY --from=builder /out/pinger-exporter /usr/local/bin/pinger-exporter

ENTRYPOINT ["/usr/local/bin/pinger-exporter"]
