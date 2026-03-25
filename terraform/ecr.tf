# ECRリポジトリ
resource "aws_ecr_repository" "app" {
  name                 = "${var.project_name}-backend"
  image_tag_mutability = "MUTABLE"
}