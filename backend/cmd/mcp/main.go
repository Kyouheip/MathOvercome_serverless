package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/Kyouheip/MathOvercome_serverless/internal/model"
	"github.com/Kyouheip/MathOvercome_serverless/internal/repository"
	"github.com/Kyouheip/MathOvercome_serverless/internal/service"
)

func main() {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-northeast-1"
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("AWS設定の読み込みに失敗: %v", err)
	}

	var opts []func(*dynamodb.Options)
	if endpoint := os.Getenv("DYNAMODB_ENDPOINT"); endpoint != "" {
		opts = append(opts, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	}
	client := dynamodb.NewFromConfig(cfg, opts...)

	repo := repository.NewRepository(client)
	testSessSvc := service.NewTestSessionService(repo)
	mypageSvc := service.NewMypageService(repo)

	s := server.NewMCPServer("mathovercome", "1.0.0")

	// create_test_session
	s.AddTool(
		mcp.NewTool("create_test_session",
			mcp.WithDescription("数学のテストセッションを作成する。セッションIDを返すので以降のツールで使う。"),
			mcp.WithString("user_sub", mcp.Required(), mcp.Description("ユーザーID")),
			mcp.WithBoolean("include_integers", mcp.Description("整数問題を含めるか（デフォルト: false）")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			userSub := req.GetString("user_sub", "")
			includeIntegers := req.GetBool("include_integers", false)

			sess, err := testSessSvc.CreateTestSess(userSub, includeIntegers)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			return mcp.NewToolResultText(fmt.Sprintf(
				"セッションを作成しました。\nセッションID: %d\n問題数: %d",
				sess.ID, len(sess.SessionProblems),
			)), nil
		},
	)

	// get_problem
	s.AddTool(
		mcp.NewTool("get_problem",
			mcp.WithDescription("テストセッションの問題を取得する。問題文と選択肢を返す。"),
			mcp.WithString("user_sub", mcp.Required(), mcp.Description("ユーザーID")),
			mcp.WithString("session_id", mcp.Required(), mcp.Description("セッションID（create_test_sessionで返された文字列をそのまま使う）")),
			mcp.WithNumber("index", mcp.Required(), mcp.Description("問題のインデックス（0始まり）")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			userSub := req.GetString("user_sub", "")
			sessionIDStr := req.GetString("session_id", "")
			sessionID, err := strconv.ParseUint(sessionIDStr, 10, 64)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("session_idが不正です: %v", err)), nil
			}
			idx := int(req.GetFloat("index", 0))

			p, err := testSessSvc.GetProblem(sessionID, userSub, idx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			result := fmt.Sprintf("[問題 %d/%d]\nQ: %s\n\n選択肢:\n", idx+1, p.Total, p.Question)
			for i, c := range p.Choices {
				result += fmt.Sprintf("%d) %s (id:%d)\n", i+1, c.ChoiceText, c.ID)
			}
			if p.Hint != "" {
				result += fmt.Sprintf("\nヒント: %s", p.Hint)
			}

			return mcp.NewToolResultText(result), nil
		},
	)

	// submit_answer
	s.AddTool(
		mcp.NewTool("submit_answer",
			mcp.WithDescription("問題に回答を送信する。選択肢IDはget_problemで取得したidを使う。"),
			mcp.WithString("user_sub", mcp.Required(), mcp.Description("ユーザーID")),
			mcp.WithString("session_id", mcp.Required(), mcp.Description("セッションID（create_test_sessionで返された文字列をそのまま使う）")),
			mcp.WithNumber("index", mcp.Required(), mcp.Description("問題のインデックス（0始まり）")),
			mcp.WithString("choice_id", mcp.Required(), mcp.Description("選択肢ID（get_problemで取得したidを文字列で渡す）")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			userSub := req.GetString("user_sub", "")
			sessionIDStr := req.GetString("session_id", "")
			sessionID, err := strconv.ParseUint(sessionIDStr, 10, 64)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("session_idが不正です: %v", err)), nil
			}
			idx := int(req.GetFloat("index", 0))
			choiceIDStr := req.GetString("choice_id", "")
			choiceID, err := strconv.ParseInt(choiceIDStr, 10, 64)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("choice_idが不正です: %v", err)), nil
			}

			if err := testSessSvc.SubmitAnswer(sessionID, userSub, idx, &choiceID); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			return mcp.NewToolResultText("回答を送信しました"), nil
		},
	)

	// get_mypage
	s.AddTool(
		mcp.NewTool("get_mypage",
			mcp.WithDescription("ユーザーの過去のテスト結果・正答率・苦手分野を取得する。"),
			mcp.WithString("user_sub", mcp.Required(), mcp.Description("ユーザーID")),
			mcp.WithString("user_name", mcp.Description("ユーザー名（表示用）")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			userSub := req.GetString("user_sub", "")
			userName := req.GetString("user_name", "")

			user := &model.User{Sub: userSub, UserName: userName}
			data, err := mypageSvc.GetUserData(user)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			result := fmt.Sprintf("ユーザー: %s\nテストセッション数: %d\n\n", data.UserName, len(data.TestSessDtos))
			for _, sess := range data.TestSessDtos {
				result += fmt.Sprintf("--- セッション %d (%s) ---\n正答率: %d/%d\n",
					sess.SessionID, sess.StartTime, sess.CorrectCount, sess.Total)
				if len(sess.WeakCategories) > 0 {
					result += fmt.Sprintf("苦手分野: %v\n", sess.WeakCategories)
				}
				result += "\n"
			}

			return mcp.NewToolResultText(result), nil
		},
	)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("MCP server error: %v", err)
	}
}
