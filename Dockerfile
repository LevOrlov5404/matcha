FROM golang:1.16 AS builder

RUN go version

COPY . /builder/
WORKDIR /builder/

RUN go mod download
RUN GOOS=linux go build -o ./.bin/matcha ./cmd/matcha/main.go

FROM debian:bullseye-slim
RUN apt-get update
RUN apt-get install -y ca-certificates && update-ca-certificates

WORKDIR /app/

COPY --from=builder /builder/.bin/matcha .
COPY --from=builder /builder/configs configs/

CMD ["./matcha"]
