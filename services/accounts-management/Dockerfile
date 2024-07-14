FROM golang:1.22.2 as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.19.1

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]
