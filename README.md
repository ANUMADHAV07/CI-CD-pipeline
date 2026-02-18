# Kubernetes Deployment with Jenkins CI/CD

A simple Go application demonstrating Kubernetes deployment with Jenkins CI/CD pipeline.

## Project Structure

```
.
├── main.go              # Go web application
├── go.mod               # Go module file
├── Dockerfile           # Multi-stage Docker build
├── Jenkinsfile          # Jenkins CI/CD pipeline
└── k8s/
    ├── deployment.yaml  # Kubernetes deployment manifest
    └── service.yaml     # Kubernetes service manifest
```

## Prerequisites

- Go 1.21+
- Docker Desktop (running)
- Kubernetes cluster (minikube, kind, or cloud)
- Jenkins with Kubernetes plugin

## Quick Start

### 1. Start Kubernetes Cluster

**If using minikube:**
```bash
# Make sure Docker Desktop is running first!
minikube start

# Verify cluster is running
kubectl get nodes
```

**If you see connection errors**, see [SETUP.md](SETUP.md) for troubleshooting.

### 2. Build and Test Locally

```bash
# Run the app locally
go run main.go

# In another terminal, test it
curl http://localhost:3000
```

### 3. Build Docker Image

```bash
# Build image
docker build -t k8s-jenkins-app:latest .

# Test locally
docker run -p 3000:3000 k8s-jenkins-app:latest
```

### 4. Deploy to Kubernetes

```bash
# Deploy application
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml

# Check status
kubectl get pods
kubectl get services

# Access the app
minikube service k8s-jenkins-app-service
# Or: kubectl port-forward service/k8s-jenkins-app-service 3000:80
```

## Local Development

1. Run the application:
```bash
go run main.go
```

2. Test the endpoints:
```bash
curl http://localhost:3000
curl http://localhost:3000/health
```

## Building Docker Image

```bash
docker build -t k8s-jenkins-app:latest .
docker run -p 3000:3000 k8s-jenkins-app:latest
```

## Kubernetes Deployment

### Manual Deployment

1. Apply Kubernetes manifests:
```bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

2. Check deployment status:
```bash
kubectl get deployments
kubectl get pods
kubectl get services
```

3. Access the application:
```bash
# Using NodePort (port 30080)
curl http://localhost:30080

# Or port-forward
kubectl port-forward service/k8s-jenkins-app-service 3000:80
```

## Jenkins CI/CD Setup

### Jenkins Configuration

1. Install required plugins:
   - Kubernetes Plugin
   - Docker Pipeline Plugin
   - Git Plugin

2. Configure Kubernetes credentials in Jenkins:
   - Go to Jenkins → Manage Jenkins → Credentials
   - Add Kubernetes credentials (kubeconfig file or service account)

3. Create a new Pipeline job:
   - Select "Pipeline" job type
   - Configure SCM (Git repository)
   - Set Pipeline script from SCM
   - Point to Jenkinsfile

### Pipeline Stages

1. **Checkout**: Clones the repository
2. **Build**: Compiles the Go application
3. **Test**: Runs Go tests
4. **Build Docker Image**: Creates Docker image
5. **Deploy to Kubernetes**: Applies Kubernetes manifests

### Running the Pipeline

1. Push code to your Git repository
2. Trigger Jenkins pipeline (manual or webhook)
3. Monitor pipeline execution in Jenkins console
4. Verify deployment:
```bash
kubectl get pods -w
kubectl logs -f deployment/k8s-jenkins-app
```

## Updating the Application

1. Make changes to `main.go`
2. Commit and push to repository
3. Jenkins pipeline will automatically:
   - Build new Docker image
   - Deploy updated version to Kubernetes
   - Perform rolling update

## Troubleshooting

- Check pod logs: `kubectl logs <pod-name>`
- Check deployment status: `kubectl describe deployment k8s-jenkins-app`
- Check service: `kubectl describe service k8s-jenkins-app-service`
- View Jenkins console output for pipeline errors

## Next Steps

- Add Ingress controller for external access
- Configure secrets for environment variables
- Add monitoring and logging (Prometheus, Grafana)
- Implement blue-green or canary deployments
- Add automated testing in CI/CD pipeline
