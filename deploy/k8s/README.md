# Totem de Pedidos - Kubernetes Deployment

## Arquivos K8S

| Arquivo | Descrição |
|---------|-----------|
| `app-namespace.yaml` | Namespace dedicado `totem-pedidos` |
| `app-configmap.yaml` | Configurações não-sensíveis da aplicação |
| `app-secret.yaml` | Credenciais e dados sensíveis (base64 encoded) |
| `app-deployment.yaml` | Deployment principal da aplicação Go |
| `app-service.yaml` | Service para a aplicação |
| `app-ingress.yaml` | Ingress para acesso externo |
| `app-hpa.yaml` | Horizontal Pod Autoscaler |
| `deploy.sh` | Script automatizado de deployment |
