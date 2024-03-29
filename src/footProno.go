package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var Version = "Development"

//VAR SQLITE DB File
var DBFile = "dbfile/footprono.sqlite"

// SQL Driver
var SDRIVER = "sqlite3"
var SCON = "./" + DBFile + "?cache=shared&mode=memory"

//Creation of the tables
var CTABLE = "DB/queries.sql"

//Filling the tables
var FTABLE = "DB/data.sql"

//Global Access variables
var stat Stats
var cookieName = "footProno-secure-cookie"

// store will hold all session data
var store *sessions.CookieStore

// tpl holds all parsed templates
var tpl *template.Template

func init() {
	authKeyOne := []byte("whatwedointheshadows")
	encryptionKeyOne := []byte("themandalorian22")
	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}
	gob.Register(User{})
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "4000" //localhost
	}
	serverHost := os.Getenv("SERVER_NAME")
	if serverHost == "" {
		serverHost = "localhost" //localhost
	}
	//updateJavaScript(serverPort, serverHost)
	stat.Version = Version

	log.Println("Version:\t", stat.Version)
	log.Println("Running Web Server on http://" + serverHost + ":" + serverPort)
	router := mux.NewRouter()

	initDatabase()

	log.Println("Preparing to Serve pages")

	router.HandleFunc("/", index)
	router.HandleFunc("/gom", getOfficialMatches)
	router.HandleFunc("/gt", getTeams)
	router.HandleFunc("/health", health)
	router.HandleFunc("/login", login)
	router.HandleFunc("/addPronos", addPronos)
	router.HandleFunc("/registerForm", registerForm)
	router.HandleFunc("/register", register)
	router.HandleFunc("/registerDone", registerDone)
	router.HandleFunc("/logout", logout)
	router.HandleFunc("/prons", getPronostics)
	router.HandleFunc("/updateStuff", updateStuff)
	router.HandleFunc("/getLeaderBoard", getLeaderBoard)
	router.HandleFunc("/about", about)

	fileServer := http.FileServer(http.Dir("static"))
	router.PathPrefix("/js").Handler(http.StripPrefix("/", fileServer))
	router.PathPrefix("/css").Handler(http.StripPrefix("/", fileServer))
	router.PathPrefix("/img").Handler(http.StripPrefix("/", fileServer))

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + serverPort,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Ready to receive calls")
	log.Fatal(srv.ListenAndServe())
}
