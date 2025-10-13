# Guide: Modern Go Development with `go.work` and Docker

This guide provides a comprehensive walkthrough for setting up a multi-module Go project using `go.work` and containerizing it with Docker.
This approach is the modern standard for Go development in a monorepo structure, offering a seamless development experience and efficient container builds.

## The "Why": Benefits of This Approach

- **Simplified Local Development:** `go.work` allows all your modules to be treated as a single unit during development,
    so changes in the `common` module are instantly available to `log_output` and `ping_pong` without any complex `replace` directives.
- **IDE Harmony:** Go language servers and IDEs have excellent support for `go.work`, which resolves many of the "missing metadata" errors that can occur with `replace` directives.
- **Efficient & Reproducible Builds:** The Docker setup is designed to be efficient by caching dependencies and creating small, secure production images.
- **Clean `go.mod` Files:** Your module's `go.mod` files are no longer cluttered with `replace` directives for local development.

---

## Part 1: Setting Up the Go Workspace with `go.work`

In this part, we will create a Go workspace, which is the modern replacement for using `replace` directives for local development.

### Step 1: Initialize the Workspace

**What:** We will create a `go.work` file at the root of our project.

**Why:** The `go.work` file tells the Go toolchain that this is a multi-module workspace. It allows us to work on multiple modules simultaneously as if they were a single module.

**How:**

Run the following command in the root directory of your project:

```bash
go work init ./common ./log_output ./ping_pong
```

This command creates a `go.work` file with the following content:

```go
go 1.25.1 // Or your Go version

use (
    ./common
    ./log_output
    ./ping_pong
)
```

### Step 2: Clean Up `go.mod` Files

**What:** We will remove the `replace` directives from our `go.mod` files.

**Why:** The `go.work` file now manages the relationship between our local modules, so the `replace` directives are no longer needed and can be removed to keep the `go.od` files clean.

**How:**

Remove the following lines from `log_output/go.mod` and `ping_pong/go.mod`:

```diff
- replace common => ../common
```

Then, run `go mod tidy` in each service's directory to apply the changes:

```bash
cd log_output
go mod tidy
cd ../ping_pong
go mod tidy
```

---

## Part 2: Containerizing the Multi-Module Application

Now that our Go workspace is set up, we need to adapt our Docker configuration to work with this multi-module structure.

### Step 1: The Challenge with Docker's Build Context

-   **What:** The "build context" is the set of files that Docker has access to during an image build. By default, it's the directory containing the `Dockerfile`.
-   **Why it's a problem:** If we build from within the `log_output` directory, Docker cannot see the `common` module because it's in a parent directory (`../common`).

### Step 2: The Solution - A Unified Build Context

-   **What:** We will set the build context to the root of the project for all services.
-   **Why it works:** This makes all modules (`common`, `log_output`, `ping_pong`) and the `go.work` file visible to Docker during the build process.

### Step 3: Modifying `docker-compose.yml`

**What:** We will update our `docker-compose.yml` to use the project root as the build context and to specify the location of each `Dockerfile`.

**Why:** This orchestrates the build process from a single, unified context.

**How:**

```yaml
# docker-compose.yml

services:
  log_output:
    image: log_output_img
    container_name: log_output_ctr
    build:
      # Set the context to the project root
      context: .
      # Specify the path to the Dockerfile relative to the context
      dockerfile: ./log_output/Dockerfile
    restart: unless-stopped
    ports:
      - "8091:8091"
    environment:
      - PORT=8091

  ping_pong:
    image: ping_pong_img
    container_name: ping_pong_ctr
    build:
      # Set the context to the project root
      context: .
      # Specify the path to the Dockerfile relative to the context
      dockerfile: ./ping_pong/Dockerfile
    restart: unless-stopped
    ports:
      - "8092:8092"
    environment:
      - PORT=8092
```

### Step 4: Creating a Multi-Stage `Dockerfile`

**What:** We will use a multi-stage `Dockerfile` for each service. This is a best practice for building Go applications.

**Why:** Multi-stage builds create small, secure, and efficient production images. The first stage (the `builder`) compiles the code,
and the final stage copies only the compiled binary, leaving behind the source code and build tools.

**How:**

Here is the new `log_output/Dockerfile`. The `ping_pong/Dockerfile` will be very similar.

```dockerfile
# log_output/Dockerfile

# --- Build Stage ---
# Use the official Go image as a builder.
# The 'AS builder' names this stage so we can reference it later.
FROM golang:1.25.1-alpine AS builder

# Set the working directory inside the container.
WORKDIR /app

# Copy the go.work file to resolve dependencies for the entire workspace.
COPY go.work ./

# Copy the go.mod files for each module to leverage Docker's layer caching.
# This step will only be re-run if a go.mod file changes.
COPY common/go.mod ./common/
COPY log_output/go.mod ./log_output/
COPY ping_pong/go.mod ./ping_pong/

# Download the dependencies for all modules in the workspace.
# 'go work vendor' is a clean way to do this.
RUN go work vendor

# Copy the source code for the modules needed for this build.
# This is done after dependency resolution to optimize caching.
COPY common/ ./common/
COPY log_output/ ./log_output/

# Build the application.
# CGO_ENABLED=0 creates a statically linked binary.
# -o specifies the output file name.
RUN CGO_ENABLED=0 go build -o /app/main ./log_output/cmd/api

# --- Final Stage ---
# Use a minimal base image for the final container.
# 'alpine' is a popular choice for its small size.
FROM alpine:latest

WORKDIR /app

# Copy only the compiled binary from the 'builder' stage.
# This is the key to a small production image.
COPY --from=builder /app/main .

# Expose the port that the application will listen on.
EXPOSE 8091

# The command to run when the container starts.
CMD ["/app/main"]
```

---

## Conclusion

You have now successfully configured your project to use a modern Go workspace with `go.work` and have adapted your Docker setup to build efficient and small container images for your multi-module application.

This setup provides a superior development experience and a robust foundation for building and deploying your microservices.

---

## Alternative Setup: Using `replace` Directives (Without `go.work`)

While `go.work` is the recommended approach for Go 1.18+, there might be scenarios where you need to use `replace` directives, such as in older CI/CD environments or with legacy tooling.

### Step 1: `go.mod` Configuration

**What:** We will add a `replace` directive to the `go.mod` file of each service that depends on the `common` module.

**Why:** The `replace` directive tells the Go compiler to use a local path for a module instead of fetching it from a remote source. This is necessary for the Go compiler inside the Docker container to find the `common` module.

**How:**

Add the following lines to `log_output/go.mod` and `ping_pong/go.mod`:

```go
replace common => ../common

require common v0.0.0-00010101000000-000000000000
```

### Step 2: `docker-compose.yml` (No Changes)

The `docker-compose.yml` file remains exactly the same as in the `go.work` setup. We still need to use the project root as the build context to make the `common` module available to the Docker build.

### Step 3: Dockerfile Modifications

**What:** We will modify the Dockerfiles to work with the `replace` directive instead of `go.work`.

**Why:** Without a `go.work` file, we need a different strategy to resolve and vendor our dependencies within the Docker build.

**How:**

#### `log_output/Dockerfile`

```dockerfile
# log_output/Dockerfile (using replace)

# --- Build Stage ---
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Copy the go.mod and go.sum files for the service and the common module.
# The paths are relative to the build context (project root).
COPY common/go.mod ./common/
COPY log_output/go.mod ./log_output/

# Copy the source code for the common module.
COPY common/ ./common/

# Copy the source code for the log_output module.
COPY log_output/ ./log_output/

# Set the working directory to the service we want to build.
# This is important because 'go mod download' and 'go build' are module-aware.
WORKDIR /app/log_output

# Download dependencies for the log_output module.
# The 'replace' directive in go.mod will be used to find the common module.
RUN go mod download

# Build the application.
RUN CGO_ENABLED=0 go build -o /app/main ./cmd/api

# --- Final Stage ---
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage.
COPY --from=builder /app/log_output/main .

EXPOSE 8091

CMD ["/app/main"]
```

#### `ping_pong/Dockerfile`

```dockerfile
# ping_pong/Dockerfile (using replace)

# --- Build Stage ---
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Copy the go.mod and go.sum files for the service and the common module.
COPY common/go.mod ./common/
COPY ping_pong/go.mod ./ping_pong/

# Copy the source code for the common module.
COPY common/ ./common/

# Copy the source code for the ping_pong module.
COPY ping_pong/ ./ping_pong/

# Set the working directory to the service we want to build.
WORKDIR /app/ping_pong

# Download dependencies for the ping_pong module.
RUN go mod download

# Build the application.
RUN CGO_ENABLED=0 go build -o /app/main ./cmd/api

# --- Final Stage ---
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage.
COPY --from=builder /app/ping_pong/main .

EXPOSE 8092

CMD ["/app/main"]
```