package handler

import (
	"encoding/json"
	"go-todo-api/model"
	"go-todo-api/repository"
	"net/http"
)

type TodoHandler struct {
	repo *repository.TodoRepository
}

func NewTodoHandler(repo *repository.TodoRepository) *TodoHandler {
	return &TodoHandler{
		repo: repo,
	}
}

// w : クライアントへ返すレスポンスを書き込む
// r : クライアントからのリクエストの情報を含む
func (h *TodoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// リクエストメソッドのチェック
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	todos, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// JSONレスポンスを返す todoのスライスをJSONにエンコードしてレスポンスボディに書き込む
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	// 1. リクエストメソッドのチェック
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 2. リクエストボディのパース
	var todo model.Todo
	// リクエストボディをJSONデコードしてtodo構造体に格納
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "無効なリクエストボディ", http.StatusBadRequest)
		return
	}

	// 3. Todoの作成
	id, err := h.repo.Create(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. レスポンスの返却
	// 作成されたTodoのIDを含めて、レスポンスとして返す
	todo.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}