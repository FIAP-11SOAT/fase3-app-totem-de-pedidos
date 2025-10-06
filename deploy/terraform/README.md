# Terraform EKS Deployment for Totem de Pedidos

This Terraform configuration deploys the Totem de Pedidos application to an existing EKS cluster using your existing Kubernetes YAML files.

## Prerequisites

1. **AWS CLI configured** with appropriate permissions
2. **Terraform >= 1.0** installed
3. **Docker** installed (for building images)
4. **kubectl** configured to access your EKS cluster
5. **Existing EKS cluster** up and running

## Required AWS Permissions

Your AWS credentials need the following permissions:
- ECR: Create repositories, push/pull images
- EKS: Describe clusters
- Secrets Manager: Read secrets (if using)
- IAM: Basic read permissions for EKS authentication

## Setup

**Simple 3-step deployment:**

1. **Navigate to terraform directory:**
   ```bash
   cd deploy/terraform
   ```

2. **Initialize Terraform:**
   ```bash
   terraform init
   ```

3. **Deploy:**
   ```bash
   terraform apply
   ```

That's it! Everything is pre-configured.

## What This Does

**Complete CI/CD Pipeline in Terraform:**

1. **Creates ECR Repository** - Stores your Docker images
2. **Builds & Pushes Docker Image** - Automatically builds from your Dockerfile and pushes to ECR
3. **Deploys Kubernetes Resources** - Uses your existing YAML files:
   - Namespace (`app-namespace.yaml`)
   - Secret (`app-secret.yaml`) 
   - Deployment (`app-deployment.yaml`) - **Automatically updated to use ECR image**
   - Service (`app-service.yaml`)
   - HPA (`app-hpa.yaml`)
   - Ingress (`app-ingress.yaml`)

## Configuration

All configuration is hard-coded in `variables.tf`:
- **AWS Region:** us-east-1
- **Project Name:** fase3-totem-de-pedidos  
- **EKS Cluster:** fase3-totem-de-pedidos-eks-cluster

No variables file needed!

## Complete Deployment Flow

When you run `terraform apply`, this happens:

1. ✅ **ECR Repository Created**
2. ✅ **Docker Image Built** (from your Dockerfile)  
3. ✅ **Image Pushed to ECR**
4. ✅ **Kubernetes Resources Deployed** (using ECR image)
5. ✅ **Application Running** on your EKS cluster

## Troubleshooting

### Docker Build Issues
```bash
# Manual build test
cd ../..  # Go to project root
docker build -t test-build .
```

### EKS Access Issues
```bash
# Update kubeconfig
aws eks update-kubeconfig --region us-east-1 --name fase3-totem-de-pedidos-eks-cluster

# Test access
kubectl get nodes
```

### Check Deployment
```bash
kubectl get pods -n totem-pedidos
kubectl get svc -n totem-pedidos
```

## Cleanup

To destroy all resources:
```bash
terraform destroy
```

**Note:** This will delete the ECR repository and all images. Your EKS cluster and other AWS resources will remain untouched.

## File Structure

```
deploy/terraform/
├── _providers.tf          # AWS, Kubernetes, and GitHub providers
├── variables.tf           # Variables and data sources
├── ecr.tf                # ECR repository and Docker build
├── outputs.tf            # Output values
├── terraform.tfvars.example  # Example configuration
└── README.md             # This file
```