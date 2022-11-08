// -- add fn --
resource "aws_lambda_function" "share_add" {
  runtime          = "go1.x"
  function_name    = var.name
  handler          = var.binary
  filename         = var.archive
  source_code_hash = filebase64sha256(var.archive)
  role             = aws_iam_role.share_add.arn

  environment {
    variables = {
      LOCAL = var.local ? "1" : null
      SHARE_COUNT_NAME = var.share_count_name
      SHARE_FILES_NAME = var.share_files_name
      SHARE_FILES_HOST = var.share_files_host
    }
  }
}

// -- add fn public url
resource "aws_lambda_function_url" "share_add" {
  function_name      = aws_lambda_function.share_add.arn
  authorization_type = "NONE"

  cors {
    allow_credentials = true
    allow_origins     = ["*"]
    allow_methods     = ["*"]
    allow_headers     = ["date", "keep-alive"]
    expose_headers    = ["keep-alive", "date"]
    max_age           = 86400
  }
}

// -- add fn role
resource "aws_iam_role" "share_add" {
  name = "${var.name}-iam"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}
