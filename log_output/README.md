# 1.1 Project log_output

## Description
 Docker Hub repository: [log-output](https://hub.docker.com/layers/michaelangelovalente/log-output/v1.0/images/sha256-d6f49f4b33d256e739cf9e8a35320af7f6a87dee84d324de739eee89ae549d54)
 Source Code: [log output source code]()
 log_output's main functionality is to log every 5 seconds a random string, this is stored in memory and then outputs it every 5 seconds with a timestamp attached.
 Golang was used to create log_output and it simulates a backend REST application.
 Note: the code base is "complex" based on the exercises main functionality, but was done so to use as a base for other backend applications for the next exercises and
       it was used to study in detail a typical golang backend application that uses go's standard net/http library

### Commands

#### 1 - docker image creation

 - 1.1 Image creation
```bash
docker build -t <>/log-output:v1.0

```

- 1.2 Image push on Docker Hub

```bash
docker push <>/log-output:v1.0
```

#### 2 - Kubernetes k3s/k3d cluster creation and management

- 2.1 Create cluster with 2 agent nodes with cluster name `log-output-cluster`

```bash
k3d cluster create log-output-cluster -a 2
```

- 2.2 Cluster list and context information
```bash
kubectl config get-contexts
```

- 2.3 Set kubectl cluster context to `log-output-cluster`
```bash
kubectl config set-context log-output-cluster
```

- 2.4 kubectl application deployment
```bash
kubectl create deployment log-output-ex-1-1 --image=michaelangelovalente/log-output:v1.0
```

- 2.5 View and confirm application deployment
```bash
kubectl get deployments
```

- 2.6 Check correct output for application logs
```bash
kubectl logs -f log-output-ex-1-1-7b98bd59b9-nctt9
```
---

### Local project setup and usage
These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```
Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```
Live reload the application:
```bash
make watch
```

Clean up binary from the last build:
```bash
make clean
```

<!-- Run build make command with tests -->
<!-- ```bash -->
<!-- make all -->
