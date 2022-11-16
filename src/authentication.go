package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log"
	"math/rand"
	"net/http"
	"regexp"

	// "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

// login authenticates the user
func login(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	check, ld := checkLogin(r.FormValue("username"), r.FormValue("code"))
	if !check {
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
		LoginDetail:   ld,
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

func checkLogin(email string, password string) (bool, loginDetail) {
	// Create the database handle, confirm driver is present
	//FORNOSQL db, _ := sql.Open("mysql", "root:@tcp(lilnas:3306)/footprono?parseTime=true")
	var ld loginDetail
	db, _ := sql.Open(SDRIVER, SCON)
	defer db.Close()
	//rows, err := db.Query("SELECT a.password,firstname,lastname,id as count FROM authentication a, users u where a.id = u.id and a.email=?;", email)
	rows, err := db.Query("SELECT password,firstname,lastname,id,token FROM users where email=?;", email)
	if err != nil {
		log.Println("No user found in the DB : " + email)
		return false, ld
	}

	var ldPwd string
	var token string
	for rows.Next() {
		rows.Scan(&ldPwd, &ld.Firstname, &ld.Lastname, &ld.Userid, &token)
		if token != "" {
			ld.Status = "token"
		}
	}

	sha_256 := sha256.New()
	sha_256.Write([]byte(password))

	sha1_hash := hex.EncodeToString(sha_256.Sum(nil))

	if ldPwd == sha1_hash {
		return true, ld
	}
	return false, ld
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

func registerForm(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := getUser(session)
	flashes := session.Flashes()
	//flashes :=
	session.Flashes()
	err = session.Save(r, w)
	tpl.ExecuteTemplate(w,
		"register.gohtml",
		M{
			// We can pass as many things as we like
			"user":  user,
			"stat":  stat,
			"flash": flashes,
		})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	firstname := r.FormValue("firstname")
	lastname := r.FormValue("lastname")
	password := r.FormValue("password")
	const minEntropyBits = 60
	errP := passwordvalidator.Validate(password, minEntropyBits)
	if errP != nil {
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		flashes = "Password is not complex enough."
		session.AddFlash(flashes)
		session.Save(r, w)
		http.Redirect(w, r, "/registerForm", http.StatusFound)
		return
	} else {
		if !check {
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
			data := make([]byte, 10)
			for i := range data {
				data[i] = byte(rand.Intn(256))
			}
			sha_256 := sha256.New()
			sha_256.Write([]byte(data))
			token := hex.EncodeToString(sha_256.Sum(nil))
			db, _ := sql.Open(SDRIVER, SCON)
			defer db.Close()
			rows, _ := db.Query("select max(id)+1 from users;")
			userid := 0
			for rows.Next() {
				rows.Scan(&userid)
			}
			sha_pwd := sha256.New()
			sha_pwd.Write([]byte(password))

			pwd_hash := hex.EncodeToString(sha_pwd.Sum(nil))

			stm, _ := db.Prepare("INSERT into users (id,firstname,lastname,email,password,token) VALUES (?,?,?,?,?,?);")
			if err != nil {
				log.Println(err.Error())
				return
			}
			_ = stm.QueryRow(userid, firstname, lastname, email, pwd_hash, token).Scan()

			if err != nil {
				log.Println(err.Error())
				return
			}
			sendMail(email, token)
			http.Redirect(w, r, "/registerDone", http.StatusFound)
		}
	}
}

func registerDone(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := getUser(session)
	flashes := session.Flashes()
	//flashes :=
	session.Flashes()
	err = session.Save(r, w)
	tpl.ExecuteTemplate(w,
		"welcome.gohtml",
		M{
			// We can pass as many things as we like
			"user":  user,
			"stat":  stat,
			"flash": flashes,
		})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func isEmailValidING(e string) bool {
	if e == "blacksadum@gmail.com" {
		return true
	}
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@ing.com`)
	return emailRegex.MatchString(e)
}
