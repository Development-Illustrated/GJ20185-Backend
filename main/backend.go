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
	router.HandleFunc("/register/client", RegisterClient).Methods("POST")
	router.HandleFunc("/register/room", RegisterRoom).Methods("POST")
	router.HandleFunc("/rooms", ReturnRooms).Methods("GET")

	log.Fatal(http.ListenAndServe(":6969", router))
}

// Endpoints!
func RegisterClient(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var t Client
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

	if AddClient(t) {
		fmt.Fprintln(w, "Registration for client: "+t.ClientId+" complete")
	} else {
		fmt.Fprintln(w, "Couldn't register client: "+t.ClientId)
	}
}

func RegisterRoom(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var t Room
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

	AddRoom(t)
	fmt.Fprintln(w, "Registration for room: "+t.RoomId+" complete")

}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func ReturnRooms(w http.ResponseWriter, r *http.Request) {

	formattedStruct, _ := json.Marshal(GetRooms())
	fmt.Fprintln(w, string(formattedStruct))

}
