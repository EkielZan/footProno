package main

//Player
type Player struct {
	ID      int       `json:"ID"`
	Name    string    `json:"Name"`
	Email   string    `json:"Email"`
	Score   int       `json:"Score"`
	Matches []PrMatch `json:"Matches"`
	Champ   string    `json:"Champ"`
	LastPos int       `json:"LastPos"`
	Status  string    `json:"Status"`
	Amount  int       `json:"Amount"`
}

type ShortPlayer struct {
	Name    string `json:"Name"`
	Score   int    `json:"Score"`
	LastPos int    `json:"LastPos"`
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

//Statistics
type Statistics struct {
	Rise     string `json:"Rise"`
	Fall     string `json:"Fall"`
	ButProno int    `json:"ButProno"`
	ButReal  int    `json:"ButReal"`
}

//spaHandler
type spaHandler struct {
	staticPath string
	indexPath  string
}
