package main

import (
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

type M map[string]interface{}

func scoreByPlayer(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/player.gohtml"))
	reload()
	playerID, _ := strconv.Atoi(mux.Vars(r)["id"])
	var player Player
	for _, p := range scoredPlayers {
		if p.ID == playerID {
			player = p
			var tempMatches []PrMatch
			LastMatchID := config.LastMatchID
			for _, m := range player.Matches {
				if m.MatchID <= LastMatchID {
					tempMatches = append(tempMatches, m)
				}
			}
			player.Matches = tempMatches
		}
	}

	tmpl.Execute(w, M{
		// We can pass as many things as we like
		"player": player,
		"stat":   stat,
	})
}

func getLeaderboard(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/leaderboard.gohtml"))
	reload()
	//tmpl.Execute(w, scoredPlayers)
	tmpl.Execute(w, M{
		// We can pass as many things as we like
		"scoredPlayers": scoredPlayers,
		"stat":          stat,
	})
}

func getTeams(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/teams.gohtml"))
	reload()
	//tmpl.Execute(w, scoredPlayers)
	tmpl.Execute(w, M{
		// We can pass as many things as we like
		"teams": teams,
		"stat":  stat,
	})
}

func getOfficialMatches(w http.ResponseWriter, r *http.Request) {
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
	})
}
