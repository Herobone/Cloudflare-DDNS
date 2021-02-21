FROM golang:1.15 as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o cloudflare-ddns

FROM alpine
ENTRYPOINT ["/cloudflare-ddns"]
COPY --from=builder /build/cloudflare-ddns /cloudflare-ddns
