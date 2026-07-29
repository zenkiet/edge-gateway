FROM golang:alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags="-w -s" \
  -o /app/gateway ./

# Stage 2: Minimal Runtime Image
FROM gcr.io/distroless/static-debian13:nonroot

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/gateway /app/gateway

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/app/gateway"]