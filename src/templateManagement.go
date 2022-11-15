// index serves the index html file
package main

//
//   file:test.db?cache=shared&mode=memory
//

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	//	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

type M map[string]interface{}

func index(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := getUser(session)
	flashes := session.Flashes()
	session.Flashes()
	err = session.Save(r, w)
	tpl.ExecuteTemplate(w,
		"index.gohtml",
		M{
			// We can pass as many things as we like
			"user":  user,
			"stat":  stat,
			"flash": flashes,
		})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func registerForm(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := getUser(session)
	flashes := session.Flashes()
	//flashes :=
	session.Flashes()
	err = session.Save(r, w)
	tpl.ExecuteTemplate(w,
		"register.gohtml",
		M{
			// We can pass as many things as we like
			"user":  user,
			"stat":  stat,
			"flash": flashes,
		})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getOfficialMatches(w http.ResponseWriter, r *http.Request) {
	// Manage Sessions and authentication
	session, err := store.Get(r, cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := getUser(session)
	if !checkAuthentication(w, r, user, session) {
		return
	}

	// Create the database handle, confirm driver is present
	//FORNOSQL db, _ := sql.Open("mysql", "root:@tcp(lilnas:3306)/footprono?parseTime=true")
	db, _ := sql.Open(SDRIVER, SCON)

	defer db.Close()
	//`id`,`stage`,`date`,`teama`,`scorea`,`pena`,`teamb`,`scoreb`,`penb`,`stadium` -> Matches
	rows, _ := db.Query("SELECT m.id, st.name as stage, date, t1.name as teama, m.scorea, m.pena, t2.name as teamb, m.scoreb, m.penb, s.name as stadium FROM matches m LEFT JOIN teams t1 on m.teama  = t1.id  LEFT JOIN teams t2 on m.teamb = t2.id LEFT JOIN stadium s on m.stadium = s.id  LEFT JOIN stage st on m.stage = st.id LIMIT 0, 1000")
	var matches []Match
	for rows.Next() {
		var match Match
		rows.Scan(&match.ID, &match.Stage, &match.Date, &match.Teama, &match.Scorea, &match.Pena, &match.Teamb, &match.Scoreb, &match.Penb, &match.Stadium)
		matches = append(matches, match)
	}
	// Pass Struct and execute template for display
	tpl.ExecuteTemplate(w,
		"matches.gohtml",
		M{
			"matches": matches,
			"user":    user,
			"stat":    stat,
		})
}

func getPronostics(w http.ResponseWriter, r *http.Request) {
	// Manage Sessions and authentication
	session, err := store.Get(r, cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := getUser(session)
	if !checkAuthentication(w, r, user, session) {
		return
	}
	// Create the database handle, confirm driver is present
	//FORNOSQL db, _ := sql.Open("mysql", "root:@tcp(lilnas:3306)/footprono?parseTime=true")
	db, _ := sql.Open(SDRIVER, SCON)
	rows0, _ := db.Query("SELECT lastmatch, refresh, stage from config;")
	var config Config
	for rows0.Next() {
		rows0.Scan(&config.Lastmatch, &config.Refresh, &config.Stage)
	}

	rows, _ := db.Query("SELECT m.id, st.name as stage, date, t1.name as teama, m.scorea, m.pena, t2.name as teamb, m.scoreb, m.penb, s.name as stadium FROM matches m LEFT JOIN teams t1 on m.teama  = t1.id  LEFT JOIN teams t2 on m.teamb = t2.id LEFT JOIN stadium s on m.stadium = s.id  LEFT JOIN stage st on m.stage = st.id where m.stage=? LIMIT 0, 1000", config.Stage)
	var matches []Match
	for rows.Next() {
		var match Match
		rows.Scan(&match.ID, &match.Stage, &match.Date, &match.Teama, &match.Scorea, &match.Pena, &match.Teamb, &match.Scoreb, &match.Penb, &match.Stadium)
		matches = append(matches, match)
	}
	defer db.Close()

	wu := user.LoginDetail.Userid

	var prons []Pronostic
	for _, m := range matches {
		Prows, _ := db.Query("SELECT p.match, st.name, m.date,  t1.name as teama, p.scorea, p.pena, t2.name as teamb, p.scoreb, p.penb FROM pronostics p LEFT JOIN matches m on m.id = p.match LEFT JOIN stage st on st.id = m.stage LEFT JOIN teams t1 on m.teama = t1.id LEFT JOIN teams t2 on m.teamb = t2.id where m.id = ? and p.user = ? LIMIT 0, 1000", m.ID, wu)
		var pron Pronostic
		for Prows.Next() {
			Prows.Scan(&pron.MatchID, &pron.Stage, &pron.Date, &pron.Team1, &pron.ScoreT1, &pron.ScoreP1, &pron.Team2, &pron.ScoreT2, &pron.ScoreP2)
		}
		if pron.MatchID == 0 {
			pron.MatchID = m.ID
			pron.Date = m.Date
			pron.Stage = m.Stage
			pron.Team1 = m.Teama
			pron.Team2 = m.Teamb
			pron.ScoreT1 = 0
			pron.ScoreT2 = 0
			pron.ScoreP1 = 0
			pron.ScoreP2 = 0
		}
		prons = append(prons, pron)
	}
	var dProns []Pronostic
	Prows, _ := db.Query("SELECT p.match, st.name, m.date,  t1.name as teama, p.scorea, p.pena, t2.name as teamb, p.scoreb, p.penb FROM pronostics p LEFT JOIN matches m on m.id = p.match LEFT JOIN stage st on st.id = m.stage LEFT JOIN teams t1 on m.teama = t1.id LEFT JOIN teams t2 on m.teamb = t2.id where st.id < ? and p.user = ? LIMIT 0, 1000", config.Stage, wu)
	for Prows.Next() {
		var dPron Pronostic
		Prows.Scan(&dPron.MatchID, &dPron.Stage, &dPron.Date, &dPron.Team1, &dPron.ScoreT1, &dPron.ScoreP1, &dPron.Team2, &dPron.ScoreT2, &dPron.ScoreP2)
		dProns = append(dProns, dPron)
	}
	// Pass Struct and execute template for display
	err = tpl.ExecuteTemplate(w,
		"prono.gohtml",
		M{
			"prons":  prons,
			"user":   user,
			"dProns": dProns,
			"stat":   stat,
		})
	if err != nil {
		log.Println(err)
	}
}

func getTeams(w http.ResponseWriter, r *http.Request) {
	// Manage Sessions and authentication
	session, err := store.Get(r, cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := getUser(session)
	if !checkAuthentication(w, r, user, session) {
		return
	}

	// Create the database handle, confirm driver is present
	//FORNOSQL db, _ := sql.Open("mysql", "root:@tcp(lilnas:3306)/footprono?parseTime=true")
	db, _ := sql.Open(SDRIVER, SCON)
	defer db.Close()
	rows, _ := db.Query("SELECT id, name, active, groupid, point, win, drawn, lose, goalfor, goalagainst FROM teams order by id;")
	var teams []Team
	for rows.Next() {
		var team Team
		rows.Scan(&team.ID, &team.Name, &team.Active, &team.Groupid, &team.Point, &team.Win, &team.Drawn, &team.Lose, &team.Goalfor, &team.Goalagainst)
		teams = append(teams, team)
	}
	// Pass Struct and execute template for display
	tpl.ExecuteTemplate(w,
		"teams.gohtml",
		M{
			"teams": teams,
			"user":  user,
			"stat":  stat,
		})
}

func addPronos(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := getUser(session)
	if !checkAuthentication(w, r, user, session) {
		return
	}
	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var npron ShortProno

	for k, _ := range r.PostForm {
		if strings.HasPrefix(k, "match-") {
			npron.MatchID, _ = strconv.Atoi(r.FormValue(k))
		}
		if strings.HasPrefix(k, "ScoreT1-") {
			npron.ScoreT1, _ = strconv.Atoi(r.FormValue(k))
		}
		if strings.HasPrefix(k, "ScoreT2-") {
			npron.ScoreT2, _ = strconv.Atoi(r.FormValue(k))
		}
		npron.UserID = user.LoginDetail.Userid
		npron.Done = true
	}
	rowid := strconv.Itoa(npron.UserID) + strconv.Itoa(npron.MatchID)
	db, _ := sql.Open(SDRIVER, SCON)
	defer db.Close()
	stm, _ := db.Prepare("INSERT or REPLACE into pronostics (id,scorea,scoreb,user,match) VALUES (?,?,?,?,?);")
	if err != nil {
		log.Println(err.Error())
		return
	}
	_ = stm.QueryRow(rowid, strconv.Itoa(npron.ScoreT1), strconv.Itoa(npron.ScoreT2), strconv.Itoa(npron.UserID), strconv.Itoa(npron.MatchID)).Scan()

	if err != nil {
		log.Println(err.Error())
		return
	}
	http.Redirect(w, r, "/prons", http.StatusFound)
}

func doRegister(w http.ResponseWriter, r *http.Request) {
	// Manage Sessions and authentication
	session, err := store.Get(r, cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := getUser(session)
	if !checkAuthentication(w, r, user, session) {
		return
	}
	tpl.ExecuteTemplate(w,
		"register.gohtml",
		M{
			"user": user,
			"stat": stat,
		})

}

func about(w http.ResponseWriter, r *http.Request) {
	// Manage Sessions and authentication
	session, err := store.Get(r, cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := getUser(session)
	// Pass Struct and execute template for display
	tpl.ExecuteTemplate(w,
		"about.gohtml",
		M{
			"user": user,
			"stat": stat,
		})
}
