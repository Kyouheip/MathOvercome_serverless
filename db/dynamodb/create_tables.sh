#!/bin/bash
# DynamoDB テーブル作成スクリプト (シングルテーブル設計)
# 使用前に REGION を変更してください

REGION="ap-northeast-1"

echo "Creating table: MathOvercome..."
aws dynamodb create-table \
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
  --region "$REGION"

echo "Table creation command sent."
echo "テーブルがACTIVEになるまで少し待ってからupload.shを実行してください。"
