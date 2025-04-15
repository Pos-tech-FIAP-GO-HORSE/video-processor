provider "aws" {
  region = var.aws_region
}

resource "aws_s3_bucket" "video_bucket" {
  bucket = var.s3_bucket_name

  force_destroy = true # cuidado em prod
}

resource "aws_dynamodb_table" "videos_table" {
  name         = var.dynamodb_table_name
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "PK"
  range_key    = "SK"

  attribute {
    name = "PK"
    type = "S"
  }

  attribute {
    name = "SK"
    type = "S"
  }
}

resource "aws_sns_topic" "success_topic" {
  name = "video-process-success"
}

resource "aws_sns_topic" "error_topic" {
  name = "video-process-error"
}

resource "aws_iam_role" "lambda_exec" {
  name = "video_processor_lambda_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect = "Allow",
      Principal = { Service = "lambda.amazonaws.com" },
      Action = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_policy" "lambda_policy" {
  name = "video_processor_access_policy"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect   = "Allow",
        Action   = ["s3:*"],
        Resource = [
          aws_s3_bucket.video_bucket.arn,
          "${aws_s3_bucket.video_bucket.arn}/*"
        ]
      },
      {
        Effect   = "Allow",
        Action   = ["dynamodb:*"],
        Resource = [aws_dynamodb_table.videos_table.arn]
      },
      {
        Effect   = "Allow",
        Action   = ["sns:Publish"],
        Resource = [
          aws_sns_topic.success_topic.arn,
          aws_sns_topic.error_topic.arn
        ]
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_custom_policy" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = aws_iam_policy.lambda_policy.arn
}

resource "aws_lambda_function" "video_processor" {
  function_name = "video-processor"

  role    = aws_iam_role.lambda_exec.arn
  handler = "main"
  runtime = "go1.x"
  timeout = 300

  filename         = "build/video-processor.zip"
  source_code_hash = filebase64sha256("build/video-processor.zip")

  environment {
    variables = {
      S3_BUCKET                       = aws_s3_bucket.video_bucket.id
      DYNAMODB_TABLE                  = aws_dynamodb_table.videos_table.id
      PROCESSAMENTO_SUCESSO_TOPIC_ARN = aws_sns_topic.success_topic.arn
      PROCESSAMENTO_ERRO_TOPIC_ARN    = aws_sns_topic.error_topic.arn
    }
  }
}

resource "aws_sns_topic_subscription" "lambda_trigger" {
  topic_arn = var.sns_trigger_topic_arn
  protocol  = "lambda"
  endpoint  = aws_lambda_function.video_processor.arn
}

resource "aws_lambda_permission" "allow_sns" {
  statement_id  = "AllowExecutionFromSNS"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.video_processor.function_name
  principal     = "sns.amazonaws.com"
  source_arn    = var.sns_trigger_topic_arn
}
