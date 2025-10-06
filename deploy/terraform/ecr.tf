resource "aws_ecr_repository" "app" {
  name                 = "${local.project_name}-app-ecr"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  encryption_configuration {
    encryption_type = "AES256"
  }

  tags = {
    Name = "${local.project_name}-ecr"
  }
}

resource "aws_ecr_lifecycle_policy" "app_repository_policy" {
  repository = aws_ecr_repository.app.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Keep last 10 images"
        selection = {
          tagStatus     = "tagged"
          tagPrefixList = ["v"]
          countType     = "imageCountMoreThan"
          countNumber   = 10
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}

resource "null_resource" "docker_build_and_push" {
  triggers = {
    dockerfile_hash = filemd5("${path.module}/../../Dockerfile")
    go_mod_hash = filemd5("${path.module}/../../go.mod")
    timestamp   = timestamp()
  }

  provisioner "local-exec" {
    command = <<-EOF
      # Navigate to project root
      cd ${path.module}/../..
      
      # Login to ECR
      aws ecr get-login-password --region ${local.aws_region} | docker login --username AWS --password-stdin ${aws_ecr_repository.app.repository_url}
      
      # Build the Docker image
      docker build -t ${local.project_name}:latest -f ./Dockerfile .
      
      # Tag for ECR
      docker tag ${local.project_name}:latest ${aws_ecr_repository.app.repository_url}:latest
      
      # Push to ECR
      docker push ${aws_ecr_repository.app.repository_url}:latest
      
      echo "✅ Docker image pushed to: ${aws_ecr_repository.app.repository_url}:latest"
    EOF
  }

  depends_on = [aws_ecr_repository.app]
}

resource "kubernetes_manifest" "app_namespace" {
  manifest = yamldecode(file("${path.module}/../k8s/app-namespace.yaml"))
}

resource "kubernetes_manifest" "app_secret" {
  manifest = yamldecode(file("${path.module}/../k8s/app-secret.yaml"))
}

locals {
  deployment_manifest_raw = yamldecode(file("${path.module}/../k8s/app-deployment.yaml"))
  deployment_manifest = merge(local.deployment_manifest_raw, {
    spec = merge(local.deployment_manifest_raw.spec, {
      template = merge(local.deployment_manifest_raw.spec.template, {
        spec = merge(local.deployment_manifest_raw.spec.template.spec, {
          containers = [
            merge(local.deployment_manifest_raw.spec.template.spec.containers[0], {
              image = "${aws_ecr_repository.app.repository_url}:latest"
            })
          ]
        })
      })
    })
  })
}

resource "kubernetes_manifest" "app_deployment" {
  manifest = local.deployment_manifest
  
  depends_on = [
    kubernetes_manifest.app_namespace,
    kubernetes_manifest.app_secret,
    null_resource.docker_build_and_push
  ]
}

resource "kubernetes_manifest" "app_service" {
  manifest = yamldecode(file("${path.module}/../k8s/app-service.yaml"))
  
  depends_on = [kubernetes_manifest.app_deployment]
}

resource "kubernetes_manifest" "app_hpa" {
  manifest = yamldecode(file("${path.module}/../k8s/app-hpa.yaml"))
  
  depends_on = [kubernetes_manifest.app_deployment]
}

resource "kubernetes_manifest" "app_ingress" {
  manifest = yamldecode(file("${path.module}/../k8s/app-ingress.yaml"))
  
  depends_on = [kubernetes_manifest.app_service]
}
