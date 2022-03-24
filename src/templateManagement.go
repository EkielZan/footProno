package main

import (
	"net/http"
	"text/template"
)

type M map[string]interface{}

func getOfficialMatches(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("templates/matches.gohtml"))
	var officialScores2 []PrMatch
	LastMatchID := config.Lastmatch
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
