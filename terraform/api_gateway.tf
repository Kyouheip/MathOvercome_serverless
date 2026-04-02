resource "aws_apigatewayv2_api" "http_api" {
  name          = "${var.project_name}-gateway"
  protocol_type = "HTTP"

  # OPTIONSプリフライトはAPI Gatewayで処理し、Lambdaを呼び出さない
  # GET/POST等の実際のCORSヘッダーはLambda（Gin）側で付与
  cors_configuration {
    allow_origins     = ["https://${aws_cloudfront_distribution.cdn.domain_name}"]
    allow_methods     = ["GET", "POST", "OPTIONS"]
    allow_headers     = ["Content-Type", "Authorization"]
    allow_credentials = true
    max_age           = 43200
  }
}

# Cognito JWT オーソライザー
resource "aws_apigatewayv2_authorizer" "cognito" {
  api_id           = aws_apigatewayv2_api.http_api.id
  authorizer_type  = "JWT"
  identity_sources = ["$request.header.Authorization"]
  name             = "${var.project_name}-cognito-authorizer"

  jwt_configuration {
    audience = [aws_cognito_user_pool_client.main.id]
    issuer   = "https://cognito-idp.${var.aws_region}.amazonaws.com/${aws_cognito_user_pool.main.id}"
  }
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.http_api.id
  name        = "$default"
  auto_deploy = true

  default_route_settings {
    throttling_burst_limit = 20
    throttling_rate_limit  = 10
  }
}

resource "aws_apigatewayv2_integration" "lambda" {
  api_id           = aws_apigatewayv2_api.http_api.id
  integration_type = "AWS_PROXY"
  integration_uri  = aws_lambda_function.backend.invoke_arn

  # JWT認証済みのCognitoクレームをヘッダーに注入（Lambda側でX-User-Subを読む）
  request_parameters = {
    "overwrite:header.X-User-Sub"  = "$context.authorizer.claims.sub"
    "overwrite:header.X-User-Name" = "$context.authorizer.claims.name"
  }
}

resource "aws_apigatewayv2_route" "any" {
  api_id             = aws_apigatewayv2_api.http_api.id
  route_key          = "ANY /{proxy+}"
  target             = "integrations/${aws_apigatewayv2_integration.lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

# Lambda側にAPI Gatewayからの呼び出しを許可
resource "aws_lambda_permission" "api_gw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.backend.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}
