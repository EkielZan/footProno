package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log"
	"net/http"
	"regexp"

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
	//rows, err := db.Query("SELECT a.password,firstname,lastname,id as count FROM authentication a, users u where a.id = u.id and a.email=?;", email)
	rows, err := db.Query("SELECT password,firstname,lastname,id as count FROM users where email=?;", email)
	if err != nil {
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

	if ld.Password == sha1_hash {
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

func register(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	email := r.FormValue("email")
	check := isEmailValidING(email)
	flashes := ""

	if check != true {
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		flashes = "Email is not ending with ing.com or isn't correctly formatted."
		session.AddFlash(flashes)
		session.Save(r, w)
		http.Redirect(w, r, "/registerForm", http.StatusFound)
		return
	} else {
		log.Println("Email is fully correct now onto temporary token generation")
	}

}
func isEmailValidING(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@ing.com`)
	return emailRegex.MatchString(e)
}
