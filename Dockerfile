FROM golang:1.23.4 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/.

FROM debian:bookworm-slim

COPY --from=builder /app/main /main

EXPOSE 3000

CMD ["/main"]
