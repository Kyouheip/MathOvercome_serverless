package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/Kyouheip/MathOvercome_serverless/internal/router"
)

func main() {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-northeast-1"
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	var opts []func(*dynamodb.Options)
	if endpoint := os.Getenv("DYNAMODB_ENDPOINT"); endpoint != "" {
		opts = append(opts, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	}
	client := dynamodb.NewFromConfig(cfg, opts...)

	r := router.New(client)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("starting server on :%s", port)
	r.Run(":" + port)
}
