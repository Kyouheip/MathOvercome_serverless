output "cloudfront_url" {
  value = "https://${aws_cloudfront_distribution.cdn.domain_name}"
}

output "api_endpoint" {
  value = aws_apigatewayv2_api.http_api.api_endpoint
}

output "ecr_repository_url" {
  value = aws_ecr_repository.app.repository_url
}