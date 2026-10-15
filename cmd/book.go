package cmd

import (
	"fmt"
	_ "github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	_ "os"
	"time"
)

type bookD struct {
	bookC []struct {
		name     string    `toml:"name"`
		createAt time.Time `toml:"create_at"`
	}
}

// bookCmd represents the book command
var bookCmd = &cobra.Command{
	Use:   "book",
	Short: "操作一个账簿",
	Long:  `操作一个账簿，使用-c创建，使用-d删除，使用-s选择`,
	Run: func(cmd *cobra.Command, args []string) {
		t, err := cmd.Flags().GetString("c")
		if t != ":*:" && err == nil {

			fmt.Printf("Success! %v创建完成", t)
		}
	},
}

func init() {
	rootCmd.AddCommand(bookCmd)
	bookCmd.Flags().StringP("c", "c", ":*:", "创建新账簿")
}
