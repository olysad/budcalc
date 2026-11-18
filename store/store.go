package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// Task 表示一条待办记录。
type Task struct {
	ID        int64
	Title     string
	Done      bool
	CreatedAt time.Time
	ClosedAt  sql.NullTime
}

// Store 封装 SQLite 存储。
type Store struct {
	db *sql.DB
}

// DefaultPath 返回默认数据库路径，Linux 下为 ~/.local/share/budcalc/budcalc.db
func DefaultPath() (string, error) {
	// XDG_DATA_HOME 是数据文件的规范位置；未设置时回退到 ~/.local/share
	if home, err := os.UserHomeDir(); err == nil {
		base := os.Getenv("XDG_DATA_HOME")
		if base == "" {
			base = filepath.Join(home, ".local", "share")
		}
		return filepath.Join(base, "budcalc", "budcalc.db"), nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "budcalc", "budcalc.db"), nil
}

// Open 打开（必要时创建）数据库，并初始化表结构。
func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("创建数据目录失败: %w", err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %w", err)
	}
	// 单连接，避免 SQLITE_BUSY
	db.SetMaxOpenConns(1)

	schema := `
CREATE TABLE IF NOT EXISTS tasks (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	title      TEXT    NOT NULL,
	done       INTEGER NOT NULL DEFAULT 0,
	created_at TEXT    NOT NULL,
	closed_at  TEXT
);
CREATE INDEX IF NOT EXISTS idx_tasks_done ON tasks(done);
`
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("初始化表结构失败: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

// Insert 新增任务，返回新 ID。
func (s *Store) Insert(title string) (int64, error) {
	res, err := s.db.Exec(`INSERT INTO tasks(title, created_at) VALUES(?, ?)`,
		title, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Update 重命名任务。
func (s *Store) Update(id int64, title string) error {
	res, err := s.db.Exec(`UPDATE tasks SET title = ? WHERE id = ?`, title, id)
	if err != nil {
		return err
	}
	return mustAffect(res, id)
}

// Delete 删除任务。
func (s *Store) Delete(id int64) error {
	res, err := s.db.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return mustAffect(res, id)
}

// Close 标记任务完成（幂等：已完成的再 close 不报错）。
func (s *Store) MarkDone(id int64) error {
	res, err := s.db.Exec(`UPDATE tasks SET done = 1, closed_at = ? WHERE id = ? AND done = 0`,
		time.Now().UTC().Format(time.RFC3339), id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n > 0 {
		return nil
	}
	// 不存在才报错；已经完成视为成功
	var exists int
	if err := s.db.QueryRow(`SELECT COUNT(1) FROM tasks WHERE id = ?`, id).Scan(&exists); err != nil {
		return err
	}
	if exists == 0 {
		return fmt.Errorf("任务 %d 不存在", id)
	}
	return nil
}

// Select 按 ID 查询单条任务。
func (s *Store) Select(id int64) (*Task, error) {
	row := s.db.QueryRow(`SELECT id, title, done, created_at, closed_at FROM tasks WHERE id = ?`, id)
	return scanOne(row)
}

// List 列出任务。today=true 只列今天创建的任务；done=false 只列未完成的。
func (s *Store) List(today, includeDone bool) ([]Task, error) {
	q := `SELECT id, title, done, created_at, closed_at FROM tasks WHERE 1=1`
	var args []any
	if !includeDone {
		q += ` AND done = 0`
	}
	if today {
		q += ` AND date(created_at) = date('now', 'localtime')`
	}
	q += ` ORDER BY done ASC, id ASC`

	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Task
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *t)
	}
	return out, rows.Err()
}

func mustAffect(res sql.Result, id int64) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("任务 %d 不存在", id)
	}
	return nil
}

type scannable interface {
	Scan(dest ...any) error
}

func scanTask(sc scannable) (*Task, error) {
	var (
		t       Task
		created string
		closed  sql.NullString
		done    int
	)
	if err := sc.Scan(&t.ID, &t.Title, &done, &created, &closed); err != nil {
		return nil, err
	}
	t.Done = done != 0
	t.CreatedAt = parseTime(created, t.CreatedAt)
	if closed.Valid {
		var ct time.Time
		if v := parseTime(closed.String, ct); !v.IsZero() {
			t.ClosedAt = sql.NullTime{Time: v, Valid: true}
		}
	}
	return &t, nil
}

func scanOne(sc scannable) (*Task, error) {
	t, err := scanTask(sc)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return t, err
}

func parseTime(s string, fallback time.Time) time.Time {
	if v, err := time.Parse(time.RFC3339, s); err == nil {
		return v.Local()
	}
	// 兜底：跳过旧格式差异
	if v, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return v.Local()
	}
	return fallback
}
