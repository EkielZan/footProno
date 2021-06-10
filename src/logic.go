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

//All working files
var stage1File = "./ressources/stage1.json"
var stage1PronoFile = "./ressources/MatchDay1Test.json"
var champFile = "./ressources/champList.json"
var statusFile = "./ressources/status.json"
var configFile = "./ressources/config.json"
var stat Statistics
var config Config

const (
	layoutISO = "2006-01-02"
)

//preLoad()
func preLoad() {
	_ = readJsonPlayers(stage1PronoFile)
}

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
	players = checkingRank(players)
	savePlayers(players)
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

//Func Logic

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
	saveConfig(config)
	return officialScores
}

func readJsonPlayers(strFile string) []Player {
	ButProno := 0
	ButReal := 0
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
	champPlayer := readChampion(champFile)
	var players []Player
	for idx, i := range m.([]interface{}) {
		var player Player
		var matchPronoSlice []PrMatch
		tempScore := 0
		player.ID = idx
		player.Email = i.(map[string]interface{})["Email"].(string)
		player.Name = i.(map[string]interface{})["Name"].(string)
		for _, v := range champPlayer {
			if v.Name == player.Name {
				player.Champ = v.Champ
			}
		}
		// We itare to create real pronostics for players
		for j, oS := range officialScores {
			date := oS.Date
			t, _ := time.Parse(layoutISO, date)
			diff := t.Before(now)
			if diff {
				config.LastMatchDate = date
				config.LastMatchID = oS.MatchID
				saveConfig(config)
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
				ButProno += score
				//We check the score if 10 then it's 0
				score, _ = strconv.Atoi(i.(map[string]interface{})[oS.Team2].(string))
				if score == 10 {
					score = 0
				}
				matchProno.Team2 = oS.Team2
				matchProno.ScoreT2 = score
				ButProno += score
				matchProno.Date = oS.Date
				//We compare Score to know who is winner according to player
				if matchProno.ScoreT1 == matchProno.ScoreT2 {
					matchProno.Winner = "PAR"
				} else if matchProno.ScoreT1 > matchProno.ScoreT2 {
					matchProno.Winner = matchProno.Team1
				} else {
					matchProno.Winner = matchProno.Team2
				}
				ButReal += oS.ScoreT1 + oS.ScoreT2
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
		stat.ButProno = ButProno
		stat.ButReal = ButReal
	}
	return players
}

func readChampion(strFile string) []Player {
	var champPlayer []Player
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
	for idx, i := range m.([]interface{}) {
		var player Player
		player.ID = idx
		player.Email = i.(map[string]interface{})["Email"].(string)
		player.Name = i.(map[string]interface{})["Name"].(string)
		player.Champ = i.(map[string]interface{})["Champ"].(string)
		champPlayer = append(champPlayer, player)
	}
	return champPlayer
}

func savePlayers(players []Player) {
	var toSaveSlice []ShortPlayer
	for idx, player := range players {
		var toSavePlayer ShortPlayer
		toSavePlayer.Name = player.Name
		toSavePlayer.LastPos = idx
		toSavePlayer.Score = player.Score
		toSaveSlice = append(toSaveSlice, toSavePlayer)
	}
	marshalled, _ := json.MarshalIndent(toSaveSlice, "", " ")
	_ = ioutil.WriteFile(statusFile, marshalled, 0644)
}

func saveConfig(config Config) {
	marshalled, _ := json.MarshalIndent(config, "", " ")
	_ = ioutil.WriteFile(configFile, marshalled, 0644)
}

func checkingRank(players []Player) []Player {
	shortPlayers := readSavedPlayers()
	for _, shortPlayer := range shortPlayers {
		for idx, player := range players {
			if shortPlayer.Name == player.Name {
				diff := idx - shortPlayer.LastPos
				fmt.Printf("%d  -  %d - %d ", idx, shortPlayer.LastPos, diff)
				fmt.Println("")
				if diff > 0 {
					player.Status = "Down"
				} else if diff < 0 {
					player.Status = "Up"
				} else {
					player.Status = "SQ"
				}
			}
		}
	}
	fmt.Println("--------")
	return players

}

func readSavedPlayers() []ShortPlayer {
	var savedPlayer []ShortPlayer
	// Open our jsonFile
	raw, err := ioutil.ReadFile(statusFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal([]byte(raw), &savedPlayer)
	if err != nil {
		panic(err)
	}
	return savedPlayer
}
