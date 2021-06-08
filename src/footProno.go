package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var Version = "Development"

/* func home(w http.ResponseWriter, r *http.Request) {
	Respond(w, []byte(`{"message": "hello people"}`))
} */

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8000" //localhost
	}

	updateJavaScript(serverPort)
	log.Println("Version:\t", Version)
	log.Println("Running Web Server Api on port " + serverPort)
	router := mux.NewRouter()
	static := spaHandler{staticPath: "static", indexPath: "index.html"}

	router.HandleFunc("/matches", getMatches)
	router.HandleFunc("/players", getPlayers)
	router.HandleFunc("/playersByScore", getOrderedPlayers)
	router.HandleFunc("/player/{id}", getPlayer)
	router.HandleFunc("/stats", getStat)

	//TODO Remove if not necessary
	router.PathPrefix("/").Handler(static)
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
	log.Fatal(srv.ListenAndServe())

}
