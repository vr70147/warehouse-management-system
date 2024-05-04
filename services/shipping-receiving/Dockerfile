FROM golang:1.18 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o inventory-management .

FROM alpine:latest  

LABEL maintainer="Your Name <your.email@example.com>"

WORKDIR /root/

COPY --from=builder /app/inventory-management .

EXPOSE 8080

CMD ["./inventory-management"]
