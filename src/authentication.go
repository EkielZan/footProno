package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

// login authenticates the user
func login(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	check := checkLogin(r.FormValue("username"), r.FormValue("code"))
	if check != true {
		session.AddFlash("Password is incorrect")
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	username := r.FormValue("username")
	user := &User{
		Username:      username,
		Authenticated: true,
	}
	session.Values["user"] = user
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// logout revokes authentication for a user
func logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["user"] = User{}
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func checkLogin(login string, password string) bool {
	// Create the database handle, confirm driver is present
	db, _ := sql.Open("mysql", "root:@tcp(lilnas:3306)/footprono?parseTime=true")
	defer db.Close()
	log.Println(login + " - " + password)
	rows, err := db.Query("SELECT a.password,firstname,lastname,userid as count FROM authentication a, users u where a.userid = u.id and a.email=?;", login)
	if err != nil {
		log.Println("The err : " + err.Error())
		log.Println("No user found in the DB : " + login)
		return false
	}

	var ld loginDetail
	for rows.Next() {
		rows.Scan(&ld.Password, &ld.Firstname, &ld.Lastname, &ld.Userid)
	}
	if ld.Password == password {
		return true
	}
	return false
}

func checkAuthentication(w http.ResponseWriter, r *http.Request, user User, session *sessions.Session) bool {
	if auth := user.Authenticated; !auth {
		session.AddFlash("You don't have access!")
		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return false
	}
	return true
}
