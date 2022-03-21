package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"
)

type M map[string]interface{}

func getOfficialMatches(w http.ResponseWriter, r *http.Request) {
	// Create the database handle, confirm driver is present
	db, _ := sql.Open("mysql", "lilnas:@/footprono")
	defer db.Close()

	// Connect and check the server version
	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)
	fmt.Println("Connected to:", version)

	tmpl := template.Must(template.ParseFiles("templates/matches.gohtml"))
	var officialScores2 []PrMatch
	LastMatchID := config.LastMatchID
	for _, p := range officialScores {
		if p.MatchID > LastMatchID && p.Winner == "Draw" {
			p.Winner = ""
		}
		officialScores2 = append(officialScores2, p)
	}
	tmpl.Execute(w, M{
		// We can pass as many things as we like
		"officialScores2": officialScores2,
		"stat":            stat,
		"config":          config,
	})
}
