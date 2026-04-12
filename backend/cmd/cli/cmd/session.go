package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "テストセッション操作",
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "テストセッションを開始する",
	RunE: func(cmd *cobra.Command, args []string) error {
		if userSub == "" {
			return fmt.Errorf("--user フラグが必要です")
		}
		includeIntegers, _ := cmd.Flags().GetBool("integers")

		sess, err := testSessSvc.CreateTestSess(userSub, includeIntegers)
		if err != nil {
			return fmt.Errorf("セッション作成失敗: %w", err)
		}

		fmt.Printf("セッションID: %d\n", sess.ID)
		fmt.Printf("問題数: %d\n", len(sess.SessionProblems))
		return nil
	},
}

var problemCmd = &cobra.Command{
	Use:   "problem <idx>",
	Short: "問題を表示する (0始まり)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if userSub == "" {
			return fmt.Errorf("--user フラグが必要です")
		}
		idx, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("idxは数値で指定してください")
		}
		sessionID, _ := cmd.Flags().GetUint64("session")

		p, err := testSessSvc.GetProblem(sessionID, userSub, idx)
		if err != nil {
			return fmt.Errorf("問題取得失敗: %w", err)
		}

		fmt.Printf("[問題 %d/%d]\n", idx+1, p.Total)
		fmt.Printf("Q: %s\n\n", p.Question)
		for i, c := range p.Choices {
			fmt.Printf("%d) %s  (id:%d)\n", i+1, c.ChoiceText, c.ID)
		}
		if p.Hint != "" {
			fmt.Printf("\nヒント: %s\n", p.Hint)
		}
		if p.SelectedID != nil {
			fmt.Printf("\n回答済み (選択肢ID: %d)\n", *p.SelectedID)
		}
		return nil
	},
}

var answerCmd = &cobra.Command{
	Use:   "answer <idx>",
	Short: "回答を送信する (0始まり)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if userSub == "" {
			return fmt.Errorf("--user フラグが必要です")
		}
		idx, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("idxは数値で指定してください")
		}
		sessionID, _ := cmd.Flags().GetUint64("session")
		choiceID, _ := cmd.Flags().GetInt64("choice")

		if err := testSessSvc.SubmitAnswer(sessionID, userSub, idx, &choiceID); err != nil {
			return fmt.Errorf("回答送信失敗: %w", err)
		}

		fmt.Println("回答を送信しました")
		return nil
	},
}

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "テストセッションを対話形式で進める",
	RunE: func(cmd *cobra.Command, args []string) error {
		if userSub == "" {
			return fmt.Errorf("--user フラグが必要です")
		}
		sessionID, _ := cmd.Flags().GetUint64("session")

		scanner := bufio.NewScanner(os.Stdin)

		for idx := 0; ; idx++ {
			p, err := testSessSvc.GetProblem(sessionID, userSub, idx)
			if err != nil {
				fmt.Println("\n全問題が終わりました")
				break
			}

			fmt.Printf("\n[問題 %d/%d]\n", idx+1, p.Total)
			fmt.Printf("Q: %s\n\n", p.Question)
			for i, c := range p.Choices {
				fmt.Printf("  %d) %s  (id:%d)\n", i+1, c.ChoiceText, c.ID)
			}
			if p.Hint != "" {
				fmt.Printf("\nヒント: %s\n", p.Hint)
			}

			fmt.Print("\n選択肢IDを入力 (スキップ: s): ")
			if !scanner.Scan() {
				break
			}
			input := strings.TrimSpace(scanner.Text())
			if input == "s" || input == "" {
				fmt.Println("スキップしました")
				continue
			}

			choiceID, err := strconv.ParseInt(input, 10, 64)
			if err != nil {
				fmt.Println("無効な入力です。スキップします")
				continue
			}

			if err := testSessSvc.SubmitAnswer(sessionID, userSub, idx, &choiceID); err != nil {
				fmt.Printf("送信エラー: %v\n", err)
				continue
			}
			fmt.Println("回答しました")
		}

		return nil
	},
}

func init() {
	createCmd.Flags().Bool("integers", false, "整数問題を含める")

	problemCmd.Flags().Uint64("session", 0, "セッションID")
	problemCmd.MarkFlagRequired("session")

	answerCmd.Flags().Uint64("session", 0, "セッションID")
	answerCmd.MarkFlagRequired("session")
	answerCmd.Flags().Int64("choice", 0, "選択肢ID")
	answerCmd.MarkFlagRequired("choice")

	playCmd.Flags().Uint64("session", 0, "セッションID")
	playCmd.MarkFlagRequired("session")

	sessionCmd.AddCommand(createCmd, problemCmd, answerCmd, playCmd)
	rootCmd.AddCommand(sessionCmd)
}
