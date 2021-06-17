package main

//Player
type Player struct {
	ID      int       `json:"ID"`
	Name    string    `json:"Name"`
	Email   string    `json:"Email"`
	Score   int       `json:"Score"`
	Matches []PrMatch `json:"Matches"`
	Champ   string    `json:"Champ"`
	Status  string    `json:"Status"`
	Amount  int       `json:"Amount"`
	Rank    int       `json:"Rank"`
}

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

//Config
type Config struct {
	LastMatchID int  `json:"LastMatchID"`
	Refresh     bool `json:"Refresh"`
}
