# Complete Explanation: Kubernetes + Jenkins CI/CD

This document explains each component and why it's needed for deploying applications with Kubernetes and Jenkins CI/CD.

---

## 🎯 **The Big Picture: What Are We Building?**

We're creating an **automated deployment pipeline** that:
1. Takes your code
2. Builds it into a container
3. Tests it
4. Deploys it to Kubernetes automatically

**Why?** Instead of manually building and deploying, Jenkins does it automatically when you push code.

---

## 📝 **Step 1: The Application (`main.go`)**

### What it does:
A simple Go web server with two endpoints:
- `/` - Returns JSON with app info
- `/health` - Health check endpoint

### Why we need it:
- **Real application**: We need something to deploy
- **Health endpoint**: Kubernetes uses this to check if your app is running
- **Environment variables**: Shows how to configure apps in containers

### Key concepts:
```go
port := os.Getenv("PORT")  // Reads from environment
hostname, _ := os.Hostname()  // Gets container hostname
```

**Why environment variables?** Different environments (dev/staging/prod) can use different configs without changing code.

---

## 🐳 **Step 2: Dockerfile (Containerization)**

### What it does:
Converts your Go application into a Docker container image.

### Why we need it:
- **Consistency**: Same app runs the same way everywhere
- **Isolation**: App runs in its own environment
- **Portability**: Run on any machine with Docker

### Multi-stage build explained:

```dockerfile
# Stage 1: Builder
FROM golang:1.21-alpine AS builder
```
**Why?** We need Go compiler to build the app, but we don't need it in final image.

```dockerfile
COPY go.mod go.sum ./
RUN go mod download
```
**Why?** Download dependencies first (Docker caches this layer, faster rebuilds).

```dockerfile
# Stage 2: Runtime
FROM alpine:latest
COPY --from=builder /app/server .
```
**Why?** Final image is tiny (only the compiled binary). Smaller = faster deployments.

### Why multi-stage?
- **Builder image**: ~300MB (has Go compiler)
- **Final image**: ~10MB (only binary)
- **Result**: Faster downloads, less storage, more secure (no compiler tools)

---

## ☸️ **Step 3: Kubernetes Deployment (`k8s/deployment.yaml`)**

### What it does:
Tells Kubernetes how to run your application.

### Why we need it:
- **Orchestration**: Kubernetes manages your containers
- **Scaling**: Run multiple copies automatically
- **Self-healing**: Restarts failed containers
- **Rolling updates**: Update without downtime

### Key parts explained:

```yaml
replicas: 3
```
**Why?** Run 3 copies of your app. If one crashes, others keep serving traffic.

```yaml
resources:
  requests:
    memory: "64Mi"
    cpu: "100m"
  limits:
    memory: "128Mi"
    cpu: "200m"
```
**Why?** 
- **Requests**: Kubernetes reserves this much (guaranteed)
- **Limits**: App can't use more than this (prevents one app from starving others)

```yaml
livenessProbe:
  httpGet:
    path: /health
```
**Why?** Kubernetes checks `/health` every 10 seconds. If it fails, Kubernetes restarts the container.

```yaml
readinessProbe:
  httpGet:
    path: /health
```
**Why?** Kubernetes checks if app is ready to receive traffic. If not ready, traffic is routed to other pods.

**Difference:**
- **Liveness**: "Is the app alive?" → Restart if dead
- **Readiness**: "Can the app handle traffic?" → Don't send traffic if not ready

---

## 🌐 **Step 4: Kubernetes Service (`k8s/service.yaml`)**

### What it does:
Creates a stable network endpoint to access your pods.

### Why we need it:
- **Pod IPs change**: When pods restart, they get new IPs
- **Service IP is stable**: Always points to your pods
- **Load balancing**: Distributes traffic across all pods

### Key parts:

```yaml
type: NodePort
```
**Why?** Makes your app accessible from outside the cluster on port 30080.

```yaml
selector:
  app: k8s-jenkins-app
```
**Why?** Service finds pods with this label and routes traffic to them.

```yaml
port: 80          # External port
targetPort: 3000  # Container port
nodePort: 30080   # Node port (accessible from outside)
```

**Flow:** External request → NodePort 30080 → Service Port 80 → Pod Port 3000

---

## 🔄 **Step 5: Jenkins Pipeline (`Jenkinsfile`)**

### What it does:
Automates the entire build and deployment process.

### Why we need it:
- **Automation**: No manual steps
- **Consistency**: Same process every time
- **Speed**: Deploy in minutes, not hours
- **Reliability**: Tests before deploying

### Pipeline stages explained:

#### Stage 1: Checkout
```groovy
checkout scm
```
**What:** Gets your code from Git repository  
**Why:** Need code to build

#### Stage 2: Build
```groovy
go mod download
go build -o app main.go
```
**What:** Downloads dependencies and compiles Go code  
**Why:** Verify code compiles before creating container

#### Stage 3: Test
```groovy
go test ./...
```
**What:** Runs automated tests  
**Why:** Catch bugs before deploying to production

#### Stage 4: Build Docker Image
```groovy
docker build -t k8s-jenkins-app:${BUILD_NUMBER} .
```
**What:** Creates Docker image with unique tag  
**Why:** 
- Containerize the app
- Tag with build number for version tracking
- Can rollback to specific version if needed

#### Stage 5: Deploy to Kubernetes
```groovy
kubectl set image deployment/k8s-jenkins-app app=k8s-jenkins-app:123
```
**What:** Updates Kubernetes deployment with new image  
**Why:** 
- Rolling update: Kubernetes gradually replaces old pods with new ones
- Zero downtime: Old pods serve traffic while new ones start
- Automatic rollback: If new version fails, Kubernetes reverts

---

## 🔄 **Complete Flow: How Everything Works Together**

```
1. Developer pushes code to Git
   ↓
2. Jenkins detects the push (webhook or polling)
   ↓
3. Jenkins runs pipeline:
   a. Checks out code
   b. Builds Go application
   c. Runs tests
   d. Builds Docker image
   e. Deploys to Kubernetes
   ↓
4. Kubernetes:
   a. Pulls new Docker image
   b. Creates new pods with new image
   c. Health checks pass
   d. Routes traffic to new pods
   e. Terminates old pods
   ↓
5. Application is live with new version!
```

---

## 🎓 **Key Concepts Explained**

### Container vs Pod vs Deployment

- **Container**: Your application running in isolation
- **Pod**: Kubernetes wraps container(s) in a pod (smallest deployable unit)
- **Deployment**: Manages pods (creates, updates, scales)

**Analogy:** 
- Container = Engine
- Pod = Car (has engine)
- Deployment = Fleet management (manages multiple cars)

### Why Kubernetes?

**Without Kubernetes:**
- Manual deployment
- No automatic scaling
- No self-healing
- Difficult to manage multiple servers

**With Kubernetes:**
- Automatic deployment
- Auto-scaling based on load
- Restarts failed containers
- Manages multiple servers as one

### Why CI/CD?

**Without CI/CD:**
```
1. Developer builds locally
2. Tests manually
3. Creates Docker image manually
4. Pushes to registry manually
5. SSH into server
6. Pull image manually
7. Stop old version
8. Start new version
9. Hope nothing breaks
```

**With CI/CD:**
```
1. Developer pushes code
2. Everything else is automatic!
```

---

## 🛠️ **What Each File Does - Quick Reference**

| File | Purpose | Why Needed |
|------|---------|------------|
| `main.go` | Your application code | The actual app to deploy |
| `go.mod` | Go dependencies | Manages libraries your app uses |
| `Dockerfile` | Container definition | Packages app for consistent deployment |
| `k8s/deployment.yaml` | Kubernetes deployment | Tells K8s how to run your app |
| `k8s/service.yaml` | Kubernetes service | Provides stable network access |
| `Jenkinsfile` | CI/CD pipeline | Automates build and deployment |
| `.dockerignore` | Docker ignore file | Excludes files from Docker build |

---

## 🚀 **Next Steps to Learn**

1. **Run locally:**
   ```bash
   go run main.go
   curl http://localhost:3000
   ```

2. **Build Docker image:**
   ```bash
   docker build -t my-app .
   docker run -p 3000:3000 my-app
   ```

3. **Deploy to Kubernetes:**
   ```bash
   kubectl apply -f k8s/deployment.yaml
   kubectl apply -f k8s/service.yaml
   kubectl get pods
   ```

4. **Set up Jenkins:**
   - Install Jenkins
   - Install Kubernetes plugin
   - Create pipeline job
   - Connect to your Git repo

5. **Make a change:**
   - Edit `main.go`
   - Push to Git
   - Watch Jenkins deploy automatically!

---

## 💡 **Common Questions**

**Q: Why not just run the app directly?**  
A: Containers provide isolation, consistency, and portability. Same app runs identically on your laptop and production.

**Q: Why Kubernetes if Docker works?**  
A: Docker runs containers. Kubernetes orchestrates them (scaling, health checks, updates, networking).

**Q: Why Jenkins? Can't I just use kubectl?**  
A: You can, but Jenkins automates it. Every push triggers automatic testing and deployment.

**Q: What if deployment fails?**  
A: Kubernetes automatically rolls back. Jenkins can also send notifications.

**Q: How do I update the app?**  
A: Just push code! Jenkins builds new image, Kubernetes does rolling update automatically.

---

## 📚 **Learning Path**

1. ✅ Understand the application code
2. ✅ Learn Docker basics (build, run, push)
3. ✅ Learn Kubernetes basics (pods, deployments, services)
4. ✅ Set up Jenkins
5. ✅ Connect everything together
6. 🎯 Practice by making changes and watching deployments!
