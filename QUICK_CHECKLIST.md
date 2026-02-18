# Quick Production Setup Checklist

## ✅ **COMPLETED**
- [x] AWS VM created
- [x] Jenkins installed (Java 25)
- [x] Go installed
- [x] Docker installed
- [x] kubectl installed
- [x] Code in GitHub

## 🔄 **DO NOW (Priority 1)**

### 1. Complete Jenkins Setup
```bash
# On Jenkins server
# Install missing plugins via Jenkins UI
# Configure Docker Hub credentials
# Test pipeline
```

### 2. Set Up Kubernetes Cluster
```bash
# Option A: k3s (quick)
curl -sfL https://get.k3s.io | sh -
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo mkdir -p /var/lib/jenkins/.kube
sudo cp ~/.kube/config /var/lib/jenkins/.kube/config
sudo chown -R jenkins:jenkins /var/lib/jenkins/.kube

# Option B: EKS (production)
eksctl create cluster --name production --region ap-southeast-2 --node-type t3.medium --nodes 2
aws eks update-kubeconfig --name production --region ap-southeast-2
```

### 3. Configure Docker Registry
- Create Docker Hub account (or use ECR)
- Add credentials to Jenkins
- Update Jenkinsfile with registry URL

### 4. Update Jenkinsfile
Add PATH and registry push stage (see PRODUCTION_SETUP_REPORT.md)

### 5. Test End-to-End
- Push code → Jenkins builds → Deploys to K8s

## 📋 **PRODUCTION FEATURES TO ADD**

### Security
- [ ] RBAC (ServiceAccount + Role)
- [ ] Network Policies
- [ ] Secrets (not hardcoded)
- [ ] Non-root containers
- [ ] SSL/TLS certificates

### High Availability
- [ ] 3+ replicas
- [ ] Pod anti-affinity
- [ ] HPA (auto-scaling)
- [ ] PDB (pod disruption budget)
- [ ] Rolling update strategy

### Monitoring
- [ ] Prometheus
- [ ] Grafana dashboards
- [ ] ServiceMonitor
- [ ] Alerts
- [ ] Logging (ELK/Loki)

### CI/CD
- [ ] Security scanning
- [ ] Registry push
- [ ] Environment separation (staging/prod)
- [ ] Rollback on failure
- [ ] Smoke tests

---

**See PRODUCTION_SETUP_REPORT.md for detailed steps**
