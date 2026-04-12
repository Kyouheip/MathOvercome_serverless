package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/spf13/cobra"

	"github.com/Kyouheip/MathOvercome_serverless/internal/repository"
	"github.com/Kyouheip/MathOvercome_serverless/internal/service"
)

var (
	userSub     string
	testSessSvc service.TestSessionServicer
	mypageSvc   service.MypageServicer
)

var rootCmd = &cobra.Command{
	Use:   "mathovercome",
	Short: "MathOvercome CLI",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return setupServices()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&userSub, "user", "", "ユーザーID")
}

func setupServices() error {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-northeast-1"
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return fmt.Errorf("AWS設定の読み込みに失敗: %w", err)
	}

	var opts []func(*dynamodb.Options)
	if endpoint := os.Getenv("DYNAMODB_ENDPOINT"); endpoint != "" {
		opts = append(opts, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	}
	client := dynamodb.NewFromConfig(cfg, opts...)

	repo := repository.NewRepository(client)
	testSessSvc = service.NewTestSessionService(repo)
	mypageSvc = service.NewMypageService(repo)

	return nil
}
