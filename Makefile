# ⚙️ Terraform commands
terraform:
	terraform init
	terraform plan
	terraform apply -auto-approve

# 🚀 Build e Deploy da Lambda
deploy:
	cd cmd && GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
	cd cmd && zip -r ../function.zip bootstrap ../internal ../go.mod ../go.sum
	aws s3 cp function.zip s3://video-processamento-fiap/function.zip --region us-east-1
	aws lambda update-function-code --function-name video-processor --s3-bucket video-processamento-fiap --s3-key function.zip --region us-east-1

# 👇 Roda terraform e deploy em sequência
all: terraform deploy