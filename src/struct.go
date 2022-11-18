package main

import "time"

// Pronostic
type Pronostic struct {
	MatchID int       `json:"MatchID"`
	Team1   string    `json:"Team1"`
	ScoreT1 int       `json:"ScoreT1"`
	ScoreP1 int       `json:"ScoreP1"`
	Team2   string    `json:"Team2"`
	ScoreT2 int       `json:"ScoreT2"`
	ScoreP2 int       `json:"ScoreP2"`
	Date    time.Time `json:"Date"`
	Stage   string    `json:"Stage"`
	Winner  string    `json:"Winner"`
	ScoreP  int       `json:"ScoreP"`
	Done    bool      `json:"Done"`
}

type ShortProno struct {
	UserID  int  `json:"UserID"`
	MatchID int  `json:"MatchID"`
	ScoreT1 int  `json:"ScoreT1"`
	ScoreP1 int  `json:"ScoreP1"`
	ScoreT2 int  `json:"ScoreT2"`
	ScoreP2 int  `json:"ScoreP2"`
	Done    bool `json:"Done"`
}

//Match
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
	Done    bool      `json:"Done"`
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
	Flag        string `json:"Flag"`
}

//Champion
type Champion struct {
	ID       int    `json:"ID"`
	Name     string `json:"Name"`
	Selected bool   `json:"Selected"`
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
	ID             int         `json:"ID"`
	Matches        []Pronostic `json:"Matches"`
	Status         string      `json:"Status"`
	Email          string      `json:"Email"`
	Firstname      string      `json:"Firstname"`
	Lastname       string      `json:"Lastname"`
	Score          int         `json:"Score"`
	Champion       string      `json:"Champion"`
	Champchange    int         `json:"Champchange"`
	Rank           int         `json:"Rank"`
	Position       int         `json:"Position"`
	Positionbefore int         `json:"Positionbefore"`
	Malus          int         `json:"Malus"`
	Bonus          int         `json:"Bonus"`
	Ngoodscores    int         `json:"Ngoodscores"`
	Ngoodwinner    int         `json:"Ngoodwinner"`
	ChampActive    bool        `json:"ChampActive"`
}

//loginDetail
type loginDetail struct {
	Firstname string `json:"Firstname"`
	Lastname  string `json:"Lastname"`
	Userid    int    `json:"Userid"`
	Status    string `json:"Status"`
}

//User
type User struct {
	Username      string      `json:"Username"`
	LoginDetail   loginDetail `json:"LoginDetail"`
	Authenticated bool        `json:"Authenticated"`
}
