FROM golang:1.22.2 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o order-processing .

FROM alpine:latest  

LABEL maintainer="Raanan Adut radesigner84@gmail.com"

WORKDIR /root/

COPY --from=builder /app/order-processing .

EXPOSE 8082

CMD ["./main"]
