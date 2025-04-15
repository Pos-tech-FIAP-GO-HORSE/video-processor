variable "aws_region" {
  default = "us-east-1"
}

variable "s3_bucket_name" {
  default = "video-processamento-fiap"
}

variable "dynamodb_table_name" {
  default = "video-processamento"
}

variable "sns_trigger_topic_arn" {
  description = "ARN do tópico SNS que dispara o processamento de vídeo"
}
