package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// db 是全局常量， 在init中处理没有办法传参数据库地址。
// func init() {
// 	// db, _ = sql.Open("mysql", "root:123456@tcp(139.198.189.162:13306)/FileBox?charset=utf8")
// 	db, _ = sql.Open("mysql", connStr)
// 	db.SetMaxOpenConns(1000)
// 	err := db.Ping()
// 	if err != nil {
// 		fmt.Println("Failed to connect to mysql, err:" + err.Error())
// 		os.Exit(1)
// 	}
// }

// connInit : 链接数据库
func connInit(connStr string) {
	db, _ = sql.Open("mysql", connStr)
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql, err:" + err.Error())
		os.Exit(1)
	}
}

// DBConn : 返回数据库连接对象
func DBConn(connStr string) *sql.DB {
	if db == nil {
		connInit(connStr)
	}
	return db
}

// ParseRows : 序列化返回结果
func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		//将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		checkErr(err)

		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
