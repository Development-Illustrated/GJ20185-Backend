package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Command struct {
	ClientId string
	Command  string

	/**
	Example:
	{"ClientID":"cl1234",
	"ClientType":"controller",
	"RoomId":"rm1234"
	}
	*/

}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/register/client", RegisterClient).Methods("POST")
	router.HandleFunc("/register/room", RegisterRoom).Methods("POST")
	router.HandleFunc("/rooms", ReturnRooms).Methods("GET")
	router.HandleFunc("/sendAction", sendAction).Methods("POST")
	router.HandleFunc("/clients", ReturnClients).Methods("GET")

	log.Fatal(http.ListenAndServe(":6969", router))
}
func ReturnClients(w http.ResponseWriter, r *http.Request) {
	formattedStruct, _ := json.Marshal(GetClients())
	fmt.Fprintln(w, string(formattedStruct))
}

func sendAction(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var a Action
	err := decoder.Decode(&a)
	if err != nil {
		panic(err)
	}
	PerformAction(a)
}

// Endpoints!
func RegisterClient(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var t Client
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

	fmt.Println(t.ClientId)

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
