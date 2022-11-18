package main

import (
	"time"
)

// This is dirty. Must do a change in the DB but DB is ~life so this will be later. the faster the better

func compareDate(matchDate time.Time) bool {
	//Probably useless and should be cleaner
	matchDate2 := matchDate.Format("02-01-2006 03:04:05")
	t, _ := time.Parse("02-01-2006 03:04:05", matchDate2)
	today := time.Now()
	today2 := today.Format("02-01-2006 03:04:05")
	t2, _ := time.Parse("02-01-2006 03:04:05", today2)

	// Using time.After() method
	g2 := t2.After(t)

	return g2
}

func getFlag(champion string) string {
	switch champion {
	case "Argentina":
		return "ar"
	case "Australia":
		return "au"
	case "Belgium":
		return "be"
	case "Brazil":
		return "br"
	case "Cameroon":
		return "cm"
	case "Canada":
		return "ca"
	case "Costa Rica":
		return "cr"
	case "Croatia":
		return "hr"
	case "Denmark":
		return "dk"
	case "Ecuador":
		return "ec"
	case "England":
		return "gb-eng"
	case "France":
		return "fr"
	case "Germany":
		return "de"
	case "Ghana":
		return "gh"
	case "Iran":
		return "ir"
	case "Japan":
		return "jp"
	case "Mexico":
		return "mx"
	case "Morocco":
		return "ma"
	case "Netherlands":
		return "nl"
	case "Poland":
		return "pl"
	case "Portugal":
		return "pt"
	case "Qatar":
		return "qa"
	case "Saudi Arabia":
		return "sa"
	case "Senegal":
		return "sn"
	case "Serbia":
		return "rs"
	case "South Korea":
		return "kr"
	case "Spain":
		return "es"
	case "Switzerland":
		return "ch"
	case "Tunisia":
		return "tn"
	case "United States":
		return "us"
	case "Uruguay":
		return "uy"
	case "Wales":
		return "gb-wls"
	}

	return ""
}
