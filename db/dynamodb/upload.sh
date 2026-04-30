#!/bin/bash
# DynamoDB データアップロードスクリプト (シングルテーブル設計)
# 使用前に REGION を変更してください

REGION="ap-northeast-1"
DATA_DIR="$(dirname "$0")/data"

upload() {
  local file="$DATA_DIR/$1"
  echo "Uploading $1..."
  aws dynamodb batch-write-item \
    --request-items "file://$file" \
    --region "$REGION"
  echo "Done: $1"
}

# 問題 (42件)
upload problems_01.json
upload problems_02.json

# 選択肢 (168件)
upload choices_01.json
upload choices_02.json
upload choices_03.json
upload choices_04.json
upload choices_05.json
upload choices_06.json
upload choices_07.json

# テストセッション
upload test_sessions.json

# セッション問題 (24件)
upload session_problems.json

echo "All data uploaded."
