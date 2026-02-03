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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Hello, Go Todo API!"})
	})
	http.HandleFunc("/todos", h.GetAll)

	// 4. サーバーの起動
	port := ":8080"
	fmt.Println("サーバー起動 port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
