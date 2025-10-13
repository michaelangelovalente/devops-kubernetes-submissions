FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Copy the workspace file first to enable workspace mode
COPY go.work ./

# Copy module files to leverage Docker's layer caching.
COPY common/go.mod ./common/
COPY log_output/go.mod ./log_output/
COPY ping_pong/go.mod ./ping_pong/

# Download all dependencies for the workspace
RUN go mod download


# Copy all source code.
# This is done after deps are downloaded to improve caching.
COPY . .

# Build the specific application for this service.
RUN go build -o /app/main ./ping_pong/cmd/api/main.go


FROM alpine:3.20.1 AS final
WORKDIR /app
COPY --from=builder /app/main /app/main

ENV PORT=8082
EXPOSE 8082

CMD ["./main"]
