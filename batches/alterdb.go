package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/qustavo/dotsql"
)

//VAR SQLITE DB File
var DBFile = "dbfile/footprono.sqlite"

// SQL Driver
var SDRIVER = "sqlite3"
var SCON = "./" + DBFile + "?cache=shared&mode=memory"

var CTABLE = "DB/update.sql"

func main() {
	log.Println("Altering DB")
	db, _ := sql.Open(SDRIVER, SCON)
	updateDB()
	defer db.Close()
}

func updateDB() {

	_, err2 := os.Stat(DBFile)
	if os.IsNotExist(err2) {
		log.Println("DB File doesn't exist")
		os.Create(DBFile)
		log.Println("DB File now exist")
		db, _ := sql.Open(SDRIVER, SCON)
		db.Close()
	} else {
		log.Println("DB File exist")
	}

	dot, err := dotsql.LoadFromFile(CTABLE)
	if err != nil {
		log.Println("SQL Files are causing the following issues:")
		log.Fatalln(err)
		return
	}
	log.Println("Filling Tables")
	createFromSqlFile(dot, "fill-tables")
}

func createFromSqlFile(dot *dotsql.DotSql, block string) bool {
	db, _ := sql.Open(SDRIVER, SCON)
	log.Println("Create table " + block)
	_, err := dot.Exec(db, block)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
