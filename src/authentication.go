package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
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

// secret displays the secret message for authorized users
func secret(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := getUser(session)
	if auth := user.Authenticated; !auth {
		session.AddFlash("You don't have access!")
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	tpl.ExecuteTemplate(w, "secret.gohtml", user.Username)
}

func checkLogin(login string, password string) bool {
	// Create the database handle, confirm driver is present
	db, _ := sql.Open("mysql", "root:@tcp(lilnas:3306)/footprono")
	defer db.Close()
	log.Println(login + " - " + password)
	rows, err := db.Query("SELECT a.password,firstname,lastname,userid as count FROM authentication a, users u where a.userid = u.id and a.email=?;", login)
	if err != nil {
		log.Println("The err : " + err.Error())
		log.Println("No user found in the DB : " + login)
		return false
	}
	type loginDetail struct {
		password  string
		firstname string
		lastname  string
		userid    int
	}
	var ld loginDetail
	for rows.Next() {
		rows.Scan(&ld.password, &ld.firstname, &ld.lastname, &ld.userid)
	}
	if ld.password == password {
		return true
	}
	return false
}
