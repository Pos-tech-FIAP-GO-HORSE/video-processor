# Video Processor

Serviço de processamento de vídeo construído em Go, usando AWS Lambda, S3, DynamoDB, SNS e Terraform como infraestrutura.

---

## Funcionalidades

- Processa vídeos e gera imagens (prints) a cada 20 segundos usando FFmpeg
- Faz upload dos frames gerados no S3
- Persiste os dados no DynamoDB
- Notifica serviços externos via SNS:
    - Em caso de erro

---

## Tecnologias

- Go 1.20+
- AWS Lambda
- Amazon S3
- Amazon DynamoDB
- Amazon SNS
- Terraform 1.3+
- FFmpeg

---

## Instalação

```bash
git clone https://github.com/seu-usuario/video-processor.git
cd video-processor
```

---

## Executando Localmente

### Pré-requisitos

- [AWS CLI configurado](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html)
- Go instalado (`go 1.20+`)
- Terraform instalado (`terraform 1.3+`)
- Bucket S3 de deploy já criado

## Scripts Disponíveis

- `go build` - Compila o projeto Go
- `go test` - Executa os testes
- `terraform apply` - Faz deploy da infraestrutura
- `terraform destroy` - Remove a infraestrutura

---

## Testes

```bash
go test
```

---

## Estrutura de pastas

```bash
.
├── cmd/                   # Ponto de entrada da aplicação
├── internal/              # Código interno do projeto
│   ├── handlers/          # Handlers da AWS Lambda
│   ├── services/          # Serviços utilizados pelo projeto
│   └── repositories/      # Repositórios de dados
├── terraform/             # Configurações de infraestrutura
└── scripts/               # Scripts auxiliares
```

---

## Membros

| Membro                        | Info     |
| ----------------------------- | -------- |
| Caio Martins Pereira          | RM357712 |
| Rafael de Souza Ribeiro       | RM357622 |
| Thaís Oliveira de Moura       | RM357737 |
| Victor Toschi                 | RM356847 |
