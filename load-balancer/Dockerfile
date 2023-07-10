# Dockerfile for load balancer

FROM golang:1.16-alpine

WORKDIR /app

COPY config.json .  
COPY . .

RUN go build -o load-balancer

EXPOSE 8080

CMD ["./load-balancer"]
