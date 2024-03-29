package main

import (
	"bytes"
	"crypto/tls"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	// _ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"github.com/qustavo/dotsql"
	"gopkg.in/gomail.v2"
)

//Respond Write to the httpWrite the content of data
func Respond(w http.ResponseWriter, data []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
	w.Write([]byte(data))
}

func getUser(s *sessions.Session) User {
	val := s.Values["user"]
	var user = User{}
	user, ok := val.(User)
	var ld loginDetail
	if !ok {
		return User{
			Username:      "",
			LoginDetail:   ld,
			Authenticated: false,
		}
	}
	return user
}

func health(w http.ResponseWriter, r *http.Request) {
	// Create the database handle, confirm driver is present
	//FORNOSQL db, _ := sql.Open("mysql", "lilnas:root@/footprono?parseTime=true")
	db, _ := sql.Open(SDRIVER, SCON)
	defer db.Close()
	// Connect and check the server version
	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)
	io.WriteString(w, "I'm healthy and The DB version is : "+version)
}

/*
func manageSession(w http.ResponseWriter, r *http.Request) (User, bool) {
	// Manage Sessions and authentication
	session, err := store.Get(r, cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	user := getUser(session)
	if !checkAuthentication(w, r, user, session) {
		return false
	}
	return user, true
}
*/
func initDatabase() {

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
	log.Println("Creating Tables")
	retStat := createFromSqlFile(dot, "create-tables")
	if !retStat {
		log.Println("Tables are already created.")
		log.Println("No need to refilled them.")
	} else {
		dot, err = dotsql.LoadFromFile(FTABLE)
		if err != nil {
			log.Println("SQL Files are causing the following issues:")
			log.Fatalln(err)
			return
		}
		log.Println("Filling Tables")
		createFromSqlFile(dot, "fill-tables")
	}
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

/*
func updateJavaScript(port string, host string) {
	input, err := ioutil.ReadFile("static/js/app.tpl.js")
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "###PORT###") {
			lines[i] = "port=\"" + port + "\""
		}
		if strings.Contains(line, "###HOST###") {
			lines[i] = "host=\"" + host + "\""
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile("static/js/app.js", []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
*/

func sendMail(to string, token string) {

	from := "blacksadum@gmail.com"
	password := os.Getenv("APP_TOKEN")

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "blacksadum@gmail.com", "footProno Manager")

	// Set E-Mail receivers
	m.SetHeader("To", to)

	// Set E-Mail subject
	m.SetHeader("Subject", "FootProno Temporary Token")

	var t bytes.Buffer

	tpl.ExecuteTemplate(&t,
		"mail.gohtml",
		M{
			"user":  to,
			"token": token,
		})

	result := t.String()
	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", result)

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
}
