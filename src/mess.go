package main

// This is dirty. Must do a change in the DB but DB is ~life so this will be later. the faster the better

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
