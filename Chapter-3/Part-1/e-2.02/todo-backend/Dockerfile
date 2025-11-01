FROM golang:1.25.3-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates tzdata

RUN addgroup -S appgroup && adduser -S -G appgroup -u 1001 appuser

WORKDIR /app

RUN chown -R appuser:appgroup /app

COPY --chown=appuser:appgroup go.work go.work.sum ./
COPY --chown=appuser:appgroup common/go.mod ./common/
COPY --chown=appuser:appgroup todo-app/go.mod todo-app/go.sum ./todo-app/
COPY --chown=appuser:appgroup todo-backend/go.mod todo-backend/go.sum ./todo-backend/

USER appuser
RUN go mod download

COPY --chown=appuser:appgroup . .

RUN CGO_ENABLED=0 go build -o /app/main ./todo-backend/cmd/api/main.go


# Final stage - minimal runtime
FROM scratch AS runtime

# Import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/group /etc/

# Copy binary
COPY --from=builder --chown=appuser:appgroup /app/main /app/main

WORKDIR /app

USER appuser

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/app/main", "-health-check"]


EXPOSE 8081
ENTRYPOINT ["/app/main"]


