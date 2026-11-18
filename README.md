# Budcalc
一个面向个人的终端预算财务管理工具。

![github](https://img.shields.io/badge/github-repo-blue?logo=github&style=for-the-badge)
![dev](https://img.shields.io/badge/STATUS-IN%20DEV-red?style=for-the-badge&logo=devbox)
![license](https://img.shields.io/badge/license-mit-brown?style=for-the-badge)
# 安装
```bash
go install github.com/olysad/budcalc@latest
```

# 快速开始
```bash
# 新增待办
budcalc note insert 交房租

# 看看还有哪些没做（默认只列未完成的）
budcalc list

# 只看今天建的
budcalc list --today

# 连已完成的也一起看
budcalc list --all

# 改个标题
budcalc note update 1 交本月房租

# 做完一条
budcalc close 1

# 看单条 / 删掉
budcalc note select 1
budcalc note delete 1
```

# 配置
数据库默认放在 `$XDG_DATA_HOME/budcalc/budcalc.db`（未设置时为 `~/.local/share/budcalc/budcalc.db`）。
想换个位置就加 `--db`：

```bash
budcalc --db ~/notes/todo.db list
```

存储使用 SQLite，单文件、可备份，也可以用任何 SQLite 工具直接查看。

# 许可证
[MIT License](LICENSE)
