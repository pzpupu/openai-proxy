package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Init() {
	// 连接数据库
	db_, err := sql.Open("postgres", "postgres://postgres:Forever0.@postgres:5432/?sslmode=disable")
	if err != nil {
		log.Fatal("Connected to PostgreSQL Error: ", err)
	}
	db = db_

	// 检查数据库连接
	err = db.Ping()
	if err != nil {
		log.Fatal("Check PostgreSQL Connected Error: ", err)
	}
	log.Println("Connected to PostgreSQL database!")

	// 检查表是否存在，不存在则自动建表
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS flow_record (id SERIAL PRIMARY KEY, name VARCHAR(20), path VARCHAR(255), req_count integer, res_count integer, created_time timestamp default CURRENT_TIMESTAMP)")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Table created successfully!")

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
