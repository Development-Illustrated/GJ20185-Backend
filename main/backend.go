package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var clients2 = make(map[*websocket.Conn]bool) // connected clients2
var action_broadcast = make(chan Action)      // action_broadcast channel
var client_broadcast = make(chan Client)      // action_broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/register/client", RegisterClient).Methods("POST")
	router.HandleFunc("/register/room", RegisterRoom).Methods("POST")
	router.HandleFunc("/rooms", ReturnRooms).Methods("GET")
	router.HandleFunc("/sendAction", sendAction).Methods("POST")
	router.HandleFunc("/clients", ReturnClients).Methods("GET")

	go func() {
		log.Print("Running http server on localhost:6969")
		log.Fatal(http.ListenAndServe(":6969", router))
	}()

	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)

	// Start listening for incoming actions
	go handleMessages()

	// Start listening for new clients wanting to join a room
	go handleClientMessages()

	// Start the server on localhost port 8000 and log any errors
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	//AddRoom(Room{"RoomId":"rm1"})
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error %v", err)
		return
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients2[ws] = true

	// Send it out to every client that is currently connected
	for client := range clients2 {
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(clients2, client)
		}
	}

	for {
		var msg Action
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients2, ws)
			break
		}
		// Send the newly received message to the action_broadcast channel
		action_broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Grab the next message from the action_broadcast channel
		msg := <-action_broadcast
		// Send it out to every client that is currently connected
		for client := range clients2 {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients2, client)
			}
		}
	}
}

func handleClientMessages() {
	for {
		// Grab the next message from the action_broadcast channel
		msg := <-client_broadcast
		// Send it out to every client that is currently connected
		for client := range clients2 {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients2, client)
			}
		}
	}
}

func ReturnClients(w http.ResponseWriter, r *http.Request) {
	log.Print("Returning clients")
	formattedStruct, _ := json.Marshal(GetClients())
	fmt.Fprintln(w, string(formattedStruct), http.StatusOK)
}

func sendAction(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var a Action
	err := decoder.Decode(&a)
	if err != nil {
		panic(err)
	}
	out := PerformAction(a)
	if out == true {
		fmt.Fprint(w, "this shit hot right now", http.StatusOK)
	} else {
		fmt.Fprint(w, "That client don't exist yo", http.StatusBadRequest)
	}
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
		fmt.Fprintln(w, "Registration for client: "+t.ClientId+" complete", http.StatusOK)
	} else {
		fmt.Fprintln(w, "Couldn't register client: "+t.ClientId, http.StatusBadRequest)
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
	fmt.Fprintln(w, "Registration for room: "+t.RoomId+" complete", http.StatusOK)

}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func ReturnRooms(w http.ResponseWriter, r *http.Request) {

	formattedStruct, _ := json.Marshal(GetRooms())
	fmt.Fprintln(w, string(formattedStruct))

}
