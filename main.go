package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

var (
	configStruct config
	db           *sql.DB
)

type Page struct {
	URL  string
	Body string
}

type config struct {
	User         string
	Password     string
	Database     string
	Host         string
	ProdPort     string
	DevPort      string
	IsProduction string
}

func init() {

	root := mux.NewRouter()
	root.StrictSlash(true)

	// root
	root.HandleFunc("/", indexHandler)
	root.NotFoundHandler = http.HandlerFunc(indexHandler)

	// templates
	root.HandleFunc("/video/{video}", videoTemplateHandler)
	root.HandleFunc("/about", aboutHandler)
	root.HandleFunc("/contact", contactHandler)
	root.HandleFunc("/fields/update", indexHandler)

	// assets
	root.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("public/assets"))))

	// API endpoints
	root.HandleFunc("/api/v1/divisions/{division}/games", divisionGamesHandler)  //Complete
	root.HandleFunc("/api/v1/divisions/{division}/teams", divisionsTeamsHandler) //Complete
	// games
	root.HandleFunc("/api/v1/games/{team}", gamesForTeamHandler) //Complete
	//teams
	root.HandleFunc("/api/v1/teams/", teamsHandler)                      //Complete
	root.HandleFunc("/api/v1/teams/{leagueId}", teamsForFacilityHandler) //Complete
	//favorites
	root.HandleFunc("/api/v1/favorites/{team}", addFavoriteTeamHandler).Methods("POST")      //Complete
	root.HandleFunc("/api/v1/facilitys/{league}/divisions", facilityDivisionsHandler)        //Complete
	root.HandleFunc("/api/v1/favorites/{team}", removeFavoriteTeamHandler).Methods("DELETE") //Complete
	root.HandleFunc("/api/v1/favorites", favoriteTeamsHandler)                               //Complete
	root.HandleFunc("/api/v1/favorites/games/", favoriteTeamsGamesHandler)                   //Complete
	//fields
	root.HandleFunc("/api/v1/fields/correction", fieldsCorrectionHandler)                         //Complete
	root.HandleFunc("/api/v1/fields/postCorrection", fieldsCorrectionPostHandler).Methods("POST") //Complete
	//notifications
	root.HandleFunc("/api/v1/notifications/register", registerPushNotifications).Methods("POST") //Complete
	//today/tomorrowgames
	root.HandleFunc("/api/v1/todaysGames/{league}", todaysGamesHandler)     //Complete
	root.HandleFunc("/api/v1/tomorrowGames/{league}", tomorrowGamesHandler) //Complete
	//videos
	root.HandleFunc("/api/v1/videoUpload", videoUploadHandler)       //Complete
	root.HandleFunc("/api/v1/videos", indexVideoHandler)             //Complete
	root.HandleFunc("/api/v1/videos/{video}/like", likeVideoHandler) //Complete
	//standings
	root.HandleFunc("/api/v1/standings/{division}", divisionStandings) //Complete

	http.Handle("/", root)
}

func main() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	ferr := decoder.Decode(&configStruct)
	if ferr != nil {
		log.Fatal("Error in decoding config: ", ferr)
	}
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("user=%v dbname=%v host=%v password=%v sslmode=disable", configStruct.User, configStruct.Database, configStruct.Host, configStruct.Password))
	if err != nil {
		log.Fatal("Error in connecting to postgres databse: ", err)
	}
	if strings.Contains("True", configStruct.IsProduction) {
		http.ListenAndServe(configStruct.ProdPort, nil)
	} else {
		http.ListenAndServe(configStruct.DevPort, nil)
	}
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
