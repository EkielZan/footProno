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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8090" //localhost
	}
	serverHost := os.Getenv("SERVER_NAME")
	if serverHost == "" {
		serverHost = "localhost" //localhost
	}

	updateJavaScript(serverPort, serverHost)
	stat.Version = Version
	log.Println("Version:\t", Version)
	log.Println("Running Web Server Api on " + serverHost + " " + serverPort)
	router := mux.NewRouter()
	static := spaHandler{staticPath: "static", indexPath: "index.html"}

	preLoad()
	log.Println("Preparing to Serve Api")
	router.HandleFunc("/matches", getMatches)
	router.HandleFunc("/getScores", getScore)
	router.HandleFunc("/player/{id}", getPlayer)
	router.HandleFunc("/getMiscData", getMiscData)
	router.HandleFunc("/refresh", refresh)
	router.HandleFunc("/scp/{id}", scoreByPlayer)
	router.HandleFunc("/GetLeaderboard", getLeaderboard)

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
	//log.Fatal(srv.ListenAndServe())
	log.Println("Ready to received calls")
	log.Fatal(srv.ListenAndServeTLS("certs/server.crt", "certs/server.key"))

}
