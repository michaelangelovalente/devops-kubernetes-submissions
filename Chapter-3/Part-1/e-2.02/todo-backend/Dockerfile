FROM golang:1.25.3-alpine AS build

WORKDIR /app

COPY go.work ./

COPY common/go.mod ./common/
COPY todo-app/go.mod ./todo-app/
COPY todo-backend/go.mod ./todo-backend/

RUN go mod download

COPY . .

RUN go build -o /app/main ./todo-backend/cmd/api/main.go

FROM alpine:3.20.1 AS dev
WORKDIR /app
COPY --from=build /app/main /app/main
EXPOSE 8081
CMD ["./main"]


