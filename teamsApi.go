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
This method handles the returning of all the teams from a given facility ID and returns the installation id if the team has been favorited by installationID
*/
func teamsForFacilityHandler(w http.ResponseWriter, r *http.Request) {
	installationId := r.FormValue("installationId")
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	// token := vars["token"]
	// err := credential_check(token)
	// if err != nil {
	// 	encoder.Encode(response_builder(403, "You don't have access to obtain that information"))
	// 	return
	// }
	err := checkInstallation(installationId)
	if err != nil {
		log.Println("Error in querying for installation ID: ", err)
		encoder.Encode(response_builder(403, "There was a problem with your installationID"))
		return
	}
	league := vars["leagueId"]
	err = verifyInteger(league)
	if err != nil {
		log.Println("Error league is not an integer : ", err)
		encoder.Encode(response_builder(403, "The league you have requested doesn't exist"))
		return
	}

	var buffer bytes.Buffer
	fmt.Println(installationId)
	buffer.WriteString("SELECT name, division, teams.teamid, installationid FROM teams LEFT OUTER JOIN favorites f1 ON f1.installationid='")
	buffer.WriteString(installationId)
	buffer.WriteString("' AND f1.teamid=teams.teamid WHERE facility=")
	buffer.WriteString(league)
	buffer.WriteString(" AND teams.deleted_at IS NULL ORDER BY name;")
	fmt.Println(buffer.String())
	rows, err := db.Query(buffer.String())
	if err != nil {
		log.Println("Error in DB query of teams for facility: ", err)
		encoder.Encode(response_builder(403, "Internal server error please try again later"))
		return
	}

	results := []installationTeam{}
	for rows.Next() {
		var t installationTeam
		rows.Scan(&t.Name, &t.Division, &t.Teamid, &t.InstalltionId)
		results = append(results, t)
	}
	fmt.Printf("%v", results)
	w.Header().Set("Content-Type", "application/json")
	encoder.Encode(&results)
}
