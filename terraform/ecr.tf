# ECRリポジトリ
resource "aws_ecr_repository" "app" {
  name                 = "${var.project_name}-backend"
  image_tag_mutability = "MUTABLE"
}

resource "aws_ecr_lifecycle_policy" "app" {
  repository = aws_ecr_repository.app.name

  policy = jsonencode({
    rules = [{
      rulePriority = 1
      description  = "最新3つ残す"
      selection = {
        tagStatus   = "any"
        countType   = "imageCountMoreThan"
        countNumber = 3
      }
      action = {
        type = "expire"
      }
    }]
  })
}