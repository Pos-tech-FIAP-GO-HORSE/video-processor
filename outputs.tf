output "s3_bucket" {
  value = aws_s3_bucket.video_bucket.id
}

output "dynamodb_table" {
  value = aws_dynamodb_table.videos_table.id
}

output "lambda_function_name" {
  value = aws_lambda_function.video_processor.function_name
}
