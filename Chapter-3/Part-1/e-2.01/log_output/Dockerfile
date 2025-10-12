FROM golang:1.25.1-alpine AS build

WORKDIR /app

COPY go.work ./


COPY common/go.mod ./common/

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

FROM alpine:3.20.1 AS dev
WORKDIR /app
COPY --from=build /app/main /app/main
EXPOSE ${LOG_OUTPUT_PORT}
CMD ["./main"]


