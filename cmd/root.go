/*
Copyright © 2026 olysad <lacak3736@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "budcalc",
	Short: "面向个人的终端预算与待办工具",
	Long: `budcalc 是一个跑在终端里的小工具。

当前包含 note 待办管理：
  budcalc note insert <标题>      新增一条待办
  budcalc note update <id> <标题> 修改标题
  budcalc note delete <id>        删除
  budcalc note select <id>        查看单条
  budcalc list [--today|--all]    列出待办
  budcalc close <id>              标记完成`,
	// 不带子命令时直接列出全部未完成待办
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	if st != nil {
		_ = st.Close()
	}
}
