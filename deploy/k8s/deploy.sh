#!/bin/bash

set -e

echo "Starting Totem de Pedidos deployment..."

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    echo "kubectl is not installed or not in PATH"
    exit 1
fi

# Check if we can connect to Kubernetes cluster
if ! kubectl cluster-info &> /dev/null; then
    echo "Cannot connect to Kubernetes cluster"
    exit 1
fi

echo "Kubernetes cluster connection verified"

# Apply manifests in the correct order
echo "Creating namespace..."
kubectl apply -f app-namespace.yaml

echo "Creating ConfigMaps..."
kubectl apply -f app-configmap.yaml

echo "Creating Secrets..."
kubectl apply -f app-secret.yaml

echo "Deploying application..."
kubectl apply -f app-deployment.yaml
kubectl apply -f app-service.yaml

echo "Waiting for application to be ready..."
kubectl wait --for=condition=Available deployment/totem-pedidos-app -n totem-pedidos --timeout=300s

echo "Creating Ingress..."
kubectl apply -f app-ingress.yaml

echo "Creating HPA..."
kubectl apply -f app-hpa.yaml
echo "   2. Update the image tag in app-deployment.yaml if needed"
echo "   3. Configure DNS or /etc/hosts for 'totem-pedidos.local'"
echo "   4. Update MP_NOTIFICATION_URL in secret.yaml with your actual domain"
