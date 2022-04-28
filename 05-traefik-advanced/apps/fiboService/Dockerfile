FROM golang:1.17-alpine AS builder

WORKDIR /app
COPY . .

RUN go build -o fiboService

FROM alpine:3.15 AS runner

WORKDIR /app
COPY --from=builder /app .

EXPOSE 3333

CMD [ "./fiboService" ]
