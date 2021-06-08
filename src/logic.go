package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var stage1File = "./ressources/stage1.json"
var stage1PronoFile = "./ressources/MatchDay1Test.json"
var stat Statistics

const (
	layoutISO = "2006-01-02"
)

// Api Calls
func getMatches(w http.ResponseWriter, r *http.Request) {
	officialScores := readJsonMatches(stage1File)
	marshalled, _ := json.Marshal(officialScores)
	Respond(w, marshalled)
}

func getPlayers(w http.ResponseWriter, r *http.Request) {
	players := readJsonPlayers(stage1PronoFile)
	marshalled, _ := json.Marshal(players)
	Respond(w, marshalled)
}

func getOrderedPlayers(w http.ResponseWriter, r *http.Request) {
	players := readJsonPlayers(stage1PronoFile)
	sort.Slice(players, func(i, j int) bool {
		return players[i].Score > players[j].Score
	})
	marshalled, _ := json.Marshal(players)
	Respond(w, marshalled)
}

func getPlayer(w http.ResponseWriter, r *http.Request) {
	playerID, _ := strconv.Atoi(mux.Vars(r)["id"])
	players := readJsonPlayers(stage1PronoFile)
	player := players[playerID]
	marshalled, _ := json.Marshal(player)
	Respond(w, marshalled)
}

func getStat(w http.ResponseWriter, r *http.Request) {
	marshalled, _ := json.Marshal(stat)
	Respond(w, marshalled)
}

func readJsonMatches(strFile string) []Match {
	var officialScores []Match
	// Open our jsonFile
	raw, err := ioutil.ReadFile(strFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal([]byte(raw), &officialScores)
	if err != nil {
		panic(err)
	}
	return officialScores
}

func readJsonPlayers(strFile string) []Player {
	stat.ButProno = 0
	stat.ButReal = 0
	stat.Fall = ""
	stat.Rise = ""
	now := time.Now()
	/*Time Debug*/
	date := "2021-06-12"
	now, _ = time.Parse(layoutISO, date)
	/* */
	officialScores := readJsonMatches(stage1File)
	// Open our jsonFile
	raw, err := ioutil.ReadFile(strFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var m interface{}
	err = json.Unmarshal([]byte(raw), &m)
	if err != nil {
		panic(err)
	}
	var players []Player
	for idx, i := range m.([]interface{}) {
		var player Player
		var matchPronoSlice []PrMatch
		tempScore := 0
		player.ID = idx
		player.Email = i.(map[string]interface{})["Email"].(string)
		player.Name = i.(map[string]interface{})["Name"].(string)
		// We itare to create real pronostics for players
		for j, oS := range officialScores {
			date := oS.Date
			t, _ := time.Parse(layoutISO, date)
			diff := t.Before(now)
			if diff {
				scoreP := 0
				var matchProno PrMatch
				matchProno.MatchID = j
				//We check the score if 10 then it's 0
				score, _ := strconv.Atoi(i.(map[string]interface{})[oS.Team1].(string))
				if score == 10 {
					score = 0
				}
				matchProno.Team1 = oS.Team1
				matchProno.ScoreT1 = score
				stat.ButProno += score
				//We check the score if 10 then it's 0
				score, _ = strconv.Atoi(i.(map[string]interface{})[oS.Team2].(string))
				if score == 10 {
					score = 0
				}
				matchProno.Team2 = oS.Team2
				matchProno.ScoreT2 = score
				stat.ButProno += score
				matchProno.Date = oS.Date
				//We compare Score to know who is winner according to player
				if matchProno.ScoreT1 == matchProno.ScoreT2 {
					matchProno.Winner = "PAR"
				} else if matchProno.ScoreT1 > matchProno.ScoreT2 {
					matchProno.Winner = matchProno.Team1
				} else {
					matchProno.Winner = matchProno.Team2
				}
				stat.ButReal += oS.ScoreT1 + oS.ScoreT2
				//We compare Score to know who is Real winner
				if oS.ScoreT1 == oS.ScoreT2 {
					oS.Winner = "PAR"
				} else if oS.ScoreT1 > oS.ScoreT2 {
					oS.Winner = oS.Team1
				} else {
					oS.Winner = oS.Team2
				}
				if matchProno.Winner == oS.Winner {
					scoreP += 1
					if matchProno.ScoreT1 == oS.ScoreT1 && matchProno.ScoreT2 == oS.ScoreT2 {
						scoreP += 2
					}
				}
				matchProno.ScoreP = scoreP
				tempScore += scoreP
				//we add the filled strucut to the slice
				matchPronoSlice = append(matchPronoSlice, matchProno)
			}
		}
		player.Matches = matchPronoSlice
		player.Score = tempScore
		players = append(players, player)
	}
	return players
}
