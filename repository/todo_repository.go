package repository

import (
	"database/sql"
	"go-todo-api/model"
	"log"
)

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository {
		db: db,
	}
}

// 全てのTodoを取得する
func (r *TodoRepository) GetAll() ([]model.Todo, error) {
	query := "SELECT id, title, content, done, created_at From todos"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []model.Todo
	for rows.Next() {
		var todo model.Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Done, &todo.CreatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *TodoRepository) Create(todo model.Todo) (int64, error) {
	// プリペアードステートメントを使ってSQLインジェクションを防ぐ
	query := "INSERT INTO todos (title, content) VALUES (?, ?)"

	// Execはクエリを実行し、結果を返す
	result, err := r.db.Exec(query, todo.Title, todo.Content)
	if err != nil {
		return 0, err
	}

	log.Println("result: ", result)

	// 挿入されたレコードのIDを取得
	return result.LastInsertId()
}