FROM golang:1.25.3-alpine AS build

WORKDIR /app

COPY go.work ./

COPY common/go.mod ./common/
COPY todo-app/go.mod ./todo-app/
COPY todo-backend/go.mod ./todo-backend/

RUN go mod download

COPY . .

RUN go build -o /app/main ./todo-app/cmd/api/main.go

FROM alpine:3.20.1 AS dev
WORKDIR /app
COPY --from=build /app/main /app/main

ENV IMAGE_DIR=${IMAGE_DIR:-/app/tmp/images} PICSUM_URL=https://picsum.photos/1200 TODO_SVC_URL=http://localhost:8080
EXPOSE 3005
CMD ["./main"]


