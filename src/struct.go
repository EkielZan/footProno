package main

import (
	"database/sql"
	"time"
)

type ShortPlayer struct {
	ID              int    `json:"ID"`
	Name            string `json:"Name"`
	Amount          int    `json:"Amount"`
	CurrentPosition int    `json:"CurrentPosition"`
	LastPosition    int    `json:"LastPosition"`
}

type LastPosition struct {
	ScoreDate string `json:"ScoreDate"`
}

// PrMatch
type PrMatch struct {
	MatchID int    `json:"MatchID"`
	Team1   string `json:"Team1"`
	ScoreT1 int    `json:"ScoreT1"`
	ScoreP1 int    `json:"ScoreP1"`
	Team2   string `json:"Team2"`
	ScoreT2 int    `json:"ScoreT2"`
	ScoreP2 int    `json:"ScoreP2"`
	Date    string `json:"Date"`
	Stage   int    `json:"Stage"`
	Winner  string `json:"Winner"`
	OWinner string `json:"OWinner"`
	ScoreP  int    `json:"ScoreP"`
	Done    bool   `json:"Done"`
}

//Team
type Team struct {
	ID          int            `json:"ID"`
	Name        sql.NullString `json:"Name"`
	Active      int            `json:"Active"`
	Groupid     string         `json:"Group"`
	Point       int            `json:"Point"`
	Win         int            `json:"Win"`
	Drawn       int            `json:"Draw"`
	Lose        int            `json:"Lose"`
	Goalfor     int            `json:"GFor"`
	Goalagainst int            `json:"GAgainst"`
}

type Match struct {
	ID      int       `json:"ID"`
	Stage   int       `json:"Stage"`
	Date    time.Time `json:"Date"`
	Teama   int       `json:"TeamA"`
	Scorea  int       `json:"ScoreA"`
	Pena    int       `json:"PenA"`
	Teamb   int       `json:"TeamB"`
	Scoreb  int       `json:"ScoreB"`
	Penb    int       `json:"PenB"`
	Stadium int       `json:"Stadium"`
}
type Pronostic struct {
	ID     int `json:"PronoID"`
	User   int `json:"User"`
	Match  int `json:"Match"`
	Scorea int `json:"ScoreA"`
	Scoreb int `json:"ScoreB"`
	Pena   int `json:"PenA"`
	Penb   int `json:"PenB"`
}

//Stadium
type Stadium struct {
	ID   int    `json:"ID"`
	Name string `json:"Name"`
}

//Stage
type Stage struct {
	ID     int    `json:"ID"`
	Name   string `json:"Name"`
	Active int    `json:"Active"`
}

//Stats
//Statistics
type Statistics struct {
	ButProno int    `json:"ButProno"`
	ButReal  int    `json:"ButReal"`
	Version  string `json:"Version"`
}
type Stats struct {
	Nombrebut  int `json:"Nombrebut"`
	Goodprono  int `json:"Goodprono"`
	Goodwinner int `json:"Goodwinner"`
	Redcard    int `json:"Redcard"`
	Yellowcard int `json:"Yellowcard"`
}

//Config
type Config struct {
	Lastmatch int `json:"LastMatch"`
	Refresh   int `json:"Refresh"`
	Stage     int `json:"Stage"`
}

//Player
type Player struct {
	ID          int       `json:"ID"`
	Name        string    `json:"Name"`
	Email       string    `json:"Email"`
	Score       int       `json:"Score"`
	Matches     []PrMatch `json:"Matches"`
	Champ       string    `json:"Champ"`
	Status      string    `json:"Status"`
	Amount      int       `json:"Amount"`
	ChangeChamp int       `json:"ChangeChamp"`
	Rank        int       `json:"Rank"`
	BonusMalus  int       `json:"BonusMalus"`
	ChampActive bool      `json:"ChampActive"`
}
type Users struct {
	ID             int    `json:"ID"`
	Firstname      string `json:"Firstname"`
	Lastname       string `json:"Lastname"`
	Score          int    `json:"Score"`
	Champion       int    `json:"Champion"`
	Champchange    int    `json:"Champchange"`
	Position       int    `json:"Position"`
	Positionbefore int    `json:"Positionbefore"`
	Malus          int    `json:"Malus"`
	Bonus          int    `json:"Bonus"`
	Ngoodscores    int    `json:"Ngoodscores"`
	Ngoodwinner    int    `json:"Ngoodwinner"`
}
