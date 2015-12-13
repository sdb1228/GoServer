package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

type team struct {
	Name     string `json:"name"`
	Division string `json:"division"`
	Teamid   string `json:"teamid"`
}

type division struct {
	Division string `json:"division"`
}

type installationTeam struct {
	Name          string `json:"name"`
	Division      string `json:"division"`
	Teamid        string `json:"teamid"`
	InstalltionId string `json:"installationId"`
}

type game struct {
	Awayteam      string    `json:"awayteam"`
	Hometeam      string    `json:"hometeam"`
	Field         string    `json:"field"`
	Address       string    `json:"address"`
	Hometeamscore *int      `json:"hometeamscore"`
	Awayteamscore *int      `json:"awayteamscore"`
	Gamesdatetime time.Time `json:"gamesdatetime"`
}

func init() {
	var err error
	db, err = sql.Open("postgres", "user=dburnett dbname=Soccer_Games host=54.68.232.199 password=doug1 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	root := mux.NewRouter()
	root.HandleFunc("/api/v1/teams/", teamsHandler)                                          //Complete
	root.HandleFunc("/api/v1/teams/{leagueId}", teamsForFacilityHandler)                     //Complete
	root.HandleFunc("/api/v1/favorites/{team}", addFavoriteTeamHandler).Methods("POST")      //Complete
	root.HandleFunc("/api/v1/favorites/{team}", removeFavoriteTeamHandler).Methods("DELETE") //Complete
	root.HandleFunc("/api/v1/favorites", favoriteTeamsHandler)                               //Complete
	root.HandleFunc("/api/v1/favorites/games/", favoriteTeamsGamesHandler)                   //Complete
	root.HandleFunc("/api/v1/todaysGames/{league}", todaysGamesHandler)                      //Complete
	root.HandleFunc("/api/v1/tomorrowGames/{league}", tomorrowGamesHandler)                  //Complete
	root.HandleFunc("/api/v1/games/{team}", gamesForTeamHandler)                             //Complete
	root.HandleFunc("/api/v1/divisions/{division}/games", divisionGamesHandeler)             //Complete
	root.HandleFunc("/api/v1/facilitys/{league}/divisions", facilityDivisionsHandler)        //Complete
	root.HandleFunc("/api/v1/divisions/{division}/teams", divisionsTeamsHandler)             //Complete
	// root.HandleFunc("/api/v1/facilitys", facilityHandler)

	http.Handle("/", root)
}

func main() {
	http.ListenAndServe(":8960", nil)
}

/*
This method returns all the divisions of a specific facility/league sorted by division name
*/
func facilityDivisionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	league := vars["league"]
	var buffer bytes.Buffer
	buffer.WriteString("SELECT DISTINCT division FROM teams WHERE facility=")
	buffer.WriteString(league)
	buffer.WriteString(" ORDER BY division;")
	fmt.Println(buffer.String())
	rows, err := db.Query(buffer.String())
	if err != nil {
		log.Fatal(err)
	}

	results := []division{}
	for rows.Next() {
		var d division
		rows.Scan(&d.Division)
		results = append(results, d)
	}
	fmt.Printf("%v", results)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(&results)
}

/*
This method handles the returning of all the teams from a given facility ID
*/
func teamsForFacilityHandler(w http.ResponseWriter, r *http.Request) {
	installationId := r.FormValue("installationId")
	vars := mux.Vars(r)
	league := vars["leagueId"]
	var buffer bytes.Buffer
	fmt.Println(installationId)
	buffer.WriteString("SELECT name, division, teams.teamid, installationid FROM teams LEFT OUTER JOIN favorites f1 ON f1.installationid='")
	buffer.WriteString(installationId)
	buffer.WriteString("' AND f1.teamid=teams.teamid WHERE facility=")
	buffer.WriteString(league)
	buffer.WriteString(" ORDER BY name;")
	fmt.Println(buffer.String())
	rows, err := db.Query(buffer.String())
	if err != nil {
		log.Fatal(err)
	}

	results := []installationTeam{}
	for rows.Next() {
		var t installationTeam
		rows.Scan(&t.Name, &t.Division, &t.Teamid, &t.InstalltionId)
		results = append(results, t)
	}
	fmt.Printf("%v", results)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(&results)
}

/*
Get the games for a sepcific division between the start date and the end date formatted like yyyy-MM-dd hh:mm:ss ex: 2015-12-10 11:48:59
*/
func divisionGamesHandeler(w http.ResponseWriter, r *http.Request) {
	facility := r.FormValue("facility")
	vars := mux.Vars(r)
	division := vars["division"]
	startDate := r.FormValue("startDate")
	endDate := r.FormValue("endDate")
	var buffer bytes.Buffer

	buffer.WriteString("SELECT f1.name AS field, f1.address AS address, a2.name AS hometeam, a1.name AS awayteam, games.gamesdatetime, games.hometeamscore, games.awayteamscore ")
	buffer.WriteString("FROM games ")
	buffer.WriteString("INNER JOIN fields f1 ON f1.id=games.field ")
	buffer.WriteString("INNER JOIN teams a1 ON games.awayteam=a1.teamid ")
	buffer.WriteString("INNER JOIN teams a2 ON games.hometeam=a2.teamid ")
	buffer.WriteString("WHERE a1.facility=")
	buffer.WriteString(facility)
	buffer.WriteString(" AND a2.facility=")
	buffer.WriteString(facility)
	buffer.WriteString(" AND a1.division='")
	buffer.WriteString(division)
	buffer.WriteString("' AND a2.division='")
	buffer.WriteString(division)
	buffer.WriteString("' AND games.gamesdatetime >= '")
	buffer.WriteString(startDate)
	buffer.WriteString("' AND games.gamesdatetime <= '")
	buffer.WriteString(endDate)
	buffer.WriteString("' ORDER BY games.gamesdatetime;")
	fmt.Println(buffer.String())

	rows, err := db.Query(buffer.String())
	if err != nil {
		log.Fatal(err)
	}

	results := []game{}
	for rows.Next() {
		var g game
		rows.Scan(&g.Field, &g.Address, &g.Hometeam, &g.Awayteam, &g.Gamesdatetime, &g.Hometeamscore, &g.Awayteamscore)
		results = append(results, g)
	}
	fmt.Printf("%v", results)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(&results)

}

/*
Get the teams for a specific facility in a specific division
*/
func divisionsTeamsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	facility := r.FormValue("facility")
	division := vars["division"]
	var buffer bytes.Buffer

	buffer.WriteString("SELECT name, division, teamid FROM teams WHERE facility=")
	buffer.WriteString(facility)
	buffer.WriteString(" AND division='")
	buffer.WriteString(division)
	buffer.WriteString("';")
	rows, err := db.Query(buffer.String())
	if err != nil {
		log.Fatal(err)
	}
	results := []team{}
	for rows.Next() {
		var t team
		rows.Scan(&t.Name, &t.Division, &t.Teamid)
		results = append(results, t)
	}
	fmt.Printf("%v", results)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(&results)

}

/*
This method handles the returning of all the teams
*/
func teamsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT name, division, teamid FROM teams ORDER BY name;")
	if err != nil {
		log.Fatal(err)
	}

	results := []team{}
	for rows.Next() {
		var t team
		rows.Scan(&t.Name, &t.Division, &t.Teamid)
		results = append(results, t)
	}
	fmt.Printf("%v", results)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(&results)
}

/*
Returns all the teams that a certain device has favorited
*/
func favoriteTeamsHandler(w http.ResponseWriter, r *http.Request) {
	installationId := r.FormValue("installationId")
	checkInstallation(installationId)
	fmt.Println(installationId)
	var buffer bytes.Buffer
	buffer.WriteString("SELECT a1.name, a1.division, favorites.teamid FROM favorites INNER JOIN teams a1 ON favorites.teamid=a1.teamid WHERE installationid='")
	buffer.WriteString(installationId)
	buffer.WriteString("';")
	rows, err := db.Query(buffer.String())
	if err != nil {
		log.Fatal(err)
	}
	results := []team{}
	for rows.Next() {
		var t team
		rows.Scan(&t.Name, &t.Division, &t.Teamid)
		results = append(results, t)
	}
	fmt.Printf("%v", results)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(&results)
}

/*
Returns all the games of an installations favorite teams including address of the field and all games that don't have a score yet.
*/

func favoriteTeamsGamesHandler(w http.ResponseWriter, r *http.Request) {
	installationId := r.FormValue("installationId")
	checkInstallation(installationId)
	var buffer bytes.Buffer

	buffer.WriteString("SELECT f1.name AS field, f1.address AS address, a2.name AS hometeam, a1.name AS awayteam, games.gamesdatetime, games.hometeamscore, games.awayteamscore ")
	buffer.WriteString("FROM favorites, games ")
	buffer.WriteString("LEFT OUTER JOIN fields f1 ON f1.id=games.field ")
	buffer.WriteString("LEFT OUTER JOIN teams a1 ON games.awayteam=a1.teamid ")
	buffer.WriteString("LEFT OUTER JOIN teams a2 ON games.hometeam=a2.teamid ")
	buffer.WriteString("WHERE favorites.installationid='")
	buffer.WriteString(installationId)
	buffer.WriteString("' AND (games.hometeam=favorites.teamid OR games.awayteam=favorites.teamid) AND games.gamesdatetime >= now()::timestamp ORDER BY games.gamesdatetime;")
	rows, err := db.Query(buffer.String())
	if err != nil {
		log.Fatal(err)
	}

	results := []game{}
	for rows.Next() {
		var g game
		rows.Scan(&g.Field, &g.Address, &g.Hometeam, &g.Awayteam, &g.Gamesdatetime, &g.Hometeamscore, &g.Awayteamscore)
		results = append(results, g)
	}
	fmt.Printf("%v", results)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(&results)

}

/*
Returns the games of a specific facility for today
*/
func todaysGamesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Failed today")
	vars := mux.Vars(r)
	league := vars["league"]
	var buffer bytes.Buffer
	t := time.Now()
	t.Format("2006-01-02")
	stringDate := t.String()

	parsedDate := stringDate[0:10]
	buffer.WriteString("SELECT f1.name AS field, f1.address AS address, a2.name AS hometeam, a1.name AS awayteam, games.gamesdatetime, games.hometeamscore, games.awayteamscore ")
	buffer.WriteString("FROM games ")
	buffer.WriteString("INNER JOIN fields f1 ON f1.id=games.field ")
	buffer.WriteString("INNER JOIN teams a1 ON games.awayteam=a1.teamid ")
	buffer.WriteString("INNER JOIN teams a2 ON games.hometeam=a2.teamid ")
	buffer.WriteString("WHERE gamesdatetime::text LIKE '")
	buffer.WriteString(parsedDate)
	buffer.WriteString("%' AND a1.facility=")
	buffer.WriteString(league)
	buffer.WriteString(" AND a2.facility=")
	buffer.WriteString(league)
	buffer.WriteString(" ORDER BY games.gamesdatetime;")

	fmt.Println(buffer.String())

	rows, err := db.Query(buffer.String())
	if err != nil {
		log.Fatal(err)
	}

	results := []game{}
	for rows.Next() {
		var g game
		rows.Scan(&g.Field, &g.Address, &g.Hometeam, &g.Awayteam, &g.Gamesdatetime, &g.Hometeamscore, &g.Awayteamscore)
		results = append(results, g)
	}
	fmt.Printf("%v", results)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(&results)
}

/*
Removes an installations favorite team
*/
func removeFavoriteTeamHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamId := vars["team"]
	installationId := r.FormValue("installationId")
	checkInstallation(installationId)
	_, err := db.Exec(
		"DELETE FROM favorites WHERE installationid=$1 AND teamid=$2;",
		installationId,
		teamId,
	)
	if err != nil {
		log.Fatal(err)
	}
	var results [0]string
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(&results)

}

/*
Adds a favorite team to an installation
*/
func addFavoriteTeamHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	teamId := vars["team"]
	installationId := r.FormValue("installationId")
	checkInstallation(installationId)
	var buffer bytes.Buffer
	buffer.WriteString("SELECT COUNT(*) FROM favorites WHERE installationid='")
	buffer.WriteString(installationId)
	buffer.WriteString("' AND teamid='")
	buffer.WriteString(teamId)
	buffer.WriteString("';")
	fmt.Println(buffer.String())
	rows, err := db.Query(buffer.String())

	if err != nil {
		log.Fatal(err)
	}
	var count int
	for rows.Next() {
		rows.Scan(&count)
	}
	if count != 0 {
		return
	}

	_, err = db.Exec(
		"INSERT INTO favorites (installationid, teamid) VALUES ($1, $2);",
		installationId,
		teamId,
	)
	fmt.Println("Inserting favorites")
	if err != nil {
		log.Fatal(err)
	}
	var results [0]string
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(&results)

}

/*
Returns all the games for tomorrow for a specific facility
*/
func tomorrowGamesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Failed Tomorrow")
	vars := mux.Vars(r)
	league := vars["league"]
	var buffer bytes.Buffer
	t := time.Now()
	tomorrowDate := t.AddDate(0, 0, 1)
	tomorrowDate.Format("2006-01-02")
	stringDate := tomorrowDate.String()

	parsedDate := stringDate[0:10]
	buffer.WriteString("SELECT f1.name AS field, f1.address AS address, a2.name AS hometeam, a1.name AS awayteam, games.gamesdatetime, games.hometeamscore, games.awayteamscore ")
	buffer.WriteString("FROM games ")
	buffer.WriteString("INNER JOIN fields f1 ON f1.id=games.field ")
	buffer.WriteString("INNER JOIN teams a1 ON games.awayteam=a1.teamid ")
	buffer.WriteString("INNER JOIN teams a2 ON games.hometeam=a2.teamid ")
	buffer.WriteString("WHERE gamesdatetime::text LIKE '")
	buffer.WriteString(parsedDate)
	buffer.WriteString("%' AND a1.facility=")
	buffer.WriteString(league)
	buffer.WriteString(" AND a2.facility=")
	buffer.WriteString(league)
	buffer.WriteString(" ORDER BY games.gamesdatetime;")

	fmt.Println(buffer.String())

	rows, err := db.Query(buffer.String())
	if err != nil {
		log.Fatal(err)
	}

	results := []game{}
	for rows.Next() {
		var g game
		rows.Scan(&g.Field, &g.Address, &g.Hometeam, &g.Awayteam, &g.Gamesdatetime, &g.Hometeamscore, &g.Awayteamscore)
		results = append(results, g)
	}
	fmt.Printf("%v", results)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(&results)
}

/*
Returns all the games fora  specific team
*/
func gamesForTeamHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	team := vars["team"]
	var buffer bytes.Buffer

	buffer.WriteString("SELECT f1.name AS field, f1.address AS address, a2.name AS hometeam, a1.name AS awayteam, games.gamesdatetime, games.hometeamscore, games.awayteamscore ")
	buffer.WriteString("FROM games ")
	buffer.WriteString("INNER JOIN fields f1 ON f1.id=games.field ")
	buffer.WriteString("INNER JOIN teams a1 ON games.awayteam=a1.teamid ")
	buffer.WriteString("INNER JOIN teams a2 ON games.hometeam=a2.teamid ")
	buffer.WriteString("WHERE games.awayteam='")
	buffer.WriteString(team)
	buffer.WriteString("' OR games.hometeam='")
	buffer.WriteString(team)
	buffer.WriteString("' ORDER BY games.gamesdatetime;")

	fmt.Println(buffer.String())

	rows, err := db.Query(buffer.String())
	if err != nil {
		log.Fatal(err)
	}

	results := []game{}
	for rows.Next() {
		var g game
		rows.Scan(&g.Field, &g.Address, &g.Hometeam, &g.Awayteam, &g.Gamesdatetime, &g.Hometeamscore, &g.Awayteamscore)
		results = append(results, g)
	}
	fmt.Printf("%v", results)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(&results)
}

/*
Checks to see if we have the current installation in the database.  If we don't we will insert it
*/

func checkInstallation(installationId string) {
	var buffer bytes.Buffer
	buffer.WriteString("SELECT COUNT(*) FROM installation where installationid='")
	buffer.WriteString(installationId)
	buffer.WriteString("';")
	rows, err := db.Query(buffer.String())
	if err != nil {
		log.Fatal(err)
	}

	var count int
	for rows.Next() {
		rows.Scan(&count)
	}
	if count == 0 {
		_, err = db.Exec(
			"INSERT INTO installation (installationid) VALUES ($1)",
			installationId,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

}
