package handler

import (
	"encoding/json"
	"go-todo-api/model"
	"go-todo-api/repository"
	"net/http"
	"strconv"
	"strings"
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

// IDに基づいてTodoを取得するハンドラー
func (h *TodoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// 1. リクエストメソッドのチェック
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 2. クエリパラメータからIDを取得
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "IDが指定されていません", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(pathParts[2])
	if err != nil {
	http.Error(w, "無効なID", http.StatusBadRequest)
		return
	}

	// 3. Todoの取得
	todo, err := h.repo.GetByID(int64(id))
	if err != nil {
		http.Error(w, "Todoが見つかりません", http.StatusNotFound)
		return
	}

	// 4. レスポンスの返却
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// IDに基づいてTodoを更新するハンドラー
func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	// 1. リクエストメソッドのチェック
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 2. クエリパラメータからIDを取得
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "IDが指定されていません", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "無効なID", http.StatusBadRequest)
		return
	}

	// 3. リクエストボディのパース
	var todo model.Todo
	// リクエストボディをJSONデコードしてtodo構造体に格納
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "無効なリクエストボディ", http.StatusBadRequest)
		return
	}
	todo.ID = id

	if err := h.repo.Update(todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// IDに基づいてTodoを削除するハンドラー
func (h *TodoHandler) Delete(e http.ResponseWriter, r *http.Request) {
	// 1. リクエストメソッドのチェック
	if r.Method != http.MethodDelete {
		e.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 2. クエリパラメータからIDを取得
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(e, "IDが指定されていません", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(e, "無効なID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(int64(id)); err != nil {
		http.Error(e, err.Error(), http.StatusInternalServerError)
		return
	}

	e.WriteHeader(http.StatusNoContent)
}