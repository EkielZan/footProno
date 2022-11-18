package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	//VAR SQLITE DB File
	var DBFile = "dbfile/footprono.sqlite"

	// SQL Driver
	var SDRIVER = "sqlite3"
	var SCON = "./" + DBFile + "?cache=shared&mode=memory"

	log.Println("Altering DB")

	db, _ := sql.Open(SDRIVER, SCON)

	defer db.Close()
}
