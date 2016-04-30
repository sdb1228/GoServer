package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	googleapi "google.golang.org/api/googleapi"
	storage "google.golang.org/api/storage/v1"
)

const (
	// This scope allows the application full control over resources in Google Cloud Storage
	scope = storage.DevstorageFullControlScope
)

var (
	projectID  = "goapi-1193"
	bucketName = "soccerlcvideostorage"
)

var db *sql.DB

type team struct {
	Name     string `json:"name"`
	Division string `json:"division"`
	Teamid   string `json:"teamid"`
}
type Field struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type PostField struct {
	Id      string `json:"id"`
	Address string `json:"address"`
	City    string `json:"city"`
	Zip     string `json:"zip"`
}
type Response struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
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
type video struct {
	Id    *int   `json:"id"`
	Url   string `json:"url"`
	Likes *int   `json:"likes"`
}

type videoInstallation struct {
	Id             *int   `json:"id"`
	Url            string `json:"url"`
	Likes          *int   `json:"likes"`
	InstallationId string `json:"installationId"`
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

type config struct {
	User     string
	Password string
	Database string
	Host     string
}

func init() {

	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config := config{}
	ferr := decoder.Decode(&config)
	if ferr != nil {
		log.Fatal("Error in decoding config: ", ferr)
	}

	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("user=%v dbname=%v host=%v password=%v sslmode=disable", config.User, config.Database, config.Host, config.Password))
	if err != nil {
		log.Fatal("Error in connecting to postgres databse: ", err)
	}
}

func credential_check(token string) error {
	if token != "XCF9-14PV-NLS1-3VCA" {
		return errors.New("Token does not match")
	}
	return nil

}
func response_builder(code int, message string) *Response {
	response := Response{}
	response.Code = code
	response.Message = message
	return &response
}

func registerPushNotifications(w http.ResponseWriter, r *http.Request) {
	installationId := r.FormValue("installationId")
	deviceToken := r.FormValue("deviceToken")
	encoder := json.NewEncoder(w)
	fmt.Println(installationId)
	fmt.Println(deviceToken)

	if installationId == "" || deviceToken == "" {
		log.Println("Devicetoken or installationId are nil")
		encoder.Encode(response_builder(403, "Internal server error please try again later"))
	}

	id, err := checkInstallationWithReturn(installationId, deviceToken)
	if err != nil {
		log.Println("Error in check instalation ", err)
		encoder.Encode(response_builder(403, "Internal server error please try again later"))
	}
	err = updateInstallationDeviceToken(id, deviceToken)
	if err != nil {
		log.Println("Error in updating devicetoken ", err)
		encoder.Encode(response_builder(403, "Internal server error please try again later"))
	}
	w.Header().Set("Content-Type", "application/json")
	encoder.Encode(response_builder(200, "Updated device token"))

}
func fieldsCorrectionHandler(w http.ResponseWriter, r *http.Request) {

	encoder := json.NewEncoder(w)
	var buffer bytes.Buffer

	buffer.WriteString("SELECT id, name FROM fields WHERE address = '' ORDER BY id;")
	fmt.Println(buffer.String())
	rows, err := db.Query(buffer.String())
	if err != nil {
		log.Println("Error in DB query of fields: ", err)
		encoder.Encode(response_builder(403, "Internal server error please try again later"))
		return
	}
	results := []Field{}
	for rows.Next() {
		var f Field
		rows.Scan(&f.Id, &f.Name)
		results = append(results, f)
	}
	fmt.Printf("%v", results)
	w.Header().Set("Content-Type", "application/json")
	encoder.Encode(&results)
}
func fieldsCorrectionPostHandler(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)
	var f PostField
	err := decoder.Decode(&f)
	if err != nil {
		fmt.Println(err)
	}
	var buffer bytes.Buffer

	buffer.WriteString("UPDATE fields SET address='")
	buffer.WriteString(f.Address)
	buffer.WriteString("', city='")
	buffer.WriteString(f.City)
	buffer.WriteString("', zip=")
	buffer.WriteString(f.Zip)
	buffer.WriteString(" WHERE id=")
	buffer.WriteString(f.Id)
	buffer.WriteString(";")
	fmt.Println(buffer.String())
	rows, err := db.Query(buffer.String())
	if err != nil {
		log.Println("Error in DB query of fields Update: ", err)
		encoder.Encode(response_builder(403, "Internal server error please try again later"))
		return
	}
	buffer.Reset()
	buffer.WriteString("SELECT id, name FROM fields WHERE address = '' ORDER BY id;")
	fmt.Println(buffer.String())
	rows, err = db.Query(buffer.String())
	results := []Field{}
	for rows.Next() {
		var f Field
		rows.Scan(&f.Id, &f.Name)
		results = append(results, f)
	}
	fmt.Printf("%v", results)
	w.Header().Set("Content-Type", "application/json")
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

/*
Verifys the Integer given is actually an integer
*/
func verifyInteger(integer string) error {
	_, interror := strconv.ParseInt(integer, 10, 64)
	if interror != nil {
		return interror
	}
	return nil
}

/*
This method handles the storage and linking of a users video
*/
func videoUploadHandler(w http.ResponseWriter, r *http.Request) {
	buffer := bytes.NewBuffer(nil)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	var videoName bytes.Buffer
	installationId := r.FormValue("installationId")
	email := r.FormValue("email")
	videoName.WriteString(uuid.NewV4().String())
	videoName.WriteString(".mp4")
	fmt.Println("off to google")
	r.ParseForm()
	uploadFile, uploadFileHeaders, err := r.FormFile("video")
	if uploadFile == nil {
		fmt.Println("uploadFIle is nil")
		encoder.Encode("{Response: Video Uploaded}")

	}
	contentLength := uploadFileHeaders.Header.Get("Content-Length")
	if contentLength != "" {
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			encoder.Encode("{Response: Video Uploaded}")
		}
	}

	_, err = io.Copy(buffer, uploadFile)
	if err != nil {
		http.Error(w, "Failed to receive entity", http.StatusInternalServerError)
		encoder.Encode("{Response: Video Uploaded}")
	}

	// Off To Google Storage

	err = googleCloudStorage(buffer, videoName.String(), installationId, email)
	if err == nil {
		response := Response{}
		response.Code = 200
		response.Message = "Video Uploaded"
		encoder.Encode(&response)
	} else {
		encoder.Encode(&err)
	}
}

/*
Uploads a byte buffer to the google cloud storage
*/
func googleCloudStorage(video *bytes.Buffer, objectName string, installationID string, email string) error {
	// Authentication is provided by the gcloud tool when running locally, and
	// by the associated service account when running on Compute Engine.
	client, err := google.DefaultClient(context.Background(), scope)
	if err != nil {
		log.Printf("Unable to get default client: %v", err)
		return err
	}
	service, err := storage.New(client)
	if err != nil {
		log.Printf("Unable to create storage service: %v", err)
		return err
	}

	if _, err := service.Buckets.Get(bucketName).Do(); err == nil {
	} else {
		// Create a bucket.
		fmt.Println("Bucket Doesn't exist.")
		return err
	}

	// Insert an object into a bucket.
	object := &storage.Object{Name: objectName}
	if err != nil {
		log.Printf("Error Saving File: %v", err)
		return err
	}
	videoContentType := googleapi.ContentType("video/mp4")
	if res, err := service.Objects.Insert(bucketName, object).PredefinedAcl("publicRead").Media(video, videoContentType).Do(); err == nil {
		log.Printf("Created object %v at location %v\n\n", res.Name, res.SelfLink)
		var buffer bytes.Buffer
		buffer.WriteString("https://storage.googleapis.com/")
		buffer.WriteString(bucketName)
		buffer.WriteString("/")
		buffer.WriteString(objectName)
		err = linkVideoToDatabse(buffer.String(), installationID, email)
		if err != nil {
			return err
		}
		return nil
	} else {
		log.Printf("Objects.Insert failed: %v", err)
		return err
	}
	return nil
}

/*
Stores video into personal database
*/
func linkVideoToDatabse(url, installationID, email string) error {
	_, err := db.Exec(
		"INSERT INTO videos (url, email, installation_id) VALUES ($1, $2, $3);",
		url,
		email,
		installationID,
	)
	if err != nil {
		return err
	}
	return nil
}

/*
Gets all videos and orders them from last inserted.  TODO make paginated.
*/
func indexVideoHandler(w http.ResponseWriter, r *http.Request) {
	installationId := r.FormValue("installationId")
	var buffer bytes.Buffer
	buffer.WriteString("SELECT videos.id, videos.url, videos.likes, lk.installationid FROM videos LEFT OUTER JOIN likes lk ON videos.id=lk.videoid AND lk.installationid='")
	buffer.WriteString(installationId)
	buffer.WriteString("' ORDER BY videos.id DESC;")
	rows, err := db.Query(buffer.String())
	if err != nil {
		log.Fatal(err)
	}

	results := []videoInstallation{}
	for rows.Next() {
		var v videoInstallation
		rows.Scan(&v.Id, &v.Url, &v.Likes, &v.InstallationId)
		results = append(results, v)
	}
	fmt.Printf("%v", results)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(&results)
}

/*
Registers a like with a video id
*/
func likeVideoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	video := vars["video"]
	installationId := r.FormValue("installationId")
	err := checkInstallation(installationId)
	var buffer bytes.Buffer
	buffer.WriteString("SELECT COUNT(*) FROM likes WHERE installationid='")
	buffer.WriteString(installationId)
	buffer.WriteString("' AND videoid='")
	buffer.WriteString(video)
	buffer.WriteString("';")
	fmt.Println(buffer.String())
	rows, err := db.Query(buffer.String())
	var count int
	for rows.Next() {
		rows.Scan(&count)
	}
	if count != 0 {
		_, err := db.Exec(
			"DELETE FROM likes WHERE installationid=$1 AND videoid=$2;",
			installationId,
			video,
		)
		if err != nil {
			fmt.Println(err)
		}
		_, err = db.Exec(
			"UPDATE videos SET likes = likes - 1 WHERE id=$1;",
			video,
		)
		var results [0]string
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.Encode(&results)
	} else {

		_, err = db.Exec(
			"INSERT INTO likes (installationid, videoid) VALUES ($1, $2);",
			installationId,
			video,
		)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec(
			"UPDATE videos SET likes = likes + 1 WHERE id=$1;",
			video,
		)
		var results [0]string
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.Encode(&results)
	}
}

/*
Get the games for a specific division between the start date and the end date formatted like yyyy-MM-dd hh:mm:ss ex: 2015-12-10 11:48:59
*/
func divisionGamesHandler(w http.ResponseWriter, r *http.Request) {
	facility := r.FormValue("facility")
	vars := mux.Vars(r)
	division := vars["division"]
	startDate := r.FormValue("startDate")
	endDate := r.FormValue("endDate")
	if facility == "" {
		facility = "0"
	}
	if startDate == "" {
		startDate = "1900-01-01"
	}
	if endDate == "" {
		endDate = "3100-01-01"
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

/*
Returns the games of a specific facility for today
*/
func todaysGamesHandler(w http.ResponseWriter, r *http.Request) {
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
Returns all the games for tomorrow for a specific facility
*/
func tomorrowGamesHandler(w http.ResponseWriter, r *http.Request) {
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

	rows, err := db.Query("SELECT f1.name AS field, f1.address AS address, a2.name AS hometeam, a1.name AS awayteam, games.gamesdatetime, games.hometeamscore, games.awayteamscore "+
		"FROM games "+
		"INNER JOIN fields f1 ON f1.id=games.field "+
		"INNER JOIN teams a1 ON games.awayteam=a1.teamid "+
		"INNER JOIN teams a2 ON games.hometeam=a2.teamid "+
		"WHERE games.awayteam=$1 "+
		"OR games.hometeam=$1 "+
		"ORDER BY games.gamesdatetime", team)
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

func checkInstallation(installationId string) error {
	var buffer bytes.Buffer
	buffer.WriteString("SELECT COUNT(*) FROM installation where installationid='")
	buffer.WriteString(installationId)
	buffer.WriteString("';")
	rows, err := db.Query(buffer.String())
	if err != nil {
		log.Println("Error in querying for installation ID: ", err)
		return err
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
			log.Println("Error in inserting installation : ", err)
			return err
		}
	}

	return nil
}

/*
Checks to see if we have the current installation in the database.  If we don't we will insert it and return the row id
*/
func checkInstallationWithReturn(installationId string, devicetoken string) (int64, error) {
	var buffer bytes.Buffer
	buffer.WriteString("SELECT id FROM installation where installationid='")
	buffer.WriteString(installationId)
	buffer.WriteString("';")
	rows, err := db.Query(buffer.String())
	if err != nil {
		log.Println("Error in querying for installation ID: ", err)
		return 0, err
	}

	var id int64
	var item sql.Result
	for rows.Next() {
		rows.Scan(&id)
	}
	fmt.Println(id)
	if id != 0 {
		return id, nil
	}
	item, err = db.Exec(
		"INSERT INTO installation (installationid, devicetoken) VALUES ($1, $2)",
		installationId,
		devicetoken,
	)
	if err != nil {
		log.Println("Error in inserting installation : ", err)
		return 0, err
	}
	value, _ := item.LastInsertId()
	fmt.Println(value)

	return value, nil
}

/*
Checks to see if we have the current installation in the database.  If we don't we will insert it and return the row id
*/
func updateInstallationDeviceToken(id int64, devicetoken string) error {
	_, err := db.Exec(
		"UPDATE installation SET devicetoken=$1 WHERE id=$2",
		devicetoken,
		id,
	)
	if err != nil {
		log.Println("Error in updating installationid : ", err)
		return err
	}
	return nil
}
