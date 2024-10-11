FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o backend-app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/backend-app ./

EXPOSE 8080

CMD ["./backend-app"]

