# Production-Grade Kubernetes Deployment Setup Report

**Date:** February 18, 2026  
**Project:** Kubernetes + Jenkins CI/CD Pipeline  
**Status:** In Progress

---

## ✅ **COMPLETED SETUP**

### Infrastructure
- [x] AWS EC2 Ubuntu VM created
- [x] Security Group configured (port 8080 for Jenkins)
- [x] SSH access configured

### Jenkins Server
- [x] Jenkins installed and running
- [x] Java 25 (LTS) installed
- [x] Jenkins initial setup completed
- [x] Jenkins accessible at: `http://15.134.38.180:8080`

### Development Tools
- [x] Go 1.21.5 installed
- [x] Git installed
- [x] Docker installed
- [x] kubectl installed

### Application Code
- [x] Go application created (`main.go`)
- [x] Dockerfile created (multi-stage build)
- [x] Kubernetes manifests created
- [x] Jenkinsfile created
- [x] Code pushed to GitHub: `ANUMADHAV07/CI-CD-pipeline`

---

## 🔄 **IN PROGRESS**

- [ ] Jenkins pipeline testing
- [ ] Docker image building
- [ ] Kubernetes cluster setup

---

## 📋 **REMAINING TASKS FOR PRODUCTION**

### Phase 1: Complete CI/CD Pipeline Setup

#### 1.1 Install Missing Tools on Jenkins Server
```bash
# Verify all tools are installed
go version
docker --version
kubectl version --client
java -version

# If any missing, install them
```

#### 1.2 Configure Jenkins Credentials
- [ ] Docker Hub credentials (for pushing images)
- [ ] Kubernetes kubeconfig (for cluster access)
- [ ] GitHub credentials (if private repo)

#### 1.3 Install Jenkins Plugins
Required plugins:
- [ ] Pipeline
- [ ] Docker Pipeline
- [ ] Kubernetes CLI Plugin
- [ ] Credentials Binding
- [ ] Git
- [ ] Build Timeout
- [ ] Timestamper

#### 1.4 Update Jenkinsfile for Production
Add to Jenkinsfile:
- [ ] Docker registry push stage
- [ ] Security scanning stage
- [ ] Environment variables (staging/production)
- [ ] Rollback on failure
- [ ] Notifications (Slack/Email)

---

### Phase 2: Kubernetes Cluster Setup

#### Option A: AWS EKS (Recommended for Production)
```bash
# Install eksctl
curl --silent --location "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/eksctl /usr/local/bin

# Install AWS CLI
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install

# Configure AWS credentials
aws configure

# Create EKS cluster
eksctl create cluster \
  --name production-cluster \
  --region ap-southeast-2 \
  --node-type t3.medium \
  --nodes 2 \
  --nodes-min 2 \
  --nodes-max 4

# Configure kubectl
aws eks update-kubeconfig --name production-cluster --region ap-southeast-2
```

#### Option B: k3s (Lightweight, for Testing)
```bash
# Install k3s
curl -sfL https://get.k3s.io | sh -

# Get kubeconfig
sudo cat /etc/rancher/k3s/k3s.yaml

# Copy to Jenkins user
sudo mkdir -p /var/lib/jenkins/.kube
sudo cp /etc/rancher/k3s/k3s.yaml /var/lib/jenkins/.kube/config
sudo chown -R jenkins:jenkins /var/lib/jenkins/.kube
```

---

### Phase 3: Docker Registry Setup

#### Option A: Docker Hub (Free/Paid)
```bash
# Create account at hub.docker.com
# Get credentials
# Add to Jenkins: Manage Jenkins → Credentials → Add
```

#### Option B: AWS ECR (Recommended for AWS)
```bash
# Create ECR repository
aws ecr create-repository --repository-name k8s-jenkins-app --region ap-southeast-2

# Get login token
aws ecr get-login-password --region ap-southeast-2 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.ap-southeast-2.amazonaws.com

# Add credentials to Jenkins
```

#### Option C: GitHub Container Registry (Free)
```bash
# Use: ghcr.io/ANUMADHAV07/k8s-jenkins-app
# Authenticate with GitHub Personal Access Token
```

---

### Phase 4: Production Kubernetes Manifests

#### 4.1 Create Production Deployment
Create `k8s/deployment.production.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: k8s-jenkins-app
  namespace: production
  labels:
    app: k8s-jenkins-app
    version: "1.0.0"
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: k8s-jenkins-app
  template:
    metadata:
      labels:
        app: k8s-jenkins-app
        version: "1.0.0"
    spec:
      serviceAccountName: app-service-account
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        fsGroup: 1000
      containers:
      - name: app
        image: your-registry.com/k8s-jenkins-app:latest
        imagePullPolicy: Always
        ports:
        - containerPort: 3000
          name: http
        env:
        - name: PORT
          value: "3000"
        - name: APP_VERSION
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: version
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /ready
            port: 3000
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
        startupProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 10
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 10
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - k8s-jenkins-app
              topologyKey: kubernetes.io/hostname
```

#### 4.2 Create Production Service
Create `k8s/service.production.yaml`:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: k8s-jenkins-app-service
  namespace: production
  labels:
    app: k8s-jenkins-app
spec:
  type: LoadBalancer
  selector:
    app: k8s-jenkins-app
  ports:
  - protocol: TCP
    port: 80
    targetPort: 3000
    name: http
```

#### 4.3 Create ConfigMap
Create `k8s/configmap.yaml`:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: production
data:
  version: "1.0.0"
  environment: "production"
```

#### 4.4 Create ServiceAccount
Create `k8s/serviceaccount.yaml`:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: app-service-account
  namespace: production
```

#### 4.5 Create HorizontalPodAutoscaler
Create `k8s/hpa.yaml`:

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: k8s-jenkins-app-hpa
  namespace: production
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: k8s-jenkins-app
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

#### 4.6 Create PodDisruptionBudget
Create `k8s/pdb.yaml`:

```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: k8s-jenkins-app-pdb
  namespace: production
spec:
  minAvailable: 2
  selector:
    matchLabels:
      app: k8s-jenkins-app
```

---

### Phase 5: Security Setup

#### 5.1 Network Policies
Create `k8s/network-policy.yaml`:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: k8s-jenkins-app-network-policy
  namespace: production
spec:
  podSelector:
    matchLabels:
      app: k8s-jenkins-app
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    ports:
    - protocol: TCP
      port: 3000
  egress:
  - to:
    - namespaceSelector:
        matchLabels:
          name: kube-system
    ports:
    - protocol: TCP
      port: 53
```

#### 5.2 RBAC (Role-Based Access Control)
Create `k8s/rbac.yaml`:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: app-role
  namespace: production
rules:
- apiGroups: [""]
  resources: ["pods", "services"]
  verbs: ["get", "list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: app-role-binding
  namespace: production
subjects:
- kind: ServiceAccount
  name: app-service-account
roleRef:
  kind: Role
  name: app-role
  apiGroup: rbac.authorization.k8s.io
```

#### 5.3 Secrets Management
```bash
# Create secrets
kubectl create secret generic app-secrets \
  --from-literal=db-password=secret123 \
  --namespace=production

# Use in deployment
# env:
# - name: DB_PASSWORD
#   valueFrom:
#     secretKeyRef:
#       name: app-secrets
#       key: db-password
```

---

### Phase 6: Monitoring & Observability

#### 6.1 Install Prometheus
```bash
# Add Prometheus Helm repo
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# Install Prometheus
helm install prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace
```

#### 6.2 Install Grafana
```bash
# Grafana comes with Prometheus stack
# Access: kubectl port-forward svc/prometheus-grafana 3000:80 -n monitoring
# Default username: admin
# Get password: kubectl get secret prometheus-grafana -n monitoring -o jsonpath="{.data.admin-password}" | base64 -d
```

#### 6.3 Add ServiceMonitor
Create `k8s/servicemonitor.yaml`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: k8s-jenkins-app-metrics
  namespace: production
spec:
  selector:
    matchLabels:
      app: k8s-jenkins-app
  endpoints:
  - port: http
    path: /metrics
    interval: 30s
```

#### 6.4 Logging (ELK Stack or Loki)
```bash
# Install Loki (lightweight)
helm repo add grafana https://grafana.github.io/helm-charts
helm install loki grafana/loki-stack --namespace logging --create-namespace
```

---

### Phase 7: Ingress & External Access

#### 7.1 Install NGINX Ingress Controller
```bash
# For EKS
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/aws/deploy.yaml

# For k3s (already included)
```

#### 7.2 Create Ingress Resource
Create `k8s/ingress.yaml`:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: k8s-jenkins-app-ingress
  namespace: production
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - app.example.com
    secretName: app-tls
  rules:
  - host: app.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: k8s-jenkins-app-service
            port:
              number: 80
```

---

### Phase 8: Production Jenkinsfile

Update `Jenkinsfile` with production features:

```groovy
pipeline {
    agent any
    
    options {
        buildDiscarder(logRotator(numToKeepStr: '50', daysToKeepStr: '30'))
        timeout(time: 30, unit: 'MINUTES')
        retry(2)
        timestamps()
    }
    
    environment {
        DOCKER_REGISTRY = credentials('docker-registry-url') ?: 'docker.io'
        DOCKER_REPO = 'your-org/k8s-jenkins-app'
        DOCKER_IMAGE = "${DOCKER_REGISTRY}/${DOCKER_REPO}"
        DOCKER_TAG = "${env.BUILD_NUMBER}-${env.GIT_COMMIT.take(7)}"
        KUBERNETES_NAMESPACE = env.BRANCH_NAME == 'main' ? 'production' : 'staging'
        PATH = "/usr/local/go/bin:/usr/local/bin:${env.PATH}"
    }
    
    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        
        stage('Build') {
            steps {
                sh 'go mod download'
                sh 'go build -o app main.go'
            }
        }
        
        stage('Test') {
            steps {
                sh 'go test -v -coverprofile=coverage.out ./...'
                publishTestResults testResultsPattern: 'test-results.xml'
            }
        }
        
        stage('Security Scan') {
            steps {
                sh 'trivy image ${DOCKER_IMAGE}:${DOCKER_TAG} || true'
            }
        }
        
        stage('Build Docker Image') {
            steps {
                sh """
                    docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
                    docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest
                """
            }
        }
        
        stage('Push to Registry') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'docker-registry-credentials',
                    usernameVariable: 'DOCKER_USER',
                    passwordVariable: 'DOCKER_PASS'
                )]) {
                    sh """
                        echo \${DOCKER_PASS} | docker login ${DOCKER_REGISTRY} -u \${DOCKER_USER} --password-stdin
                        docker push ${DOCKER_IMAGE}:${DOCKER_TAG}
                        docker push ${DOCKER_IMAGE}:latest
                    """
                }
            }
        }
        
        stage('Deploy to Kubernetes') {
            steps {
                sh """
                    kubectl set image deployment/k8s-jenkins-app \
                        app=${DOCKER_IMAGE}:${DOCKER_TAG} \
                        -n ${KUBERNETES_NAMESPACE} || \
                    kubectl apply -f k8s/deployment.production.yaml -n ${KUBERNETES_NAMESPACE}
                """
                sh "kubectl apply -f k8s/service.production.yaml -n ${KUBERNETES_NAMESPACE}"
                sh "kubectl rollout status deployment/k8s-jenkins-app -n ${KUBERNETES_NAMESPACE} --timeout=5m"
            }
        }
        
        stage('Smoke Tests') {
            steps {
                sh 'curl -f http://k8s-jenkins-app-service.${KUBERNETES_NAMESPACE}.svc.cluster.local/health || exit 1'
            }
        }
    }
    
    post {
        success {
            echo 'Deployment successful!'
        }
        failure {
            sh 'kubectl rollout undo deployment/k8s-jenkins-app -n ${KUBERNETES_NAMESPACE}'
            echo 'Deployment failed, rolled back!'
        }
        always {
            cleanWs()
        }
    }
}
```

---

### Phase 9: Backup & Disaster Recovery

#### 9.1 Jenkins Backup
```bash
# Create backup script
cat > /backup/jenkins-backup.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/backup/jenkins"
DATE=$(date +%Y%m%d_%H%M%S)
tar -czf $BACKUP_DIR/jenkins_home_$DATE.tar.gz /var/lib/jenkins
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete
EOF

chmod +x /backup/jenkins-backup.sh

# Schedule daily backup
echo "0 2 * * * /backup/jenkins-backup.sh" | crontab -
```

#### 9.2 Kubernetes Backup
```bash
# Install Velero (backup tool)
# For EKS: Use AWS Backup or Velero
```

---

### Phase 10: SSL/TLS Certificates

#### 10.1 Install cert-manager
```bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml
```

#### 10.2 Create ClusterIssuer
```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: your-email@example.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
```

---

## 📊 **PRODUCTION CHECKLIST**

### Infrastructure
- [ ] Kubernetes cluster created (EKS/k3s)
- [ ] Cluster nodes configured
- [ ] Network policies configured
- [ ] Load balancer configured
- [ ] DNS configured

### Security
- [ ] RBAC configured
- [ ] Network policies applied
- [ ] Secrets management setup
- [ ] SSL/TLS certificates configured
- [ ] Security scanning enabled
- [ ] Non-root containers
- [ ] Pod security policies

### Application
- [ ] Production deployment manifests
- [ ] Resource limits configured
- [ ] Health checks (liveness/readiness/startup)
- [ ] Pod anti-affinity configured
- [ ] Rolling update strategy
- [ ] ConfigMaps for configuration
- [ ] Secrets for sensitive data
- [ ] ServiceAccount configured

### CI/CD
- [ ] Jenkins plugins installed
- [ ] Docker registry configured
- [ ] Kubernetes credentials configured
- [ ] Production Jenkinsfile updated
- [ ] Pipeline tested end-to-end
- [ ] Rollback mechanism tested

### Monitoring
- [ ] Prometheus installed
- [ ] Grafana installed
- [ ] ServiceMonitor created
- [ ] Alerts configured
- [ ] Logging stack (ELK/Loki)
- [ ] Dashboard created

### High Availability
- [ ] Multiple replicas (3+)
- [ ] Pod anti-affinity
- [ ] HPA configured
- [ ] PDB configured
- [ ] Multi-zone deployment

### Backup & Recovery
- [ ] Jenkins backup configured
- [ ] Kubernetes backup configured
- [ ] Disaster recovery plan
- [ ] Backup tested

---

## 🎯 **IMMEDIATE NEXT STEPS**

### Priority 1: Complete Basic Setup
1. **Install Docker on Jenkins server** (if not done)
2. **Install kubectl on Jenkins server** (if not done)
3. **Set up Kubernetes cluster** (EKS or k3s)
4. **Configure Docker registry credentials in Jenkins**
5. **Test pipeline end-to-end**

### Priority 2: Production Hardening
1. **Update Kubernetes manifests** with production settings
2. **Add security configurations** (RBAC, Network Policies)
3. **Set up monitoring** (Prometheus/Grafana)
4. **Configure SSL/TLS** (cert-manager)

### Priority 3: Advanced Features
1. **Set up auto-scaling** (HPA)
2. **Configure ingress** (NGINX)
3. **Set up logging** (ELK/Loki)
4. **Implement backup strategy**

---

## 📝 **COMMANDS TO RUN NOW**

### On Jenkins Server:

```bash
# 1. Verify all tools
go version
docker --version
kubectl version --client
java -version

# 2. Set up Kubernetes cluster (choose one)

# Option A: k3s (quick setup)
curl -sfL https://get.k3s.io | sh -
sudo cat /etc/rancher/k3s/k3s.yaml > ~/.kube/config
sudo mkdir -p /var/lib/jenkins/.kube
sudo cp ~/.kube/config /var/lib/jenkins/.kube/config
sudo chown -R jenkins:jenkins /var/lib/jenkins/.kube

# Option B: EKS (production)
# Follow EKS setup steps above

# 3. Create production namespace
kubectl create namespace production

# 4. Test kubectl access
kubectl get nodes
kubectl get namespaces
```

---

## 💰 **ESTIMATED COSTS**

### AWS Resources (Monthly)
- EC2 Instance (t3.medium): ~$30-50
- EKS Control Plane: ~$73
- EKS Worker Nodes (2x t3.medium): ~$60-100
- ECR Storage: ~$1-5
- Load Balancer: ~$20-30
- **Total: ~$184-258/month**

### Alternative (k3s on single VM)
- EC2 Instance (t3.large): ~$60-80
- Docker Hub: Free (public) or $7/month (private)
- **Total: ~$60-87/month**

---

## 🚀 **QUICK START GUIDE**

### Step 1: Set up Kubernetes (5 minutes)
```bash
# Install k3s
curl -sfL https://get.k3s.io | sh -

# Configure kubectl
mkdir -p ~/.kube
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo chown $USER:$USER ~/.kube/config

# Test
kubectl get nodes
```

### Step 2: Configure Jenkins (10 minutes)
1. Install plugins (Pipeline, Docker, Kubernetes)
2. Add Docker Hub credentials
3. Configure kubectl access
4. Create pipeline job

### Step 3: Deploy Application (5 minutes)
```bash
# Create namespace
kubectl create namespace production

# Apply manifests
kubectl apply -f k8s/deployment.production.yaml -n production
kubectl apply -f k8s/service.production.yaml -n production

# Verify
kubectl get pods -n production
kubectl get svc -n production
```

### Step 4: Test Pipeline (5 minutes)
1. Push code to GitHub
2. Jenkins detects change
3. Pipeline runs automatically
4. Verify deployment

---

## 📚 **LEARNING RESOURCES**

### Kubernetes
- Official Docs: https://kubernetes.io/docs/
- kubectl Cheat Sheet: https://kubernetes.io/docs/reference/kubectl/cheatsheet/

### Jenkins
- Jenkins Pipeline Docs: https://www.jenkins.io/doc/book/pipeline/
- Jenkinsfile Examples: https://www.jenkins.io/doc/pipeline/examples/

### Production Best Practices
- Kubernetes Best Practices: https://kubernetes.io/docs/concepts/security/
- 12-Factor App: https://12factor.net/

---

## ✅ **SUCCESS CRITERIA**

Your setup is production-ready when:
- [x] Jenkins is running and accessible
- [ ] Pipeline builds Docker images successfully
- [ ] Images are pushed to registry
- [ ] Kubernetes cluster is running
- [ ] Application deploys to Kubernetes
- [ ] Zero-downtime deployments work
- [ ] Monitoring shows metrics
- [ ] Alerts are configured
- [ ] Backups are automated
- [ ] Security policies are enforced

---

## 🎓 **INTERVIEW TALKING POINTS**

When explaining your setup:

1. **Infrastructure**: "I set up Jenkins on AWS EC2 with all required tools (Go, Docker, kubectl) for CI/CD."

2. **CI/CD Pipeline**: "The pipeline automates build, test, security scan, Docker image creation, registry push, and Kubernetes deployment."

3. **High Availability**: "I configured 3 replicas with pod anti-affinity, rolling updates with maxUnavailable: 0 for zero downtime."

4. **Security**: "Implemented RBAC, network policies, non-root containers, and secrets management."

5. **Monitoring**: "Set up Prometheus for metrics, Grafana for visualization, and configured alerts."

6. **Best Practices**: "Used ConfigMaps for configuration, Secrets for sensitive data, HPA for auto-scaling, and PDB for availability."

---

**Next Action:** Complete Phase 1 (CI/CD Pipeline Setup) → Then move to Phase 2 (Kubernetes Cluster)
