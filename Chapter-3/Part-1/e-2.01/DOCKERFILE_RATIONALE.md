
# Rationale: One Dockerfile Per Service

This document addresses the question: "Can I use a single Dockerfile for multiple microservices in a monorepo?"

While it is technically possible, this guide explains why the best practice is to maintain one Dockerfile per service.

---

## The Initial Observation

In a multi-module Go project, the Docker build process for each service often looks very similar. They share a common build context, copy many of the same files (like a `common` module), and follow similar build steps.
This leads to the natural question of whether a single, "universal" Dockerfile could be used to reduce file duplication.

## The "Universal" Dockerfile: How It Would Work

To have a single Dockerfile, you would need to parameterize it using build arguments (`ARG`). It would look something like this:

```dockerfile
# Dockerfile.universal (not recommended)

# --- Build Stage ---
# Use build arguments to specify which service to build
ARG SERVICE_PATH
ARG EXPOSED_PORT

FROM golang:1.25.1-alpine AS builder
WORKDIR /app

# (Copy all go.mod, go.work, and source files as before)
# ...

# Build the specific application passed in as an argument
RUN CGO_ENABLED=0 go build -o /app/main ./$SERVICE_PATH/cmd/api

# --- Final Stage ---
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .

# Expose the port passed in as an argument
EXPOSE $EXPOSED_PORT

CMD ["/app/main"]
```

Your `docker-compose.yml` would then pass these arguments during the build:

```yaml
# docker-compose.yml (with a universal Dockerfile)
services:
  log_output:
    build:
      context: .
      dockerfile: Dockerfile.universal
      args:
        SERVICE_PATH: log_output
        EXPOSED_PORT: 8091
  ping_pong:
    build:
      context: .
      dockerfile: Dockerfile.universal
      args:
        SERVICE_PATH: ping_pong
        EXPOSED_PORT: 8092
```

## Why Separate Dockerfiles Are Better

While the "universal" approach seems to reduce file duplication, it introduces more problems than it solves:

1.  **Increased Complexity:** The single Dockerfile is now more complex and harder to understand at a glance. The `docker-compose.yml` also becomes more verbose and tightly coupled to the Dockerfile's implementation details.
2.  **Violation of Single Responsibility:** A Dockerfile should have one clear purpose: to build a specific service. A universal Dockerfile is trying to do too many things, which makes it brittle and harder to maintain.
3.  **Reduced Flexibility:** What if, in the future, the `ping_pong` service needs a specific system dependency (like `gcc` or `git`) that `log_output` doesn't?
    You would have to add conditional logic to your universal Dockerfile, making it even more complex. With separate files, you just add `RUN apk add --no-cache gcc` to the `ping_pong/Dockerfile`.
4.  **Less Clarity:** Having a `Dockerfile` inside each service's directory (`log_output/Dockerfile`) makes it immediately obvious how that service is built. A shared `Dockerfile` at the root of the project makes the build process for any given service less transparent.
5.  **Minimal Duplication:** The duplication between your current Dockerfiles is minimal and, more importantly, it's *healthy* duplication. It ensures that the build process for each service is independent and self-contained.

## Conclusion

Your instinct to question the duplication is a good one, as it shows you're thinking about efficiency. However, in this case, the established best practice is to **keep one Dockerfile per service**.

The small amount of duplication is a worthwhile trade-off for the significant gains in clarity, flexibility, and long-term maintainability.
The setup of having one Dockerfile per service, with the build context set to the project root, is the standard and most effective way to handle this common scenario in monorepo development.
