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
Get the teams for a specific facility in a specific division
*/
func divisionsTeamsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	facility := r.FormValue("facility")
	division := vars["division"]
	var buffer bytes.Buffer

	if facility == "" || division == "" {
		log.Println("They are missing parameters in Division Teams")
		encoder.Encode(response_builder(422, "You are missing parameters"))
		return
	}

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
	encoder.Encode(&results)

}

/*
Get the games for a specific division between the start date and the end date formatted like yyyy-MM-dd hh:mm:ss ex: 2015-12-10 11:48:59
*/

func divisionGamesHandler(w http.ResponseWriter, r *http.Request) {
	facility := r.FormValue("facility")
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	division := vars["division"]
	startDate := r.FormValue("startDate")
	endDate := r.FormValue("endDate")
	if facility == "" || division == "" || startDate == "" || endDate == "" {
		log.Println("They are missing paramaters in Division Games")
		encoder.Encode(response_builder(422, "You are missing parameters"))
		return
	}

	rows, err := db.Query("SELECT f1.name AS field, f1.address AS address, a2.name AS hometeam, a1.name AS awayteam, games.gamesdatetime, games.hometeamscore, games.awayteamscore FROM games "+
		"INNER JOIN fields f1 ON f1.id=games.field "+
		"INNER JOIN teams a1 ON games.awayteam=a1.teamid "+
		"INNER JOIN teams a2 ON games.hometeam=a2.teamid "+
		"WHERE a1.facility=$1 "+
		"AND a2.facility=$1 "+
		"AND a1.division=$2 "+
		"AND a2.division=$2 "+
		"AND games.gamesdatetime >= $3 "+
		"AND games.gamesdatetime <= $4 "+
		"ORDER BY games.gamesdatetime", facility, division, startDate, endDate)
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
	encoder.Encode(&results)

}
