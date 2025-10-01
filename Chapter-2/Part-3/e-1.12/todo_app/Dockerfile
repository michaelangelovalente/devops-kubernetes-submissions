FROM golang:1.24.5-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

FROM alpine:3.20.1 AS dev
WORKDIR /app
COPY --from=build /app/main /app/main

ENV IMAGE_DIR=${IMAGE_DIR:-./tmp/images}

EXPOSE ${APP_PORT:-3005}
CMD ["./main"]


