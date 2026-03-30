# Cognito ユーザープール
resource "aws_cognito_user_pool" "main" {
  name = "${var.project_name}-user-pool"

  # emailでサインイン
  username_attributes      = ["email"]
  auto_verified_attributes = ["email"]

  password_policy {
    minimum_length    = 8
    require_lowercase = true
    require_numbers   = true
    require_symbols   = false
    require_uppercase = true
  }

  # メール検証設定
  verification_message_template {
    default_email_option = "CONFIRM_WITH_CODE"
  }

  # アカウント復旧
  account_recovery_setting {
    recovery_mechanism {
      name     = "verified_email"
      priority = 1
    }
  }
}

# Cognito アプリクライアント（フロントエンド用、シークレットなし）
resource "aws_cognito_user_pool_client" "main" {
  name         = "${var.project_name}-client"
  user_pool_id = aws_cognito_user_pool.main.id

  generate_secret = false

  explicit_auth_flows = [
    "ALLOW_USER_SRP_AUTH",
    "ALLOW_REFRESH_TOKEN_AUTH",
  ]

  # IDトークンにname, emailクレームを含める
  read_attributes = ["email", "name"]
}
