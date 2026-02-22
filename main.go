package main

import (
	"encoding/json"
	"fmt"
	"go-todo-api/handler"
	"go-todo-api/repository"
	"log"
	"net/http"
)

func main() {
	log.Println("start Go Todo API")

	// 1. DBの初期化
	db, err := repository.NewDB()
	if err != nil {
		log.Fatal("DB接続失敗:", err)
	}
	defer db.Close()

	// 2. 依存関係の注入
	repo := repository.NewTodoRepository(db)
	h := handler.NewTodoHandler(repo)

	// 3. ルーティングの設定
	// 第二引数にはハンドラー関数の「型」が一致するものを渡す
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Hello, Go Todo API!"})
	})
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetAll(w, r)
		case http.MethodPost:
			h.Create(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	// URLパスが/todos/で始まるリクエストを処理するハンドラー
	// 例: /todos/1, /todos/2など
	http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		// URLパスにIDがなければ、Todoの一覧取得や新規作成のハンドラーに処理を渡す
		if r.URL.Path == "/todos/" || r.URL.Path == "/todos" {
			switch r.Method {
			case http.MethodGet:
				h.GetAll(w, r)
			case http.MethodPost:
				h.Create(w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		} else {
			// IDに基づくTodoの取得
			h.GetByID(w, r)
		}
	})

	// 4. サーバーの起動
	port := ":8080"
	fmt.Println("サーバー起動 port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
