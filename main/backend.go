package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/register", Register).Methods("POST")
	//router.HandleFunc("/test", test).Methods("POST")

	log.Fatal(http.ListenAndServe(":6969", router))
}

func LogObject(input Client) {
	formattedStruct, _ := json.Marshal(input)
	log.Println("Logging object: " + string(formattedStruct))
}

// Endpoints!
func Register(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var t Client
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	LogObject(t)

	AddClient(t)
	fmt.Fprintln(w, "Registration for "+t.ClientId+" complete")

	// Show current registered clients
	PrintClients()

}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}
