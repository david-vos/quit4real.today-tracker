FROM golang:1.23 as builder

RUN apt-get update && apt-get install -y gcc

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -v -o quit4real.today

FROM golang:1.23

WORKDIR /app

COPY --from=builder /app/quit4real.today /app/

COPY --from=builder /app/src/db/migrations /app/src/db/migrations

EXPOSE 8080

CMD ["./quit4real.today"]

