# lambda.tf: コンテナイメージを使用したLambda
resource "aws_lambda_function" "backend" {
  function_name = "${var.project_name}-api"
  role          = aws_iam_role.lambda_exec.arn
  package_type  = "Image"
  # 注意：初回デプロイ時はECRにイメージが存在する必要があります
  image_uri     = "${aws_ecr_repository.app.repository_url}:latest"

  environment {
    variables = {
      DYNAMODB_TABLE     = aws_dynamodb_table.main.name
      SESSION_SECRET = var.session_secret
    }
  }
}