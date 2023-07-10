FROM golang:1.16-alpine as build

WORKDIR /app

COPY . .

RUN go build -o service3 .

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/service3 .

CMD ["./service3"]
