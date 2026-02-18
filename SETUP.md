# Kubernetes Cluster Setup Guide

## Problem: Connection Refused Error

If you see this error:
```
error: error validating "k8s/deployment.yaml": error validating data: failed to download openapi: Get "https://127.0.0.1:56234/openapi/v2?timeout=32s": dial tcp 127.0.0.1:56234: connect: connection refused
```

**This means your Kubernetes cluster is not running.**

---

## Solution: Start Minikube Cluster

You have minikube installed. Follow these steps:

### Step 1: Start Docker Desktop

1. Open Docker Desktop application
2. Wait for it to fully start (whale icon should be steady)
3. Verify Docker is running:
   ```bash
   docker ps
   ```

### Step 2: Start Minikube

```bash
minikube start
```

This will:
- Download Kubernetes images (first time only)
- Create a virtual machine/container
- Start the Kubernetes cluster
- Configure kubectl to use minikube

**Expected output:**
```
😄  minikube v1.x.x on Darwin
✨  Using the docker driver based on existing profile
👍  Starting control plane node minikube in cluster minikube
🚜  Pulling base image ...
🔥  Creating docker container (CPUs=2, Memory=4000MB) ...
🐳  Preparing Kubernetes v1.x.x on Docker ...
    ▪ Generating certificates and keys ...
    ▪ Booting control plane ...
    ▪ Configuring RBAC rules ...
✅  Verifying Kubernetes components...
    ▪ Using image gcr.io/k8s-minikube/storage-provisioner:v5
🌟  Enabled addons: storage-provisioner, default-storageclass
🏄  Done! kubectl is now configured to use "minikube" cluster
```

### Step 3: Verify Cluster is Running

```bash
# Check cluster status
minikube status

# Check nodes
kubectl get nodes

# Check cluster info
kubectl cluster-info
```

**Expected output:**
```
kubectl get nodes
NAME       STATUS   ROLES           AGE   VERSION
minikube   Ready    control-plane   1m    v1.x.x
```

### Step 4: Now Deploy Your Application

```bash
# Apply deployment
kubectl apply -f k8s/deployment.yaml

# Apply service
kubectl apply -f k8s/service.yaml

# Check pods
kubectl get pods

# Check services
kubectl get services
```

---

## Common Issues and Solutions

### Issue 1: Docker Permission Denied

**Error:**
```
permission denied while trying to connect to the Docker daemon socket
```

**Solution:**
1. Make sure Docker Desktop is running
2. Add your user to docker group (Linux) or restart Docker Desktop (Mac/Windows)
3. Try: `minikube delete` then `minikube start`

### Issue 2: Minikube Won't Start

**Try:**
```bash
# Delete existing cluster
minikube delete

# Start fresh
minikube start

# If still failing, try with specific driver
minikube start --driver=docker
```

### Issue 3: Port Already in Use

**Error:**
```
Error: port 56234 is already in use
```

**Solution:**
```bash
# Find what's using the port
lsof -i :56234

# Kill the process or restart minikube
minikube stop
minikube start
```

### Issue 4: Out of Memory

**Error:**
```
Error: insufficient memory
```

**Solution:**
```bash
# Start with less memory
minikube start --memory=2048

# Or increase Docker Desktop memory limit in settings
```

---

## Quick Commands Reference

```bash
# Start cluster
minikube start

# Stop cluster
minikube stop

# Delete cluster
minikube delete

# Check status
minikube status

# View dashboard (optional)
minikube dashboard

# Get cluster IP
minikube ip

# Access service URL
minikube service k8s-jenkins-app-service --url
```

---

## Alternative: Using Kind (Kubernetes in Docker)

If minikube doesn't work, you can use Kind:

```bash
# Install kind
brew install kind

# Create cluster
kind create cluster --name k8s-jenkins

# Verify
kubectl get nodes
```

---

## Next Steps After Cluster is Running

1. **Deploy your app:**
   ```bash
   kubectl apply -f k8s/deployment.yaml
   kubectl apply -f k8s/service.yaml
   ```

2. **Check deployment:**
   ```bash
   kubectl get pods -w
   kubectl get services
   ```

3. **Access your app:**
   ```bash
   # Using minikube service
   minikube service k8s-jenkins-app-service
   
   # Or port-forward
   kubectl port-forward service/k8s-jenkins-app-service 3000:80
   ```

4. **View logs:**
   ```bash
   kubectl logs -f deployment/k8s-jenkins-app
   ```

---

## Verify Everything Works

Run these commands to verify your setup:

```bash
# 1. Cluster is running
kubectl get nodes

# 2. Deploy application
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml

# 3. Wait for pods to be ready
kubectl wait --for=condition=ready pod -l app=k8s-jenkins-app --timeout=60s

# 4. Check pods are running
kubectl get pods

# 5. Test the application
curl $(minikube service k8s-jenkins-app-service --url)
```

If all commands succeed, your Kubernetes cluster is ready! 🎉
