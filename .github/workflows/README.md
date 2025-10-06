# GitHub Actions Setup for Terraform Deployment

This document explains how to configure GitHub Actions to automatically deploy your Totem de Pedidos application to EKS using Terraform.

## Required GitHub Secrets

You need to configure the following secrets in your GitHub repository:

### AWS Credentials

1. Go to your GitHub repository
2. Navigate to **Settings** → **Secrets and variables** → **Actions**
3. Add the following secrets:

| Secret Name | Description | Example |
|------------|-------------|---------|
| `AWS_ACCESS_KEY_ID` | AWS Access Key ID with ECR and EKS permissions | `AKIAIOSFODNN7EXAMPLE` |
| `AWS_SECRET_ACCESS_KEY` | AWS Secret Access Key | `wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY` |

### Required AWS Permissions

Your AWS user/role needs the following permissions:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ecr:*",
                "eks:DescribeCluster",
                "eks:DescribeNodegroup",
                "eks:ListClusters",
                "sts:GetCallerIdentity",
                "iam:PassRole"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:GetObject",
                "s3:PutObject",
                "s3:DeleteObject",
                "s3:ListBucket"
            ],
            "Resource": [
                "arn:aws:s3:::your-terraform-state-bucket",
                "arn:aws:s3:::your-terraform-state-bucket/*"
            ]
        }
    ]
}
```

## Workflow Triggers

The GitHub Actions workflow will trigger on:

### Automatic Triggers

1. **Push to main branch** with changes to:
   - `deploy/terraform/**` (Terraform files)
   - `cmd/**`, `internal/**`, `pkg/**` (Go source code)
   - `Dockerfile`, `go.mod`, `go.sum` (Build configuration)

2. **Pull Requests** to main branch with changes to:
   - `deploy/terraform/**` (Shows plan in PR comments)

### Manual Triggers

3. **Manual dispatch** with options:
   - `plan` - Run terraform plan only
   - `apply` - Run terraform apply (deploy)
   - `destroy` - Destroy all resources ⚠️

## Workflow Jobs

### 1. terraform-validate
- Runs on all triggers
- Validates Terraform syntax and configuration
- Checks formatting with `terraform fmt`

### 2. terraform-plan  
- Runs after validation
- Creates Terraform execution plan
- Comments plan on Pull Requests
- Uploads plan as artifact

### 3. terraform-apply
- Runs only on main branch pushes or manual apply
- Downloads and applies the Terraform plan
- Builds Docker image and pushes to ECR
- Deploys to EKS cluster
- Verifies deployment status
- Provides deployment summary

### 4. terraform-destroy
- Runs only on manual trigger with 'destroy' option
- ⚠️ **WARNING:** Destroys all resources including ECR images

## Usage Examples

### Deploying Changes

1. **Make changes** to your code
2. **Commit and push** to main branch:
   ```bash
   git add .
   git commit -m "Update application code"
   git push origin main
   ```
3. **GitHub Actions will automatically:**
   - Validate Terraform
   - Build Docker image
   - Push to ECR
   - Deploy to EKS

### Manual Deployment

1. Go to **Actions** tab in your GitHub repository
2. Select **Deploy Totem de Pedidos to EKS** workflow
3. Click **Run workflow**
4. Choose action: `plan`, `apply`, or `destroy`
5. Click **Run workflow**

### Reviewing Changes (Pull Requests)

1. **Create a branch** and make changes
2. **Open a Pull Request**
3. **GitHub Actions will:**
   - Validate Terraform
   - Show plan in PR comments
   - Allow review before merging

## Monitoring Deployment

After deployment, you can monitor your application:

```bash
# Get pods status
kubectl get pods -n totem-pedidos

# Get services
kubectl get svc -n totem-pedidos

# Check application logs
kubectl logs -n totem-pedidos deployment/totem-pedidos-app

# Port forward to test locally
kubectl port-forward -n totem-pedidos svc/totem-pedidos-service 8080:80
```

## Troubleshooting

### Common Issues

1. **AWS Permissions Error**
   - Check AWS credentials in GitHub Secrets
   - Verify IAM permissions

2. **EKS Cluster Access**
   - Ensure cluster name is correct: `fase3-totem-de-pedidos-eks-cluster`
   - Check cluster exists in `us-east-1` region

3. **Docker Build Failure**
   - Check Dockerfile syntax
   - Verify all dependencies in `go.mod`

4. **Kubernetes Deployment Issues**
   - Check YAML files in `deploy/k8s/`
   - Verify secrets and configmaps

### Manual Debugging

If the workflow fails, you can run Terraform locally:

```bash
cd deploy/terraform
terraform init
terraform plan
terraform apply
```

## Security Notes

- ✅ AWS credentials are stored as GitHub Secrets (encrypted)
- ✅ Terraform state should use remote backend (S3 + DynamoDB)
- ✅ ECR images are scanned for vulnerabilities
- ⚠️ Be careful with `destroy` action - it deletes everything!

## Next Steps

After successful deployment:

1. **Configure monitoring** (CloudWatch, Prometheus)
2. **Set up alerts** for application health
3. **Configure domain/ingress** for external access
4. **Set up log aggregation** (ELK stack, CloudWatch Logs)