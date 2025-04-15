# 🎥 Video Processor - FIAP Pós-Graduação

Serviço de processamento de vídeo construído em Go, usando AWS Lambda, S3, DynamoDB, SNS e Terraform como infraestrutura. Parte do projeto final da pós-graduação na FIAP.

---

## 📦 Funcionalidades

- Processa vídeos e gera imagens (prints) a cada 20 segundos usando FFmpeg.
- Faz upload dos frames gerados no S3.
- Persiste os dados no DynamoDB.
- Notifica serviços externos via SNS:
    - Em caso de sucesso.
    - Em caso de erro.

---

## 🧱 Arquitetura

- **Linguagem:** Go
- **Banco:** DynamoDB
- **Infraestrutura:** Terraform
- **Processamento:** AWS Lambda (com FFmpeg via Layer ou Docker)
- **Mensageria:** SNS
- **Armazenamento:** S3

---

## 🚀 Como rodar o projeto

### ✅ Pré-requisitos

- [AWS CLI configurado](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html)
- Go instalado (`go 1.20+`)
- Terraform instalado (`terraform 1.3+`)
- Bucket S3 de deploy já criado

### ⚙️ Deploy

1. **Clone o projeto e entre na pasta:**

```bash
git clone https://github.com/seu-usuario/video-processor.git
cd video-processor
