package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func init() {
	root := mux.NewRouter()

	// root
	root.HandleFunc("/", indexHandler)

	// assets
	root.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("public/assets"))))

	// API endpoints
	root.HandleFunc("/api/v1/teams/", teamsHandler)                                          //Complete
	root.HandleFunc("/api/v1/teams/{leagueId}", teamsForFacilityHandler)                     //Complete
	root.HandleFunc("/api/v1/favorites/{team}", addFavoriteTeamHandler).Methods("POST")      //Complete
	root.HandleFunc("/api/v1/favorites/{team}", removeFavoriteTeamHandler).Methods("DELETE") //Complete
	root.HandleFunc("/api/v1/favorites", favoriteTeamsHandler)                               //Complete
	root.HandleFunc("/api/v1/favorites/games/", favoriteTeamsGamesHandler)                   //Complete
	root.HandleFunc("/api/v1/todaysGames/{league}", todaysGamesHandler)                      //Complete
	root.HandleFunc("/api/v1/tomorrowGames/{league}", tomorrowGamesHandler)                  //Complete
	root.HandleFunc("/api/v1/games/{team}", gamesForTeamHandler)                             //Complete
	root.HandleFunc("/api/v1/divisions/{division}/games", divisionGamesHandler)              //Complete
	root.HandleFunc("/api/v1/facilitys/{league}/divisions", facilityDivisionsHandler)        //Complete
	root.HandleFunc("/api/v1/divisions/{division}/teams", divisionsTeamsHandler)             //Complete
	root.HandleFunc("/api/v1/videoUpload", videoUploadHandler)

	http.Handle("/", root)
}

func main() {
	http.ListenAndServe(":8960", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t := lookupTemplate("index")
	t.Execute(w, nil)
}
