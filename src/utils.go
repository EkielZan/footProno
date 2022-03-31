package main

import (
	"database/sql"
	"io"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
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
	if !ok {
		return User{
			Username:      "",
			Authenticated: false,
		}
	}
	return user
}

func health(w http.ResponseWriter, r *http.Request) {
	// Create the database handle, confirm driver is present
	db, _ := sql.Open("mysql", "lilnas:root@/footprono?parseTime=true")
	defer db.Close()
	// Connect and check the server version
	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)
	io.WriteString(w, "I'm healthy and The DB version is : "+version)
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
