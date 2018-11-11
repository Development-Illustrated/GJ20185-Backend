package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var client_connections = make(map[*websocket.Conn]bool) // connected client_connections
var action_broadcast = make(chan Action)                // action_broadcast channel
var client_broadcast = make(chan Client)                // action_broadcast channel

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
	go handleActionMessages()

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
	client_connections[ws] = true

	// Send it out to every client that is currently connected
	for client := range client_connections {
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(client_connections, client)
		}
	}

	for {
		var msg Action
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(client_connections, ws)
			break
		}
		// Send the newly received message to the action_broadcast channel
		action_broadcast <- msg
	}
}

func handleActionMessages() {
	for {
		// Grab the next message from the action_broadcast channel
		msg := <-action_broadcast
		// Send it out to every client that is currently connected
		for client := range client_connections {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(client_connections, client)
			}
		}
	}
}

func handleClientMessages() {
	for {
		// Grab the next message from the action_broadcast channel
		msg := <-client_broadcast
		// Send it out to every client that is currently connected
		for client := range client_connections {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(client_connections, client)
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
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Action sent"))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Client doesn't exist, or client isn't in room"))
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
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Registration for client: " + t.ClientId + " complete"))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Can't add client: " + t.ClientId + " to room " + t.RoomId + " because room doesn't exist!"))

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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Registration for room: " + t.RoomId + " complete"))

}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hey there partner"))
}

func ReturnRooms(w http.ResponseWriter, r *http.Request) {

	formattedStruct, _ := json.Marshal(GetRooms())
	fmt.Fprintln(w, string(formattedStruct))

}
