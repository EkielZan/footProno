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
)

//All working files
var stage1File = "./ressources/stage1.json"
var stage1PronoFile = "./ressources/MatchDay1Test.json"
var stage2PronoFile = "./ressources/MatchDay2Test.json"
var stage3PronoFile2 = "./ressources/MatchDay3Test.json"

//var stage3PronoFile = "./ressources/MatchDay3Test.json"
var champFile = "./ressources/champList.json"
var statusFile = "./ressources/status.json"
var configFile = "./ressources/config.json"

//
var stat Statistics
var config Config
var officialScores []PrMatch
var players []Player
var champPlayer []Player
var scoredPlayers []Player

//Preload json File in Memory
func preLoad() {
	log.Println("Reading Players Champions")
	champPlayer = readChampion(champFile)
	log.Println("Reading Official Match Scores")
	officialScores = readJsonMatches(stage1File)
	log.Println("Reading Saved Players")
	//shortPlayers = readSavedPlayers()
	log.Println("Reading Players Pronostics")
	load()
	log.Println("Reading Config")
	config = loadConfig()
	saveConfig(config)
}

//reload json File in Memory
func reload() {
	officialScores = readJsonMatches(stage1File)
	load()
	config = loadConfig()
	saveConfig(config)
}

func load() {
	players = initJsonPlayers(stage1PronoFile, 0)
	players = updatePlayers(stage2PronoFile, players, 1)
	players = updatePlayers(stage3PronoFile2, players, 2)
	scoredPlayers = calculateScore()
	sort.Slice(scoredPlayers, func(i, j int) bool {
		return scoredPlayers[i].Score > scoredPlayers[j].Score
	})
	scoredPlayers = setRank(scoredPlayers)
}

//Read datas to Fill Struct
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

func readJsonMatches(strFile string) []PrMatch {
	var officialScores []PrMatch
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
	return officialScores
}

func initJsonPlayers(strFile string, stage int) []Player {
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
			if oS.Stage == stage {
				var matchProno PrMatch
				matchProno.MatchID = oS.MatchID
				//We check the score if 10 then it's 0
				score, done := convertInterface(i.(map[string]interface{})[oS.Team1])
				if score == 10 {
					score = 0
				}
				matchProno.Done = done
				matchProno.Team1 = oS.Team1
				matchProno.ScoreP1 = score
				//We check the score if 10 then it's 0
				score, done = convertInterface(i.(map[string]interface{})[oS.Team2])
				if score == 10 {
					score = 0
				}
				matchProno.Done = done
				matchProno.Team2 = oS.Team2
				matchProno.ScoreP2 = score
				matchProno.Date = oS.Date
				//We compare Score to know who is winner according to player
				if matchProno.ScoreP1 == matchProno.ScoreP2 {
					matchProno.Winner = "Draw"
				} else if matchProno.ScoreP1 > matchProno.ScoreP2 {
					matchProno.Winner = matchProno.Team1
				} else {
					matchProno.Winner = matchProno.Team2
				}
				//we add the filled strucut to the slice
				matchPronoSlice = append(matchPronoSlice, matchProno)
			}
		}
		player.Matches = matchPronoSlice
		players = append(players, player)
	}
	return players
}
func convertInterface(myInterface interface{}) (int, bool) {
	done := false
	test := myInterface
	tmpStr := ""
	switch v := test.(type) {
	case int:
		fmt.Println(v)
	case string:
		tmpStr = test.(string)
		done = true
	default:
		tmpStr = ""
	}
	score, _ := strconv.Atoi(tmpStr)
	return score, done
}
func updatePlayers(strFile string, pPlayers []Player, stage int) []Player {
	// Open our jsonFile
	var tPlayers []Player
	fmt.Println(strFile)
	playersT := initJsonPlayers(strFile, stage)
	for _, player := range playersT {
		for _, pPlayer := range pPlayers {
			if pPlayer.Name == player.Name {
				pPlayer.Matches = append(pPlayer.Matches, player.Matches...)
				tPlayers = append(tPlayers, pPlayer)
			}
		}
	}
	return tPlayers
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

//Other func
func calculateScore() []Player {
	var tempPlayers []Player
	savedPlayers := readSavedPlayers()
	LastMatchID := config.LastMatchID
	for _, player := range players {
		playerScoreTemp := 0
		var tempMatches []PrMatch
		for _, match := range player.Matches {
			MatchscoreTemp := 0
			for _, oS := range officialScores {
				MatchscoreTemp = 0
				if match.MatchID <= LastMatchID {
					if match.Team1 == oS.Team1 && match.Team2 == oS.Team2 {
						if match.Done {
							if match.Winner == oS.Winner {
								MatchscoreTemp += 1
								if match.ScoreP1 == oS.ScoreT1 && match.ScoreP2 == oS.ScoreT2 {
									MatchscoreTemp += 2
								}
							}
						}
						match.ScoreP = MatchscoreTemp
						match.ScoreT1 = oS.ScoreT1
						match.ScoreT2 = oS.ScoreT2
						match.OWinner = oS.Winner
						playerScoreTemp += match.ScoreP
					}
				}
			}
			tempMatches = append(tempMatches, match)
		}
		for _, savedPlayer := range savedPlayers {
			if savedPlayer.ID == player.ID {
				player.Amount = savedPlayer.Amount
				if player.Amount > 0 {
					player.Status = "Up"
				} else if player.Amount == 0 {
					player.Status = "Stay"
				} else {
					player.Status = "Down"
				}
			}
		}
		player.Matches = tempMatches
		player.Score = playerScoreTemp
		tempPlayers = append(tempPlayers, player)
	}
	return tempPlayers
}

func setRank(players []Player) []Player {
	for idx_p := range players {
		players[idx_p].Rank = idx_p + 1
	}
	return players
}

// Save datas
func savePlayers(players []Player) {
	savedPlayers := readSavedPlayers()
	var toSaveSlice []ShortPlayer
	//max := len(players)
	for _, player := range scoredPlayers {
		var toSavePlayer ShortPlayer
		for _, savedPlayer := range savedPlayers {
			if savedPlayer.ID == player.ID {
				toSavePlayer.LastPosition = savedPlayer.CurrentPosition
			}
		}
		toSavePlayer.Name = player.Name
		toSavePlayer.CurrentPosition = player.Rank
		toSavePlayer.Amount = toSavePlayer.LastPosition - toSavePlayer.CurrentPosition
		toSavePlayer.ID = player.ID
		toSaveSlice = append(toSaveSlice, toSavePlayer)
	}
	sort.Slice(toSaveSlice, func(i, j int) bool {
		return toSaveSlice[i].ID > toSaveSlice[j].ID
	})
	marshalled, _ := json.MarshalIndent(toSaveSlice, "", " ")
	_ = ioutil.WriteFile(statusFile, marshalled, 0644)
}

//Config Management
func loadConfig() Config {
	// Open our jsonFile
	raw, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal([]byte(raw), &config)
	if err != nil {
		panic(err)
	}
	return config
}

func saveConfig(config Config) {
	marshalled, _ := json.MarshalIndent(config, "", " ")
	_ = ioutil.WriteFile(configFile, marshalled, 0644)
}

func refresh(w http.ResponseWriter, r *http.Request) {
	config = loadConfig()
	if config.Refresh {
		fmt.Println("We refresh the rank difference")
		savePlayers(players)
		config.Refresh = false
		saveConfig(config)
	}
}
