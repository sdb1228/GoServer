package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type Page struct {
	URL  string
	Body string
}

func init() {
	root := mux.NewRouter()

	// root
	root.HandleFunc("/", indexHandler)

	// templates
	root.HandleFunc("/video/{video}", videoTemplateHandler)
	root.HandleFunc("/about", aboutHandler)
	root.HandleFunc("/contact", contactHandler)
	root.HandleFunc("/fields/update", fieldsHandler)

	// assets
	root.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("public/assets"))))

	// API endpoints
	root.HandleFunc("/api/v1/teams/", teamsHandler)                                               //Complete
	root.HandleFunc("/api/v1/fields/correction", fieldsCorrectionHandler)                         //Complete
	root.HandleFunc("/api/v1/fields/postCorrection", fieldsCorrectionPostHandler).Methods("POST") //Complete
	root.HandleFunc("/api/v1/teams/{leagueId}", teamsForFacilityHandler)                          //Complete
	root.HandleFunc("/api/v1/favorites/{team}", addFavoriteTeamHandler).Methods("POST")           //Complete
	root.HandleFunc("/api/v1/favorites/{team}", removeFavoriteTeamHandler).Methods("DELETE")      //Complete
	root.HandleFunc("/api/v1/favorites", favoriteTeamsHandler)                                    //Complete
	root.HandleFunc("/api/v1/favorites/games/", favoriteTeamsGamesHandler)                        //Complete
	root.HandleFunc("/api/v1/todaysGames/{league}", todaysGamesHandler)                           //Complete
	root.HandleFunc("/api/v1/tomorrowGames/{league}", tomorrowGamesHandler)                       //Complete
	root.HandleFunc("/api/v1/games/{team}", gamesForTeamHandler)                                  //Complete
	root.HandleFunc("/api/v1/divisions/{division}/games", divisionGamesHandler)                   //Complete
	root.HandleFunc("/api/v1/facilitys/{league}/divisions", facilityDivisionsHandler)             //Complete
	root.HandleFunc("/api/v1/divisions/{division}/teams", divisionsTeamsHandler)                  //Complete
	root.HandleFunc("/api/v1/videoUpload", videoUploadHandler)                                    //Complete
	root.HandleFunc("/api/v1/videos", indexVideoHandler)                                          //Complete
	root.HandleFunc("/api/v1/videos/{video}/like", likeVideoHandler)                              //Complete
	root.HandleFunc("/api/v1/notifications/register", registerPushNotifications).Methods("POST")  //Complete
	root.HandleFunc("/api/v1/standings/{division}", divisionStandings)                            //Complete

	http.Handle("/", root)
}

func main() {
	http.ListenAndServe(":80", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t := lookupTemplate("index")
	t.Execute(w, nil)
}
func contactHandler(w http.ResponseWriter, r *http.Request) {
	t := lookupTemplate("contact")
	t.Execute(w, nil)
}
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	t := lookupTemplate("about")
	t.Execute(w, nil)
}
func fieldsHandler(w http.ResponseWriter, r *http.Request) {
	t := lookupTemplate("fields")
	t.Execute(w, nil)
}
func videoTemplateHandler(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	videoId := vars["video"]
	_, interror := strconv.ParseInt(videoId, 10, 64)
	if interror != nil {
		e := lookupTemplate("error")
		e.Execute(w, nil)
		return
	}
	t := lookupTemplate("video")
	var buffer bytes.Buffer
	buffer.WriteString("SELECT id, url, likes FROM videos WHERE id='")
	buffer.WriteString(videoId)
	buffer.WriteString("';")
	rows, err := db.Query(buffer.String())
	if err != nil {
		response := Response{}
		response.Code = 500
		response.Message = "Ooops Internal Server Error!"
		encoder.Encode(&response)
		log.Printf(err.Error())
	}
	var v video
	for rows.Next() {
		rows.Scan(&v.Id, &v.Url, &v.Likes)
	}
	if v.Url == "" {
		e := lookupTemplate("error")
		e.Execute(w, nil)
		return
	}
	p := &Page{URL: v.Url}
	t.Execute(w, p)
}
