ZIP_PATH = build/video-processor.zip
LAMBDA_NAME = video-processor
S3_BUCKET = video-processor-postechfiap
REGION = us-east-1

build:
	GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/main.go
	mkdir -p build
	zip -r build/video-processor.zip bootstrap ffmpeg internal go.mod go.sum

terraform: build
	terraform init
	terraform plan
	terraform apply -auto-approve

deploy: build
	aws s3 cp $(ZIP_PATH) s3://$(S3_BUCKET)/function.zip --region $(REGION)
	aws lambda update-function-code \
		--function-name $(LAMBDA_NAME) \
		--s3-bucket $(S3_BUCKET) \
		--s3-key function.zip \
		--region $(REGION)

rebuild:
	rm -rf bootstrap build
	make build
	make deploy