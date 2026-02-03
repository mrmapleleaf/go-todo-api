package handler

import (
	"encoding/json"
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