# DynamoDBテーブル
resource "aws_dynamodb_table" "main" {

  name         = "${var.project_name}-table"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "pk" # Partition Key
  range_key    = "sk" # Sort Key

  attribute {
    name = "pk"
    type = "S"
  }

  attribute {
    name = "sk"
    type = "S"
  }

  attribute {
    name = "gsi1pk"
    type = "S"
  }

  attribute {
    name = "gsi1sk"
    type = "S"
  }

  attribute {
    name = "user_id"
    type = "S"
  }

  # --- グローバルセカンダリインデックス (GSI) 1 ---
  global_secondary_index {
    name            = "GSI1"
    hash_key        = "gsi1pk"
    range_key       = "gsi1sk"
    projection_type = "ALL"
  }

  # --- グローバルセカンダリインデックス (GSI) 2 ---
  global_secondary_index {
    name            = "GSI2"
    hash_key        = "user_id"
    projection_type = "ALL"
  }
}