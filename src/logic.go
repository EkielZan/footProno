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
var datefile = "ressources/matchday.json"
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
	//shortPlayers = readSavedPlayers()
	log.Println("Reading Players Pronostics")
	players = readJsonPlayers(stage1PronoFile)
}

//reload json File in Memory
func reload() {
	log.Println("Reading Official Match Scores")
	officialScores = readJsonMatches(stage1File)
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
		player.ID = idx + 1
		player.Email = i.(map[string]interface{})["Email"].(string)
		player.Name = i.(map[string]interface{})["Name"].(string)
		for _, v := range champPlayer {
			if v.Name == player.Name {
				player.Champ = v.Champ
			}
		}
		// We itare to create real pronostics for players
		for _, oS := range officialScores {
			var matchProno PrMatch
			matchProno.MatchID = oS.MatchID
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
				matchProno.Winner = "Draw"
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
	reload()
	marshalled, _ := json.MarshalIndent(stat, "", " ")
	Respond(w, marshalled)
}

func getScore(w http.ResponseWriter, r *http.Request) {
	reload()
	scoredPlayers := calculateScore()
	sort.Slice(scoredPlayers, func(i, j int) bool {
		return scoredPlayers[i].Score > scoredPlayers[j].Score
	})
	scoredPlayers = setRank(scoredPlayers)
	//players = computeRank(players)
	savePlayers(scoredPlayers)
	marshalled, _ := json.MarshalIndent(scoredPlayers, "", " ")
	Respond(w, marshalled)
}

//Other func
func calculateScore() []Player {
	var tempPlayers []Player
	/* now := time.Now() */
	LastMatchID := 3
	/*Time Debug
	date := "2021-06-12"
	now, _ = time.Parse(layoutISO, date)
	*/
	for _, player := range players {
		playerScoreTemp := 0
		var tempMatches []PrMatch
		for _, match := range player.Matches {
			//var tempMatch PrMatch
			MatchscoreTemp := 0
			for _, oS := range officialScores {
				MatchscoreTemp = 0
				/* date := oS.Date */
				/* t, _ := time.Parse(layoutISO, date) */
				/* diff := t.Before(now) */
				if match.MatchID <= LastMatchID {
					if match.Team1 == oS.Team1 && match.Team2 == oS.Team2 {
						if match.Winner == oS.Winner {
							MatchscoreTemp += 1
							if match.ScoreT1 == oS.ScoreT1 && match.ScoreT2 == oS.ScoreT2 {
								MatchscoreTemp += 2
							}
						}
						match.ScoreP = MatchscoreTemp
						playerScoreTemp += match.ScoreP
					}
				}
			}
			//tempMatch = match
			tempMatches = append(tempMatches, match)
		}
		player.Matches = tempMatches
		player.Score = playerScoreTemp
		tempPlayers = append(tempPlayers, player)
	}
	return tempPlayers
}

func setRank(players []Player) []Player {
	for idx_p, _ := range players {
		players[idx_p].Rank = idx_p + 1
	}
	return players
}

// Save datas
func savePlayers(players []Player) {
	now := time.Now()
	var toSaveSlice []ShortPlayer
	for idx, player := range players {
		var toSavePlayer ShortPlayer
		toSavePlayer.Name = player.Name
		var tempPosition LastPosition
		tempPositionSlice := player.Positions
		tempPosition.Position = idx + 1
		tempPosition.ScoreDate = now.Format("2006-01-02")
		toSavePlayer.Positions = append(tempPositionSlice, tempPosition)
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
