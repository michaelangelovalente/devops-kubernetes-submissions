FROM golang:1.25.3-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates tzdata

# Create non-root user for building
RUN addgroup -S appgroup && adduser -S -G appgroup -u 1001 appuser

WORKDIR /app

RUN chown -R appuser:appgroup /app

COPY --chown=appuser:appgroup go.work go.work.sum ./
COPY --chown=appuser:appgroup common/go.mod ./common/
COPY --chown=appuser:appgroup todo-app/go.mod todo-app/go.sum ./todo-app/
COPY --chown=appuser:appgroup todo-backend/go.mod todo-backend/go.sum ./todo-backend/

# Download dependencies as non-root
USER appuser
RUN go mod download

# Copy src code as non-root user
COPY --chown=appuser:appgroup . .

RUN CGO_ENABLED=0 go build -o /app/main ./todo-app/cmd/api/main.go


FROM alpine:3.20 AS runtime

RUN addgroup -S appgroup && adduser -S -G appgroup -u 1001 appuser

RUN mkdir -p /app/tmp/images && \
        chown -R appuser:appgroup /app

# Import from builder
COPY --from=builder --chown=appuser:appgroup /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder --chown=appuser:appgroup /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder --chown=appuser:appgroup /etc/passwd /etc/group /etc/

# Copy binary
COPY --from=builder --chown=appuser:appgroup /app/main /app/main


WORKDIR /app

USER appuser

EXPOSE 3005
ENTRYPOINT ["/app/main"]


