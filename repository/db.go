package repository

import (
	"database/sql"
	"fmt"
	"log"

	// Importすると、database/sqlパッケージからMySQLドライバが利用可能になる
	_ "github.com/go-sql-driver/mysql"
)

func NewDB()(*sql.DB, error) {
	// 接続情報の構築
	// paseTime=trueをつけることで、MySQLのDatetTime, TimeStamp型をGOのtime.Time型で値を受け取れるようになる
	dsn := "root:marumaru39@tcp(127.0.0.1:3306)/todo_db?parseTime=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// 接続確認
	if err := db.Ping(); err != nil {
		log.Fatal("DB接続失敗\n", err)
		return nil, err
	}

	fmt.Println("DB接続成功")
	return db, nil
}