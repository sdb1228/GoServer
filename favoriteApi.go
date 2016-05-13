package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

/*
Adds a favorite team to an installation
*/
func addFavoriteTeamHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	teamId := vars["team"]
	installationId := r.FormValue("installationId")
	err := checkInstallation(installationId)
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
This method returns all the divisions of a specific facility/league sorted by division name
*/
func facilityDivisionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoder := json.NewEncoder(w)
	// token := vars["token"]
	// err := credential_check(token)
	// if err != nil {
	// 	encoder.Encode(response_builder(403, "You don't have access to obtain that information"))
	// 	return
	// }
	league := vars["league"]
	err := verifyInteger(league)
	if err != nil {
		log.Println("Error league is not an integer : ", err)
		encoder.Encode(response_builder(403, "The league you have requested doesn't exist"))
		return
	}
	var buffer bytes.Buffer
	buffer.WriteString("SELECT DISTINCT division FROM teams WHERE facility=")
	buffer.WriteString(league)
	buffer.WriteString(" ORDER BY division;")
	fmt.Println(buffer.String())
	rows, err := db.Query(buffer.String())
	if err != nil {
		log.Println("Error in DB query of facility division: ", err)
		encoder.Encode(response_builder(403, "Internal server error please try again later"))
		return
	}
	results := []division{}
	for rows.Next() {
		var d division
		rows.Scan(&d.Division)
		results = append(results, d)
	}
	fmt.Printf("%v", results)
	w.Header().Set("Content-Type", "application/json")
	encoder.Encode(&results)
}

/*
Removes an installations favorite team
*/
func removeFavoriteTeamHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamId := vars["team"]
	installationId := r.FormValue("installationId")
	err := checkInstallation(installationId)
	_, err = db.Exec(
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
Returns all the teams that a certain device has favorited
*/
func favoriteTeamsHandler(w http.ResponseWriter, r *http.Request) {
	installationId := r.FormValue("installationId")
	err := checkInstallation(installationId)
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
	err := checkInstallation(installationId)
	var buffer bytes.Buffer
	buffer.WriteString("SELECT f1.name AS field, f1.address AS address, a2.name AS hometeam, a1.name AS awayteam, games.gamesdatetime, games.hometeamscore, games.awayteamscore ")
	buffer.WriteString("FROM favorites, games ")
	buffer.WriteString("LEFT OUTER JOIN fields f1 ON f1.id=games.field ")
	buffer.WriteString("LEFT OUTER JOIN teams a1 ON games.awayteam=a1.teamid ")
	buffer.WriteString("LEFT OUTER JOIN teams a2 ON games.hometeam=a2.teamid ")
	buffer.WriteString("WHERE favorites.installationid='")
	buffer.WriteString(installationId)
	buffer.WriteString("' AND (games.hometeam=favorites.teamid OR games.awayteam=favorites.teamid) AND games.gamesdatetime >= (now()::timestamp - '1 day'::interval) ORDER BY games.gamesdatetime;")

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
