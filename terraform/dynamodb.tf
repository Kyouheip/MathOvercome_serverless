# DynamoDBテーブル
resource "aws_dynamodb_table" "main" {
  name         = "${var.project_name}-table"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "PK" # Partition Key

  attribute {
    name = "PK"
    type = "S"
  }
}