FROM golang:1.17-alpine AS builder

WORKDIR /app
COPY . .

RUN go build -o pingservice

FROM alpine:3.15 AS runner

WORKDIR /app
COPY --from=builder /app .

CMD [ "./pingservice" ]
