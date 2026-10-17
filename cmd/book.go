package cmd

import (
	"fmt"
	_ "github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
)

type BookD struct {
	BookC []struct {
		Name     string    `toml:"name"`
		CreateAt time.Time `toml:"create_at"`
	}
}

// bookCmd represents the book command
var bookCmd = &cobra.Command{
	Use:   "book",
	Short: "操作一个账簿",
	Long:  `操作一个账簿，使用-c创建，使用-d删除，使用-s选择`,
	Run: func(cmd *cobra.Command, args []string) {
		t, err := cmd.Flags().GetString("c")
		// FIXME: 逻辑有问题！！！
		// FIXME: 此处应该改用Changed("c")来判断！
		if t != ":*:" && err == nil {
			uhd, _ := os.UserHomeDir()
			os.MkdirAll(filepath.Join(uhd, ".config", "budcalc"), 0755)
			fs, err := os.OpenFile(filepath.Join(uhd, ".config", "budcalc", "book.toml"), os.O_CREATE|os.O_RDWR, 0644)
			defer fs.Close()
			if err == nil {
				bt := make([]byte, 1024)

				//NOTE:文件操作--------------------
				if i, _ := fs.Stat(); i.Size() != 0 {
					//已写过该文件
					// NOTE:(注意truncate,seek和sync!)
					n, _ := fs.Read(bt)
					_ = n
				} else {
					//空文件就直接Write
					// NOTE:读写操作应当抽函数，只有增操作直接append，改查删都全量重存！
				}
				//NOTE:文件操作结束-----------------

			}
			fmt.Printf("Success! %v创建完成", t)
		}
	},
}

func init() {
	rootCmd.AddCommand(bookCmd)
	bookCmd.Flags().StringP("c", "c", ":*:", "创建新账簿")
}
