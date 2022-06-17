terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.9"
    }
  }

  required_version = ">= 0.14.9"
}

// localstack
provider "aws" {
  region                      = var.share_region
  s3_use_path_style           = true
  access_key                  = "test"
  secret_key                  = "test"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  // configure all services to use localstack url
  endpoints {
    apigateway   = "http://localhost:4566"
    apigatewayv2 = "http://localhost:4566"
    dynamodb     = "http://localhost:4566"
    iam          = "http://localhost:4566"
    lambda       = "http://localhost:4566"
    s3           = "http://localhost:4566"
  }
}

// s3: share url bucket
resource "aws_s3_bucket" "share_files" {
  bucket = var.share_files_name
}

resource "aws_s3_bucket_acl" "share_files" {
  bucket = aws_s3_bucket.share_files.id
  acl    = "public-read"
}

resource "aws_s3_bucket_website_configuration" "share_files" {
  bucket = aws_s3_bucket.share_files.bucket

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "index.html"
  }
}

resource "aws_s3_bucket_policy" "share_files" {
  bucket = aws_s3_bucket.share_files.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "PublicReadGetObject"
        Effect    = "Allow"
        Principal = "*"
        Action    = "s3:GetObject"
        Resource = [
          "${aws_s3_bucket.share_files.arn}",
          "${aws_s3_bucket.share_files.arn}/*",
        ]
      },
    ]
  })
}

// dynamo: share kv store
resource "aws_dynamodb_table" "share_count" {
  name           = var.share_count_name
  billing_mode   = "PROVISIONED"
  read_capacity  = 25
  write_capacity = 25
  hash_key       = "Id"

  attribute {
    name = "Id"
    type = "S"
  }
}

// lambda: share service
resource "aws_lambda_function" "share_add" {
  runtime          = "go1.x"
  function_name    = var.share_add_name
  handler          = var.share_add_binary
  filename         = var.share_add_archive
  source_code_hash = filebase64sha256(var.share_add_archive)
  role             = aws_iam_role.share_add.arn
}

# resource "aws_lambda_function_url" "share_add_url" {
#   function_name      = aws_lambda_function.share_add.arn
#   authorization_type = "NONE"

#   cors {
#     allow_credentials = true
#     allow_origins     = ["*"]
#     allow_methods     = ["*"]
#     allow_headers     = ["date", "keep-alive"]
#     expose_headers    = ["keep-alive", "date"]
#     max_age           = 86400
#   }
# }

resource "aws_iam_role" "share_add" {
  name = var.share_add_iam

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

// api gateway: just for now
# resource "aws_apigatewayv2_api" "share_add" {
#   name          = "${aws_lambda_function.share_add.function_name}-api"
#   protocol_type = "HTTP"
# }

resource "aws_api_gateway_rest_api" "share_add" {
  name = "${aws_lambda_function.share_add.function_name}-api"
}

# resource "aws_apigatewayv2_stage" "share_add" {
#   api_id = aws_apigatewayv2_api.share_add.id

#   name        = "${aws_lambda_function.share_add.function_name}-stage"
#   auto_deploy = true

#   access_log_settings {
#     destination_arn = aws_cloudwatch_log_group.share_add.arn

#     format = jsonencode({
#       requestId               = "$context.requestId"
#       sourceIp                = "$context.identity.sourceIp"
#       requestTime             = "$context.requestTime"
#       protocol                = "$context.protocol"
#       httpMethod              = "$context.httpMethod"
#       resourcePath            = "$context.resourcePath"
#       routeKey                = "$context.routeKey"
#       status                  = "$context.status"
#       responseLength          = "$context.responseLength"
#       integrationErrorMessage = "$context.integrationErrorMessage"
#     })
#   }
# }

resource "aws_api_gateway_resource" "share_add_proxy" {
  rest_api_id = "${aws_api_gateway_rest_api.share_add.id}"
  parent_id   = "${aws_api_gateway_rest_api.share_add.root_resource_id}"
  path_part   = "{proxy+}"
}

resource "aws_api_gateway_method" "share_add_proxy" {
  rest_api_id   = "${aws_api_gateway_rest_api.share_add.id}"
  resource_id   = "${aws_api_gateway_resource.share_add_proxy.id}"
  http_method   = "ANY"
  authorization = "NONE"
}

# resource "aws_apigatewayv2_integration" "share_add" {
#   api_id = aws_apigatewayv2_api.share_add.id

#   integration_uri    = aws_lambda_function.share_add.invoke_arn
#   integration_type   = "AWS_PROXY"
#   integration_method = "POST"
# }

resource "aws_api_gateway_integration" "share_add" {
  rest_api_id = "${aws_api_gateway_rest_api.share_add.id}"
  resource_id = "${aws_api_gateway_method.share_add_proxy.resource_id}"
  http_method = "${aws_api_gateway_method.share_add_proxy.http_method}"

  uri                     = "${aws_lambda_function.share_add.invoke_arn}"
  type                    = "AWS_PROXY"
  integration_http_method = "POST"
}

# resource "aws_apigatewayv2_route" "share_add" {
#   api_id = aws_apigatewayv2_api.share_add.id

#   route_key = "GET /hello"
#   target    = "integrations/${aws_apigatewayv2_integration.share_add.id}"
# }

# resource "aws_cloudwatch_log_group" "share_add" {
#   name = "/aws/share_add/${aws_apigatewayv2_api.share_add.name}"

#   retention_in_days = 30
# }

resource "aws_api_gateway_method" "share_add_proxy_root" {
  rest_api_id   = "${aws_api_gateway_rest_api.share_add.id}"
  resource_id   = "${aws_api_gateway_rest_api.share_add.root_resource_id}"
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "share_add_root" {
  rest_api_id = "${aws_api_gateway_rest_api.share_add.id}"
  resource_id = "${aws_api_gateway_method.share_add_proxy_root.resource_id}"
  http_method = "${aws_api_gateway_method.share_add_proxy_root.http_method}"

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "${aws_lambda_function.share_add.invoke_arn}"
}

resource "aws_api_gateway_deployment" "share_add" {
  depends_on = [
    aws_api_gateway_integration.share_add,
    aws_api_gateway_integration.share_add_root,
  ]

  rest_api_id = "${aws_api_gateway_rest_api.share_add.id}"
  stage_name  = "test"
}

resource "aws_lambda_permission" "share_add" {
  # statement_id  = "AllowExecutionFromAPIGateway"
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.share_add.arn
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.share_add.execution_arn}/*/*"
}

output "share_add_url" {
  value = "${aws_api_gateway_deployment.share_add.invoke_url}"
}