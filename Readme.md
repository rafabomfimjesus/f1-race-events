# 🏎️ F1 Race Events

API REST desenvolvida em Go para gerenciamento de eventos de corrida de Fórmula 1.

O projeto demonstra a construção de uma aplicação cloud-native utilizando:

- Go
- DynamoDB
- Docker
- Amazon ECS Fargate
- Amazon ECR
- Application Load Balancer (ALB)
- Terraform

Toda a infraestrutura é provisionada como código utilizando Terraform e a aplicação é executada em containers na AWS.

## Visão Geral

Este serviço expõe uma API HTTP simples com a rota `/race-events`:
- `POST /race-events` para criar um novo evento de corrida.
- `GET /race-events` para listar eventos com filtros opcionais.

A solução segue uma arquitetura organizada em camadas para separar responsabilidades.

## Arquitetura

```text
Client
   │
   ▼
Application Load Balancer
   │
   ▼
Amazon ECS Fargate
   │
   ▼
Container Go
   │
   ▼
Amazon DynamoDB
```

### Camadas principais

- `main.go` - inicializa o logger, cliente DynamoDB, repositório, casos de uso e registra as rotas HTTP.
- `routes/` - define os handlers HTTP e transforma requisições/respostas JSON.
- `usecases/` - implementa as regras de negócio:
  - criar evento de corrida
  - consultar eventos filtrando por `raceId`, `driverId`, `scuderiaId`, `eventType`, `lap`, `createdAtFrom` e `createdAtTo`
- `repositories/` - abstrai o acesso ao DynamoDB.
- `dtos/` - contém os objetos de transferência de dados usados na entrada e saída da API.
- `models/` - representa o domínio de negócio, incluindo o tipo `RaceEvent` e os tipos de evento válidos.
- `utils/` - contém utilitários como o cliente DynamoDB e logger.

### Fluxo de dados

1. O servidor HTTP recebe a requisição em `routes/`.
2. Os dados são validados e convertidos para DTOs.
3. Os casos de uso (`usecases/`) aplicam regras e constroem o modelo de domínio.
4. O repositório (`repositories/`) persiste ou consulta o DynamoDB.
5. A resposta é serializada como JSON e retornada ao cliente.

### DynamoDB

A persistência usa DynamoDB, com a tabela `race_event`.
Os eventos são armazenados com atributos como:
- `eventId`
- `raceId`
- `driverId`
- `scuderiaId`
- `lap`
- `eventType`
- `description`
- `createdAt`

### Infraestrutura com Terraform

A pasta `infra/` contém a infraestrutura como código para AWS:
- `main.tf` - configura o provider AWS
- `alb.tf` - Application Load Balancer
- `ecs.tf` - ECS cluster, tarefas e serviços
- `ecr.tf` - repositório ECR
- `iam.tf` - permissões e roles
- `security_groups.tf` - regras de segurança
- `data.tf` - dados remotos e recursos auxiliares
- `output.tf` - saídas do Terraform
- `variables.tf` - variáveis de configuração

## Passo a passo manual do Terraform

Antes de executar, verifique se você está com as credenciais AWS configuradas no ambiente.

No diretório `infra/`, execute:

```bash
terraform init
terraform validate
terraform apply
```

### Explicação dos comandos

- `terraform init` - inicializa o diretório de trabalho e baixa os providers necessários.
- `terraform validate` - verifica se a configuração do Terraform está sintaticamente correta.
- `terraform apply` - aplica as mudanças e provisiona a infraestrutura na AWS.

## Como executar localmente

1. Configure as credenciais AWS localmente.
2. Compile o projeto:

```bash
go build ./...
```

3. Execute o binário gerado.
4. A API estará disponível em `http://localhost:8080`.

## Observações

- O cliente DynamoDB usa a região `us-east-2` no código Go (`utils/dynamodb.go`).
- A API aceita e retorna JSON.
- O endpoint principal é `/race-events`.

## Desafios Técnicos

Durante o desenvolvimento foram implementados:

- Containerização da aplicação com Docker
- Provisionamento de infraestrutura com Terraform
- Deploy de containers no ECS Fargate
- Configuração de Load Balancer para exposição da API
- Integração com DynamoDB
- Gerenciamento de permissões IAM
- Health checks e monitoramento de containers
