package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
var officialScores []Match
var players []Player
var champPlayer []Player
var shortPlayers []ShortPlayer

const (
	layoutISO = "2006-01-02"
)

//Preload json File in Memory
func preLoad() {
	log.Println("Reading Players Champions")
	champPlayer = readChampion(champFile)
	log.Println("Reading Official Match Scores")
	officialScores = readJsonMatches(stage1File)
	log.Println("Reading Saved Players")
	shortPlayers = readSavedPlayers()
	log.Println("Reading Players Pronostics")
	players = readJsonPlayers(stage1PronoFile)
}

//Function to Fill Struct
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

func readJsonPlayers(strFile string) []Player {
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
			var matchProno PrMatch
			matchProno.MatchID = j
			//We check the score if 10 then it's 0
			score, _ := strconv.Atoi(i.(map[string]interface{})[oS.Team1].(string))
			if score == 10 {
				score = 0
			}
			matchProno.Team1 = oS.Team1
			matchProno.ScoreT1 = score
			//We check the score if 10 then it's 0
			score, _ = strconv.Atoi(i.(map[string]interface{})[oS.Team2].(string))
			if score == 10 {
				score = 0
			}
			matchProno.Team2 = oS.Team2
			matchProno.ScoreT2 = score
			matchProno.Date = oS.Date
			//We compare Score to know who is winner according to player
			if matchProno.ScoreT1 == matchProno.ScoreT2 {
				matchProno.Winner = "PAR"
			} else if matchProno.ScoreT1 > matchProno.ScoreT2 {
				matchProno.Winner = matchProno.Team1
			} else {
				matchProno.Winner = matchProno.Team2
			}
			//we add the filled strucut to the slice
			matchPronoSlice = append(matchPronoSlice, matchProno)
		}
		player.Matches = matchPronoSlice
		players = append(players, player)
	}
	return players
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
	//We compare Score to know who is Real winner
	for idx, oS := range officialScores {
		if oS.ScoreT1 == oS.ScoreT2 {
			oS.Winner = "Draw"
		} else if oS.ScoreT1 > oS.ScoreT2 {
			oS.Winner = oS.Team1
		} else {
			oS.Winner = oS.Team2
		}
		officialScores[idx] = oS
	}
	saveConfig(config)
	return officialScores
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

// Api Calls
func getMatches(w http.ResponseWriter, r *http.Request) {
	marshalled, _ := json.MarshalIndent(officialScores, "", " ")
	Respond(w, marshalled)
}

func getPlayers(w http.ResponseWriter, r *http.Request) {
	marshalled, _ := json.MarshalIndent(players, "", " ")
	Respond(w, marshalled)
}

func getOrderedPlayers(w http.ResponseWriter, r *http.Request) {
	sort.Slice(players, func(i, j int) bool {
		return players[i].Score > players[j].Score
	})
	players = checkingRank(players)
	savePlayers(players)
	marshalled, _ := json.MarshalIndent(players, "", " ")
	Respond(w, marshalled)
}

func getPlayer(w http.ResponseWriter, r *http.Request) {
	playerID, _ := strconv.Atoi(mux.Vars(r)["id"])
	players := readJsonPlayers(stage1PronoFile)
	player := players[playerID]
	marshalled, _ := json.MarshalIndent(player, "", " ")
	Respond(w, marshalled)
}

func getStat(w http.ResponseWriter, r *http.Request) {
	marshalled, _ := json.MarshalIndent(stat, "", " ")
	Respond(w, marshalled)
}

func calculateScore(w http.ResponseWriter, r *http.Request) {
	var tempPlayers []Player
	now := time.Now()
	/*Time Debug*/
	date := "2021-06-12"
	now, _ = time.Parse(layoutISO, date)
	/* */
	for _, player := range players {
		var tempPlayer Player
		var tempMatches []Match
		for _, match := range player.Matches {
			var tempMatch Match
			for _, oS := range officialScores {
				date := oS.Date
				t, _ := time.Parse(layoutISO, date)
				diff := t.Before(now)
				if diff {
					if matchProno.ScoreT1 == matchProno.ScoreT2 {
						matchProno.Winner = "PAR"
					} else if matchProno.ScoreT1 > matchProno.ScoreT2 {
						matchProno.Winner = matchProno.Team1
					} else {
						matchProno.Winner = matchProno.Team2
					}
				}
			}
			tempMatches = append(tempMatches, tempMatch)
		}
		tempPlayers = append(tempPlayers, tempPlayer)
	}
	marshalled, _ := json.MarshalIndent(tempPlayers, "", " ")
	Respond(w, marshalled)
}

func checkingRank(players []Player) []Player {
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

// Save datas
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
