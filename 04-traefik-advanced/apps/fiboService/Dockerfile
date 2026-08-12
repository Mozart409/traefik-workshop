FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o fiboService

FROM alpine:3.21 AS runner
WORKDIR /app
COPY --from=builder /app/fiboService .
EXPOSE 3333
CMD ["./fiboService"]
