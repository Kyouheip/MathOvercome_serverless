#!/bin/sh
# Docker Compose 用初期化スクリプト
# DynamoDB Local にテーブル作成 + データ投入を行う
set -e

ENDPOINT="http://dynamodb:8000"
DATA_DIR="/dynamodb/data"

echo "Checking table: MathOvercome..."
if aws dynamodb describe-table --endpoint-url "$ENDPOINT" --table-name MathOvercome --region ap-northeast-1 > /dev/null 2>&1; then
  echo "Table already exists, skipping creation and data upload."
  exit 0
fi

echo "Creating table: MathOvercome..."
aws dynamodb create-table \
  --endpoint-url "$ENDPOINT" \
  --table-name MathOvercome \
  --attribute-definitions \
    AttributeName=pk,AttributeType=S \
    AttributeName=sk,AttributeType=S \
    AttributeName=gsi1pk,AttributeType=S \
    AttributeName=gsi1sk,AttributeType=S \
  --key-schema \
    AttributeName=pk,KeyType=HASH \
    AttributeName=sk,KeyType=RANGE \
  --global-secondary-indexes '[
    {
      "IndexName": "GSI1",
      "KeySchema": [
        {"AttributeName": "gsi1pk", "KeyType": "HASH"},
        {"AttributeName": "gsi1sk", "KeyType": "RANGE"}
      ],
      "Projection": {"ProjectionType": "ALL"}
    }
  ]' \
  --billing-mode PAY_PER_REQUEST \
  --region ap-northeast-1

echo "Waiting for table to be active..."
aws dynamodb wait table-exists \
  --endpoint-url "$ENDPOINT" \
  --table-name MathOvercome \
  --region ap-northeast-1

echo "Uploading data..."
for f in \
  problems_01.json problems_02.json \
  choices_01.json choices_02.json choices_03.json choices_04.json choices_05.json choices_06.json choices_07.json \
  test_sessions.json \
  session_problems.json; do
  echo "  $f"
  aws dynamodb batch-write-item \
    --endpoint-url "$ENDPOINT" \
    --request-items "file://$DATA_DIR/$f" \
    --region ap-northeast-1
done

echo "Setup complete."
