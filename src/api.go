package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Api Calls
func getMatches(w http.ResponseWriter, r *http.Request) {
	var officialScores2 []PrMatch
	LastMatchID := config.LastMatchID
	for _, p := range officialScores {
		if p.MatchID > LastMatchID && p.Winner == "Draw" {
			p.Winner = ""
		}
		officialScores2 = append(officialScores2, p)
	}
	marshalled, _ := json.MarshalIndent(officialScores2, "", " ")
	Respond(w, marshalled)
}

func getPlayer(w http.ResponseWriter, r *http.Request) {
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
	marshalled, _ := json.MarshalIndent(player, "", " ")
	Respond(w, marshalled)
}

func getStat(w http.ResponseWriter, r *http.Request) {
	reload()
	marshalled, _ := json.MarshalIndent(stat, "", " ")
	Respond(w, marshalled)
}

func getScore(w http.ResponseWriter, r *http.Request) {
	reload()
	marshalled, _ := json.MarshalIndent(scoredPlayers, "", " ")
	Respond(w, marshalled)
}
