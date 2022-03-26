package main

import (
	"database/sql"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type M map[string]interface{}

func getOfficialMatches(w http.ResponseWriter, r *http.Request) {
	// Create the database handle, confirm driver is present
	db, _ := sql.Open("mysql", "root:@tcp(lilnas:3306)/footprono")
	defer db.Close()
	rows, _ := db.Query("SELECT m.id, st.name as stage, date, t1.name as teama, t2.name as teamb, s.name as stadium FROM matches m LEFT JOIN teams t1 on m.teama  = t1.id  LEFT JOIN teams t2 on m.teamb = t2.id LEFT JOIN stadium s on m.stadium = s.id  LEFT JOIN stage st on m.stage = st.id LIMIT 0, 1000")
	var amams []Mama
	for rows.Next() {
		var amam Mama
		rows.Scan(&amam.ID, &amam.Stage, &amam.Date, &amam.Teama, &amam.Teamb, &amam.Stadium)
		amams = append(amams, amam)
	}
	tmpl := template.Must(template.ParseFiles("templates/matches.gohtml"))
	tmpl.Execute(w, M{
		// We can pass as many things as we like
		"amams": amams,
		"stat":  stat,
	})
}

func getTeams(w http.ResponseWriter, r *http.Request) {
	// Create the database handle, confirm driver is present
	db, _ := sql.Open("mysql", "root:@tcp(lilnas:3306)/footprono")
	defer db.Close()
	rows, _ := db.Query("SELECT id, name, active, groupid, point, win, drawn, lose, goalfor, goalagainst FROM footprono.teams order by id;")
	var teams []Team
	for rows.Next() {
		var team Team
		rows.Scan(&team.ID, &team.Name, &team.Active, &team.Groupid, &team.Point, &team.Win, &team.Drawn, &team.Lose, &team.Goalfor, &team.Goalagainst)
		teams = append(teams, team)
	}
	tmpl := template.Must(template.ParseFiles("templates/teams.gohtml"))
	tmpl.Execute(w, M{
		// We can pass as many things as we like
		"teams": teams,
		"stat":  stat,
	})
}
