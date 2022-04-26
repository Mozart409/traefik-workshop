FROM golang:1.17-alpine AS builder

WORKDIR /app
COPY . .

RUN go build -o jaegerService

FROM alpine:3.15 AS runner

WORKDIR /app
COPY --from=builder /app .

EXPOSE 3334

CMD [ "./jaegerService" ]
