package main

//Player
type Player struct {
	ID      string    `json:"ID"`
	Name    string    `json:"Name"`
	Email   string    `json:"Email"`
	Score   int       `json:"Score"`
	Matches []PrMatch `json:"Matches"`
}

//Match from json files
type Match struct {
	MatchID int    `json:"MatchID"`
	Team1   string `json:"Team1"`
	ScoreT1 int    `json:"ScoreT1"`
	Team2   string `json:"Team2"`
	ScoreT2 int    `json:"ScoreT2"`
	Date    string `json:"Date"`
	Stage   string `json:"Stage"`
	Winner  string `json:"Winner"`
}

// PrMatch
type PrMatch struct {
	MatchID int    `json:"MatchID"`
	Team1   string `json:"Team1"`
	ScoreT1 int    `json:"ScoreT1"`
	Team2   string `json:"Team2"`
	ScoreT2 int    `json:"ScoreT2"`
	Date    string `json:"Date"`
	Stage   string `json:"Stage"`
	Winner  string `json:"Winner"`
	ScoreP  int    `json:"ScoreP"`
}

//spaHandler
type spaHandler struct {
	staticPath string
	indexPath  string
}
