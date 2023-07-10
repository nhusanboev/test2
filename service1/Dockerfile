FROM golang:1.16-alpine as build

WORKDIR /app

COPY . .

RUN go build -o service1 .

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/service1 .

CMD ["./service1"]
