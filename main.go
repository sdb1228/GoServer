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

	// polls
	polls := root.PathPrefix("/polls").Subrouter()
	polls.Methods("POST").Path("/").HandlerFunc(httpCreatePoll)
	polls.Methods("GET").Path("/{key}/").HandlerFunc(httpShowPoll)

	http.Handle("/", root)
}

func main() {
	http.ListenAndServe(":8960", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t := lookupTemplate("index")
	t.Execute(w, nil)
}
