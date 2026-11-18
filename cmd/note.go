package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olysad/budcalc/store"
	"github.com/spf13/cobra"
)

var dbPath string

// st 是延迟打开的存储句柄。
var st *store.Store

func openStore() (*store.Store, error) {
	if st != nil {
		return st, nil
	}
	path := dbPath
	if path == "" {
		p, err := store.DefaultPath()
		if err != nil {
			return nil, err
		}
		path = p
	}
	s, err := store.Open(path)
	if err != nil {
		return nil, err
	}
	st = s
	return st, nil
}

func mustStore() *store.Store {
	s, err := openStore()
	if err != nil {
		fail(err)
	}
	return s
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "budcalc:", err)
	os.Exit(1)
}

// parseID 把参数解析成任务 ID，支持 "3" 或 "#3"。
func parseID(s string) int64 {
	s = strings.TrimPrefix(strings.TrimSpace(s), "#")
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil || id <= 0 {
		fail(fmt.Errorf("无效的任务编号 %q", s))
	}
	return id
}

// noteCmd 是 note 的父命令，本身可独立运行（相当于 list --all）。
var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "管理待办事项",
	Long:  `管理待办事项：insert 新增、update 修改、delete 删除、select 查看单条。不带子命令时等同于 list --all。`,
	Run: func(cmd *cobra.Command, args []string) {
		printTasks(mustStore(), false, false)
	},
}

var noteInsertCmd = &cobra.Command{
	Use:   "insert <标题>",
	Short: "新增一条待办",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := strings.TrimSpace(strings.Join(args, " "))
		if title == "" {
			fail(fmt.Errorf("标题不能为空"))
		}
		id, err := mustStore().Insert(title)
		if err != nil {
			fail(err)
		}
		fmt.Printf("已新增 #%d %s\n", id, title)
	},
}

var noteUpdateCmd = &cobra.Command{
	Use:   "update <id> <新标题>",
	Short: "修改待办标题",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		id := parseID(args[0])
		title := strings.TrimSpace(strings.Join(args[1:], " "))
		if title == "" {
			fail(fmt.Errorf("标题不能为空"))
		}
		if err := mustStore().Update(id, title); err != nil {
			fail(err)
		}
		fmt.Printf("已更新 #%d → %s\n", id, title)
	},
}

var noteDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "删除一条待办",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := parseID(args[0])
		if err := mustStore().Delete(id); err != nil {
			fail(err)
		}
		fmt.Printf("已删除 #%d\n", id)
	},
}

var noteSelectCmd = &cobra.Command{
	Use:   "select <id>",
	Short: "查看一条待办",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := parseID(args[0])
		t, err := mustStore().Select(id)
		if err != nil {
			fail(err)
		}
		if t == nil {
			fail(fmt.Errorf("任务 %d 不存在", id))
		}
		printTask(*t)
	},
}

var (
	listToday bool
	listAll   bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出待办事项",
	Long:  `列出待办事项。默认列出全部未完成任务；--today 只看今天创建的，--all 包含已完成的。`,
	Run: func(cmd *cobra.Command, args []string) {
		printTasks(mustStore(), listToday, listAll)
	},
}

var closeCmd = &cobra.Command{
	Use:   "close <id>",
	Short: "标记一条待办为已完成",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := parseID(args[0])
		if err := mustStore().MarkDone(id); err != nil {
			fail(err)
		}
		fmt.Printf("已完成 #%d\n", id)
	},
}

func printTasks(s *store.Store, today, all bool) {
	tasks, err := s.List(today, all)
	if err != nil {
		fail(err)
	}
	if len(tasks) == 0 {
		fmt.Println("没有待办事项")
		return
	}
	for _, t := range tasks {
		printTask(t)
	}
}

func printTask(t store.Task) {
	mark := "[ ]"
	if t.Done {
		mark = "[x]"
	}
	line := fmt.Sprintf("%s #%-3d %s", mark, t.ID, t.Title)
	if t.Done && t.ClosedAt.Valid {
		line += "  (" + t.ClosedAt.Time.Format("2006-01-02 15:04") + " 完成)"
	}
	fmt.Println(line)
}

func init() {
	rootCmd.AddCommand(noteCmd, listCmd, closeCmd)

	noteCmd.AddCommand(noteInsertCmd, noteUpdateCmd, noteDeleteCmd, noteSelectCmd)

	listCmd.Flags().BoolVar(&listToday, "today", false, "只列出今天创建的待办")
	listCmd.Flags().BoolVar(&listAll, "all", false, "包含已完成的待办 (默认只看未完成)")

	rootCmd.PersistentFlags().StringVar(&dbPath, "db", "", "指定数据库文件路径 (默认 ~/.config/budcalc/budcalc.db)")
}
