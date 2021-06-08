package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

var stage1File = "./ressources/stage1.json"
var stage1PronoFile = "./ressources/MatchDay1Test.json"

const (
	layoutISO = "2006-01-02"
)

func getScores(marshalled []byte, id int) {
	myString := string(marshalled)
	Team1Name := gjson.Get(myString, "#(MatchID=="+strconv.Itoa(id)+").Team1")
	Team1Score := gjson.Get(myString, "#(MatchID=="+strconv.Itoa(id)+").ScoreT1")
	Team2Name := gjson.Get(myString, "#(MatchID=="+strconv.Itoa(id)+").Team2")
	Team2Score := gjson.Get(myString, "#(MatchID=="+strconv.Itoa(id)+").ScoreT2")
	fmt.Println(Team1Name.String() + ":" + Team1Score.String())
	fmt.Println(Team2Name.String() + ":" + Team2Score.String())
	if Team2Score.Int() == Team1Score.Int() {
		fmt.Println("Egalité")
	} else if Team2Score.Int() > Team1Score.Int() {
		fmt.Printf("%s Wins", Team2Name.String())
	} else {
		fmt.Printf("%s Wins", Team1Name.String())
	}
}

func getMatches(w http.ResponseWriter, r *http.Request) {
	officialScores := readJsonMatches(stage1File)
	marshalled, _ := json.Marshal(officialScores)
	getScores(marshalled, 0)
	Respond(w, marshalled)
}

func getPlayers(w http.ResponseWriter, r *http.Request) {
	players := readJson(stage1PronoFile)
	marshalled, _ := json.Marshal(players)
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

func readJson(strFile string) []Player {
	now := time.Now()
	date := "2012-06-12"
	t, _ := time.Parse(layoutISO, date)
	diff := t.Sub(now)
	if diff.Hours() > 0 {
		fmt.Println("avant")
		fmt.Println(diff)
	} else {
		fmt.Println("apres")
		fmt.Println(diff.Hours())
	}

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
	for _, i := range m.([]interface{}) {
		var player Player
		var matchPronoSlice []PrMatch
		tempScore := 0
		player.ID = i.(map[string]interface{})["ID"].(string)
		player.Email = i.(map[string]interface{})["Email"].(string)
		player.Name = i.(map[string]interface{})["Name"].(string)
		// We itare to create real pronostics for players
		for j, oS := range officialScores {
			scoreP := 0
			var matchProno PrMatch
			matchProno.MatchID = j
			matchProno.Team1 = oS.Team1
			matchProno.ScoreT1, _ = strconv.Atoi(i.(map[string]interface{})[oS.Team1].(string))
			matchProno.Team2 = oS.Team2
			matchProno.ScoreT2, _ = strconv.Atoi(i.(map[string]interface{})[oS.Team2].(string))
			matchProno.Date = oS.Date
			//We compare Score to know who is winner according to player
			if matchProno.ScoreT1 == matchProno.ScoreT2 {
				matchProno.Winner = "PAR"
			} else if matchProno.ScoreT1 > matchProno.ScoreT2 {
				matchProno.Winner = matchProno.Team1
			} else {
				matchProno.Winner = matchProno.Team2
			}
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
		player.Matches = matchPronoSlice
		player.Score = tempScore
		players = append(players, player)
	}
	return players
}

/*
for k, v := range i.(map[string]interface{}) {
	if k == "ID" {
		fmt.Println(k + " - " + v.(string))
		player.ID = k
	}
}*/
