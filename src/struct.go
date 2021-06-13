package main

//Player
type Player struct {
	ID        int            `json:"ID"`
	Name      string         `json:"Name"`
	Email     string         `json:"Email"`
	Score     int            `json:"Score"`
	Matches   []PrMatch      `json:"Matches"`
	Champ     string         `json:"Champ"`
	Positions []LastPosition `json:"Position"`
	Status    string         `json:"Status"`
	Amount    int            `json:"Amount"`
	Rank      int            `json:"Rank"`
}

type ShortPlayer struct {
	Name      string         `json:"Name"`
	Score     int            `json:"Score"`
	Positions []LastPosition `json:"Position"`
}

type LastPosition struct {
	Position  int    `json:"LastPos"`
	ScoreDate string `json:"ScoreDate"`
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

//Config
type Config struct {
	LastSave      string `json:"LastSave"`
	LastMatchDate string `json:"LastMatchDate"`
	LastMatchID   int    `json:"LastMatchID"`
}
