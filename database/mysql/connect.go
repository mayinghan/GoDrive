package mysql

import (
	"database/sql"
	"fmt"
	"os"
	// blank import, bind it to database/sql
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:13306)/filedb?charset=utf8")
	if err != nil {
		fmt.Println("Failed to connect to mysql server")
		os.Exit(1)
	}
	// if open the server successfully
	db.SetMaxOpenConns(100)
	// check if the db connection is dead
	err = db.Ping()

	if err != nil {
		fmt.Println("Failed to ping the mysql server")
		os.Exit(1)
	}
}

// DBConnection : return the mysql connection obj
func DBConnection() *sql.DB {
	return db
}
