package mydb

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	// blank import, bind it to database/sql
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	tmp, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:13306)/fileserver?charset=utf8")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Failed to connect to mysql server")
		os.Exit(1)
	}
	db = tmp
	// if open the server successfully
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	// check if the db connection is dead
	err = db.Ping()

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Failed to ping the mysql server")
		os.Exit(1)
	}
}

// DBConn : return the mysql connection obj
func DBConn() *sql.DB {
	return db
}
