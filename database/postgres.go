package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Init() {
	// 连接数据库
	db_, err := sql.Open("postgres", "postgres://postgres:Forever0.@localhost:5432/?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db = db_

	// 检查数据库连接
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to PostgreSQL database!")
}

func Insert(user string, path string, reqCount int64, resCount int64) {
	stmt, err := db.Prepare("INSERT INTO flow_record (name, path, req_count, res_count) VALUES ($1, $2, $3, $4)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user, path, reqCount, resCount)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Data inserted successfully!")
}

func Close() {
	if db != nil {
		db.Close()
	}
}
