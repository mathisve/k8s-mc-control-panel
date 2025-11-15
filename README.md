# k8s-mc-control-panel
A control panel for a Minecraft server running on K8s

**⚠️ DISCLAIMER: This project is provided as-is with no support. You will need to manually configure several files before deployment.**

## Overview

This project consists of:
- **frontend-control**: A React frontend that polls the backend server
- **backend**: A Golang backend that talks to the Kubernetes API server and responds to HTTP requests from the frontend

## Required Manual Configuration

Before deploying this project, you **must** modify the following files and configurations:

### 1. Backend Configuration (`backend/main.go`)

- **Line 23**: Update `NAMESPACE` constant to match your Kubernetes namespace (currently `"mc"`)
- **Line 24**: Update `DEPLOYMENT_NAME` constant to match your Minecraft deployment name (currently `"mc"`)
- **Line 25**: Update `IMAGE_NAME` constant to match your Minecraft server container image name (currently `"minecraft-server"`)
- **Line 31**: Update `authorizationToken` to your own secure token (currently `"Bearer b2xsaWUxMjMK"`)

### 2. Frontend Configuration (`frontend-control/src/App.js`)

- **Line 15**: Update `authorizationToken` to match the backend token (currently `"b2xsaWUxMjMK"`)
- **Line 17**: Update `url` to your backend API endpoint (currently `"https://mc-control-panel.homek8s.com/api/"`)

### 3. Kubernetes Deployment Manifests

#### Backend Deployment (`manifests/backend-deploy.yaml`)
- **Line 27**: Update the ECR image URL to your own container registry:
  - Replace `043039367084.dkr.ecr.us-east-1.amazonaws.com` with your AWS account ID and region
  - Or replace with your own container registry URL
- **Line 25**: Ensure `ecr-registry-secret` exists in your cluster, or remove `imagePullSecrets` if using public images
- **Line 7**: Update namespace if not using `mc` namespace

#### Frontend Deployment (`manifests/frontend-deploy.yaml`)
- **Line 26**: Update the ECR image URL to your own container registry:
  - Replace `043039367084.dkr.ecr.us-east-1.amazonaws.com` with your AWS account ID and region
  - Or replace with your own container registry URL
- **Line 24**: Ensure `ecr-registry-secret` exists in your cluster, or remove `imagePullSecrets` if using public images
- **Line 7**: Update namespace if not using `mc` namespace

#### Ingress Resources
- **`manifests/backend-ing.yaml` (Line 14)**: Update `host` field to your domain name (currently `mc-control-panel.homek8s.com`)
- **`manifests/frontend-ing.yaml` (Line 15)**: Update `host` field to your domain name (currently `mc-control-panel.homek8s.com`)

#### All Manifests
- Update all `namespace: mc` references throughout all YAML files if you're using a different namespace

### 4. OAuth2 Proxy Configuration (`manifests/proxy/oauth2-proxy.cfg`)

- **Line 2**: Set `cookie_secret` - generate a secure random string (32 bytes base64 encoded)
- **Line 7**: Set `client_secret` from your OAuth provider (GitHub, Google, etc.)
- **Line 8**: Set `client_id` from your OAuth provider
- **Line 9**: Update `redirect_url` to match your domain (currently `https://mc-control-panel.homek8s.com/oauth2/callback`)
- **Line 12**: Update `github_org` to your GitHub organization (currently `homek8s-com`), or change `provider` if using a different OAuth provider

### 5. Build and Deploy Script (`build-and-deploy.sh`)

- **Line 6**: Update `AWS_REGION` to your AWS region (currently `us-east-1`)
- **Line 8-9**: Update repository names if using different ECR repository names
- **Line 10-11**: Verify paths are correct for your project structure

### 6. Kubernetes Secrets

You will need to create the following secrets in your cluster:

1. **ECR Registry Secret** (if using private ECR images):
   ```bash
   kubectl create secret docker-registry ecr-registry-secret \
     --docker-server=<your-ecr-url> \
     --docker-username=AWS \
     --docker-password=$(aws ecr get-login-password --region <region>) \
     -n <namespace>
   ```

2. **OAuth2 Proxy Config Secret**:
   ```bash
   kubectl create secret generic frontend-proxy-config \
     --from-file=oauth2-proxy.cfg=manifests/proxy/oauth2-proxy.cfg \
     -n <namespace>
   ```
   Or use the provided script: `manifests/proxy/create-proxy.sh` (update namespace if needed)

### 7. RBAC Configuration

- **`manifests/role.yaml`**: Review and adjust permissions as needed for your security requirements
- **`manifests/sa.yaml`**: Ensure the service account name matches what's used in the backend deployment
- **`manifests/rolebinding.yaml`**: Verify namespace and service account references

## Deployment Steps

1. Build and push Docker images to your container registry
2. Update all configuration files listed above
3. Create required Kubernetes secrets
4. Apply all manifests in the `manifests/` directory
5. Ensure your ingress controller is properly configured
6. Set up DNS records pointing to your ingress

## Current TODOs
- Remove static ips, replace with something better
- Templatize deployment with Helm
- Add RCON proxy container to talk to the minecraft server directly
- Add user list in UI