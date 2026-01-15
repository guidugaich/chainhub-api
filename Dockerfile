FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./
RUN go build -o /bin/chainhub-api ./cmd/api

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /bin/chainhub-api /app/chainhub-api

EXPOSE 8080

CMD ["/app/chainhub-api"]
