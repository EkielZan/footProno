package main

import "time"

// PrMatch
type Pronostics struct {
	MatchID int       `json:"MatchID"`
	Team1   string    `json:"Team1"`
	ScoreT1 int       `json:"ScoreT1"`
	ScoreP1 int       `json:"ScoreP1"`
	Team2   string    `json:"Team2"`
	ScoreT2 int       `json:"ScoreT2"`
	ScoreP2 int       `json:"ScoreP2"`
	Date    time.Time `json:"Date"`
	Stage   int       `json:"Stage"`
	Winner  string    `json:"Winner"`
	ScoreP  int       `json:"ScoreP"`
	Done    bool      `json:"Done"`
}

type Match struct {
	ID      int       `json:"ID"`
	Stage   string    `json:"Stage"`
	Date    time.Time `json:"Date"`
	Teama   string    `json:"TeamA"`
	Scorea  int       `json:"ScoreA"`
	Pena    int       `json:"PenA"`
	Teamb   string    `json:"TeamB"`
	Scoreb  int       `json:"ScoreB"`
	Penb    int       `json:"PenB"`
	Stadium string    `json:"Stadium"`
}

//Match
type Mama struct {
	ID      int       `json:"ID"`
	Stage   string    `json:"Stage"`
	Date    time.Time `json:"Date"`
	Teama   string    `json:"TeamA"`
	Teamb   string    `json:"TeamB"`
	Stadium string    `json:"Stadium"`
}

//Team
type Team struct {
	ID          int    `json:"ID"`
	Name        string `json:"Name"`
	Active      int    `json:"Active"`
	Groupid     string `json:"Group"`
	Point       int    `json:"Point"`
	Win         int    `json:"Win"`
	Drawn       int    `json:"Draw"`
	Lose        int    `json:"Lose"`
	Goalfor     int    `json:"GFor"`
	Goalagainst int    `json:"GAgainst"`
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
type Stats struct {
	Nombrebut  int    `json:"Nombrebut"`
	Goodprono  int    `json:"Goodprono"`
	Goodwinner int    `json:"Goodwinner"`
	Redcard    int    `json:"Redcard"`
	Yellowcard int    `json:"Yellowcard"`
	Version    string `json:"Version"`
}

//Config
type Config struct {
	Lastmatch int `json:"LastMatch"`
	Refresh   int `json:"Refresh"`
	Stage     int `json:"Stage"`
}

//Player
type Player struct {
	ID             int          `json:"ID"`
	Matches        []Pronostics `json:"Matches"`
	Status         string       `json:"Status"`
	Email          string       `json:"Email"`
	Firstname      string       `json:"Firstname"`
	Lastname       string       `json:"Lastname"`
	Score          int          `json:"Score"`
	Champion       int          `json:"Champion"`
	Champchange    int          `json:"Champchange"`
	Rank           int          `json:"Rank"`
	Position       int          `json:"Position"`
	Positionbefore int          `json:"Positionbefore"`
	Malus          int          `json:"Malus"`
	Bonus          int          `json:"Bonus"`
	Ngoodscores    int          `json:"Ngoodscores"`
	Ngoodwinner    int          `json:"Ngoodwinner"`
	ChampActive    bool         `json:"ChampActive"`
}

type loginDetail struct {
	Password  string
	Firstname string
	Lastname  string
	Userid    int
}
type User struct {
	Username      string
	Authenticated bool
}
