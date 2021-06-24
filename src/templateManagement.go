package main

import (
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

type M map[string]interface{}

func scoreByPlayer(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/player.html"))
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
	tmpl.Execute(w, player)
}

func getLeaderboard(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("templates/leaderboard.html"))
	reload()
	//tmpl.Execute(w, scoredPlayers)
	tmpl.Execute(w, M{
		// We can pass as many things as we like
		"scoredPlayers": scoredPlayers,
		"stat":          stat,
	})
}
