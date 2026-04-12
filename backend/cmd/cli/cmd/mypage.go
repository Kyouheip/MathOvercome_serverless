package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Kyouheip/MathOvercome_serverless/internal/model"
)

var mypageCmd = &cobra.Command{
	Use:   "mypage",
	Short: "マイページ情報を表示する",
	RunE: func(cmd *cobra.Command, args []string) error {
		if userSub == "" {
			return fmt.Errorf("--user フラグが必要です")
		}
		userName, _ := cmd.Flags().GetString("name")

		user := &model.User{
			Sub:      userSub,
			UserName: userName,
		}

		data, err := mypageSvc.GetUserData(user)
		if err != nil {
			return fmt.Errorf("マイページ取得失敗: %w", err)
		}

		fmt.Printf("ユーザー: %s\n", data.UserName)
		fmt.Printf("テストセッション数: %d\n\n", len(data.TestSessDtos))

		for _, sess := range data.TestSessDtos {
			fmt.Printf("--- セッション %d (%s) ---\n", sess.SessionID, sess.StartTime)
			fmt.Printf("正答率: %d/%d\n", sess.CorrectCount, sess.Total)
			if len(sess.WeakCategories) > 0 {
				fmt.Printf("苦手分野: %v\n", sess.WeakCategories)
			}
			fmt.Println()
		}

		return nil
	},
}

func init() {
	mypageCmd.Flags().String("name", "", "ユーザー名 (表示用)")
	rootCmd.AddCommand(mypageCmd)
}
