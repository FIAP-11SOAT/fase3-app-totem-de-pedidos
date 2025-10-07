# 🍔 SOAT Tech Challenge - Totem de Pedido (Fase 1)

Este repositório contém a implementação do sistema de autoatendimento para uma lanchonete em expansão. O projeto é o desafio prático da **Fase 1 do curso de Pós-graduação em Arquitetura de Software da FIAP**, integrando os conhecimentos adquiridos em todas as disciplinas do módulo.

## 🚀 Objetivo

Desenvolver um sistema de backend monolítico que permita:

- Clientes realizarem pedidos personalizados de forma autônoma.
- Pagamento dos pedidos via QR Code (Mercado Pago).
- Acompanhamento em tempo real do status dos pedidos.
- Administração de clientes, produtos e pedidos.

### 📺 Vídeo da Apresentação

[Fase 1 - Totem de Pedido (DDD + Hexagonal + Docker)](https://github.com/FIAP-11SOAT/totem-de-pedidos)

## 📋 Documentação DDD

[Miro DDD Documentação](https://miro.com/app/board/uXjVIHWL0sE=/?share_link_id=24901001533)

## 💻 Executar serviço

```
docker compose up --build
```

## 📖 Documentação da API

O projeto disponibiliza três visualizações para a documentação OpenAPI:

- **Swagger UI**:  
  [http://localhost:8080/docs/swagger](http://localhost:8080/docs/swagger)

- **Redoc**:  
  [http://localhost:8080/docs/redoc](http://localhost:8080/docs/redoc)

- **Scalar**:  
  [http://localhost:8080/docs/scalar](http://localhost:8080/docs/scalar)

O arquivo OpenAPI bruto pode ser acessado em:  
[http://localhost:8080/docs/openapi.yaml](http://localhost:8080/docs/openapi.yaml)

## 📋 Funcionalidades

### Cliente
- Cadastro e identificação via CPF
- Montagem de pedido com:
  - Lanche
  - Acompanhamento
  - Bebida
  - Sobremesa
- Pagamento via QR Code
- Acompanhamento de status:
  - Recebido
  - Em preparação
  - Pronto
  - Finalizado

### Administração
- Gerenciamento de clientes
- Gerenciamento de produtos e categorias
- Acompanhamento de pedidos em tempo real

## 🧠 Aprendizados Aplicados

### Domain-Driven Design (DDD)
- Introdução e fundamentos de DDD
- Domain Storytelling e descoberta de conhecimento
- Contextos delimitados (Bounded Contexts)
- Event Storming dos fluxos:
  - Pedido e pagamento
  - Preparação e entrega
- Refinamento técnico com Definition of Ready e Done

### Arquitetura de Software
- Arquitetura Hexagonal
- Modularização, testabilidade e escalabilidade
- Documentação de decisões arquiteturais

### Qualidade de Software
- Testes de unidade, integração, carga
- TDD e BDD aplicados
- Ferramentas de cobertura e relatórios de testes

### Dockerização
- Dockerfile e docker-compose configurados
- Melhores práticas de containerização
- Segurança de containers e uso de ECS

### Desenvolvimento Seguro
- OWASP TOP 10
- Análise estática de código
- Proteções contra ataques comuns como XSS, SQL Injection e Buffer Overflow

## 🧪 Tecnologias Utilizadas

- Docker / Docker Compose
- Golang
- Swagger/OpenAPI
- Banco de Dados: 'Postgres'
- Ferramentas de testes: 'testing + testify'
- Mercado Pago (Integracao com Mercado Pago)
