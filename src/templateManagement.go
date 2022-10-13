// index serves the index html file
package main

//
//   file:test.db?cache=shared&mode=memory
//

import (
	"database/sql"
	"net/http"

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
			"user":    user,
			"stat":    stat,
			"flashes": flashes,
		})
	if err != nil {
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
	if !checkAuthentication(w, r, user, session) {
		return
	}
	// Pass Struct and execute template for display
	tpl.ExecuteTemplate(w,
		"about.gohtml",
		M{
			"user": user,
			"stat": stat,
		})
}
