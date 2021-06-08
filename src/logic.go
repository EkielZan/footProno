package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/tidwall/gjson"
)

var stage1File = "./ressources/stage1.json"
var stage1PronoFile = "./ressources/MatchDay1Test.json"

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
	matches := readJsonMatches(stage1File)
	marshalled, _ := json.Marshal(matches)
	getScores(marshalled, 0)
	Respond(w, marshalled)
}

func getPlayers(w http.ResponseWriter, r *http.Request) {
	players := readJson(stage1PronoFile)
	marshalled, _ := json.Marshal(players)
	Respond(w, marshalled)
}

func readJsonMatches(strFile string) []Match {
	var matches []Match
	// Open our jsonFile
	raw, err := ioutil.ReadFile(strFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal([]byte(raw), &matches)
	if err != nil {
		panic(err)
	}
	return matches
}

func readJson(strFile string) []Player {
	matches := readJsonMatches(stage1File)
	// Open our jsonFile
	raw, err := ioutil.ReadFile(strFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	//we itarete the played Matches
	// Code is for debug can be removed at the end when happy
	for _, s := range matches {
		match := s
		if match.ScoreT1 == match.ScoreT2 {
			fmt.Printf("Egalité")
		} else if match.ScoreT2 > match.ScoreT1 {
			fmt.Printf("%s Wins", match.Team2)
		} else {
			fmt.Printf("%s Wins", match.Team1)
		}
		fmt.Println()
	}

	var m interface{}
	err = json.Unmarshal([]byte(raw), &m)
	if err != nil {
		panic(err)
	}
	var players []Player
	for _, i := range m.([]interface{}) {
		var player Player
		player.ID = i.(map[string]interface{})["ID"].(string)
		player.Email = i.(map[string]interface{})["Email"].(string)
		player.Name = i.(map[string]interface{})["Name"].(string)
		// We itare to create real pronostics for players
		for j, s := range matches {
			var matchProno PrMatch
			matchProno.MatchID = j
			matchProno.Team1 = s.Team1
			matchProno.ScoreT1, _ = strconv.Atoi(i.(map[string]interface{})[s.Team1].(string))
			matchProno.Team2 = s.Team2
			matchProno.ScoreT2, _ = strconv.Atoi(i.(map[string]interface{})[s.Team2].(string))
			fmt.Println(matchProno)
		}
		player.Score = 0
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
