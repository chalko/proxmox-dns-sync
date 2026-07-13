# Build stage
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o proxmox-dns-sync ./cmd/proxmox-dns-sync

# Final stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/proxmox-dns-sync /app/proxmox-dns-sync
ENTRYPOINT ["/app/proxmox-dns-sync"]
