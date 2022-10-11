package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log"
	"net/http"

	// "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
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

func checkLogin(email string, password string) bool {
	// Create the database handle, confirm driver is present
	//FORNOSQL db, _ := sql.Open("mysql", "root:@tcp(lilnas:3306)/footprono?parseTime=true")
	db, _ := sql.Open(SDRIVER, SCON)
	defer db.Close()
	log.Println(email + " - " + password)
	//rows, err := db.Query("SELECT a.password,firstname,lastname,id as count FROM authentication a, users u where a.id = u.id and a.email=?;", email)
	rows, err := db.Query("SELECT password,firstname,lastname,id as count FROM users where email=?;", email)
	if err != nil {
		log.Println("The err : " + err.Error())
		log.Println("No user found in the DB : " + email)
		return false
	}

	var ld loginDetail
	for rows.Next() {
		rows.Scan(&ld.Password, &ld.Firstname, &ld.Lastname, &ld.Userid)
	}

	sha_256 := sha256.New()
	sha_256.Write([]byte(password))

	sha1_hash := hex.EncodeToString(sha_256.Sum(nil))
	log.Println(sha1_hash)

	if ld.Password == sha1_hash {
		log.Println("Good Password")
		return true
	}
	log.Println("Wrong Password")
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
