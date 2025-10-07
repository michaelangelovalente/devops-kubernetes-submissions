# Exercise 1.11 - Persisting Data
Two Go applications sharing data using Kubernetes Persistent Volumes for the DevOps with Kubernetes Course.

## 🔗 Links

- **Log Writer Docker Hub:** [michaelangelovalente/log_writer:ex1.11](https://hub.docker.com/r/michaelangelovalente/log_writer)
- **Log Reader Docker Hub:** [michaelangelovalente/log_reader:ex1.11](https://hub.docker.com/r/michaelangelovalente/log_reader)
- **Ping Pong Docker Hub:** [michaelangelovalente/ping_pong:ex1.11](https://hub.docker.com/r/michaelangelovalente/ping_pong)
- **Source Code:** [GitHub Repository](https://github.com/michaelangelovalente/devops-kubernetes-submissions/tree/main/Chapter-2/Part-3/e-1.11)

## 📋 Description
This exercise demonstrates how to persist and share data between applications in Kubernetes using a `PersistentVolume` and a `PersistentVolumeClaim`. The "log-output" and "ping-pong" applications from the previous exercise are now deployed in a single multi-container pod, sharing a single storage volume.

### Log Writer Application
- **Functionality:** Periodically writes a timestamp and a random UUID to a `logs.txt` file located in the shared persistent volume.
- **Backend:** Go application running as a background service.

### Ping Pong Application
- **Path:** `/pingpong`
- **Functionality:** Increments a counter on each request. The counter's value is saved to a `pingpong.txt` file in the shared volume.
- **Backend:** Go HTTP server.

### Log Reader Application
- **Path:** `/`
- **Functionality:** Reads the latest entry from `logs.txt` and the current count from `pingpong.txt` from the shared volume and displays them together.
- **Backend:** Go HTTP server.

## 🚀 Features

### Shared Infrastructure
- **Persistent Data:** A `PersistentVolume` (PV) of type `local` is used to store data on the node's filesystem.
- **Persistent Volume Claim:** A `PersistentVolumeClaim` (PVC) is used to request and bind to the PV.
- **Shared Volume:** The PVC is mounted as a shared volume into all three containers in the pod, allowing them to share files.
- **Multi-container Pod:** A single `Deployment` manages one pod that runs the `log-writer`, `log-reader`, and `ping-pong` containers.
- **Shared Ingress:** A single Ingress resource routes external traffic to the `log-reader` and `ping-pong` services.

## 🐳 Docker Commands

### Run Locally with Docker Compose
The services can be run locally for development using Docker Compose, which simulates the shared volume.
```bash
# From the e-1.11/log_output_ping_pong/ directory
docker-compose up --build
```

## ☸️ Kubernetes Deployment

### 1. Create k3d Cluster
```bash
k3d cluster create --port 8081:80@loadbalancer --agents 1
```

### 2. Prepare the Node for Local Persistent Volume
Because we are using a `local` PersistentVolume, the directory must be manually created on the node.
```bash
# Create the directory inside the k3d agent node container
docker exec k3d-k3s-default-agent-0 mkdir -p /tmp/kube
```

### 3. Deploy Applications
```bash
# Deploy all resources: PV, PVC, Deployment, Services, and Ingress
kubectl apply -f manifests/
```

### 4. Verify Deployment
```bash
# Check all resources
kubectl get pv,pvc,deployment,pods,svc,ing

# View logs from the multi-container pod (use the actual pod name)
kubectl logs -f <pod-name> -c log-reader-ctr
kubectl logs -f <pod-name> -c log-writer-ctr
kubectl logs -f <pod-name> -c ping-pong-ctr
```

### 5. Access Applications
```bash
# Increment the ping-pong counter
curl localhost:8081/pingpong

# View the combined output from the log-reader
curl localhost:8081/
```

## 📊 Expected Responses

### Ping Pong Application (`/pingpong`)
Each request increments the counter.
```
Result: 1
```
```
Result: 2
```

### Log Reader Application (`/`)
Displays the latest log from the writer and the current count from ping-pong.
```
2025-09-29T10:00:05Z: 1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d
Ping / Pongs: 2
```

## 🛠️ Local Development

### Live Reload with Air
You can run each service with live-reloading for development. Open a separate terminal for each service.
```bash
# In log_writer directory
make watch

# In log_reader directory
make watch

# In ping_pong directory
make watch
```

## ⚙️ Configuration

### Environment Variables
- **log_writer:**
  - `WRITER_PORT`: Server port (default: 8091)
  - `LOG_FILE_PATH`: Path to the log file in the shared volume.
- **log_reader:**
  - `READER_PORT`: Server port (default: 8092)
  - `LOG_FILE_PATH`: Path to the log file in the shared volume.
  - `PING_PONG_FILE_PATH`: Path to the ping-pong file in the shared volume.
- **ping_pong:**
  - `PING_PONG_PORT`: Server port (default: 8093)
  - `PING_PONG_FILE_PATH`: Path to the ping-pong file in the shared volume.

### Kubernetes Resources
- **PersistentVolume (`ping-pong-pv`):** A 1Gi `local` volume that uses the `/tmp/kube` directory on the node.
- **PersistentVolumeClaim (`ping-pong-claim`):** A request for 1Gi of storage that binds to the `ping-pong-pv`.
- **Deployment (`ping-pong-dep`):** A single deployment that manages a pod with three containers (`log-writer-ctr`, `log-reader-ctr`, `ping-pong-ctr`).
- **Services (`log-reader-svc`, `ping-pong-svc`):** `ClusterIP` services that provide internal endpoints for the `log-reader` and `ping-pong` containers.
- **Ingress (`ping-pong-ingress`):** Routes external traffic to the correct service based on the URL path.

### Network Architecture
```
Internet → k3d LoadBalancer:8081 → Shared Ingress:80
                                      ├─ / → log-reader-svc:2346 → Pod(log-reader-ctr:8092)
                                      └─ /pingpong → ping-pong-svc:2346 → Pod(ping-pong-ctr:8093)
                                                 ↑
                                         (Shared Persistent Volume)
                                                 ↓
                                      Pod(log-writer-ctr)
```

## 🏗️ Project Structure
```
e-1.11/
└── log_output_ping_pong/
    ├── README.md
    ├── docker-compose.yaml
    ├── .env
    ├── log_reader/
    │   ├── cmd/api/main.go
    │   ├── internal/
    │   ├── Dockerfile
    │   └── go.mod
    ├── log_writer/
    │   ├── cmd/api/main.go
    │   ├── internal/
    │   ├── Dockerfile
    │   └── go.mod
    ├── ping_pong/
    │   ├── cmd/api/main.go
    │   ├── internal/
    │   ├── Dockerfile
    │   └── go.mod
    └── manifests/
        ├── deployment-persistent.yaml
        ├── ingress.yaml
        ├── log-reader-service.yaml
        ├── persistentvolume.yaml
        ├── persistentvolumeclaim.yaml
        └── pingpong-service.yaml
```