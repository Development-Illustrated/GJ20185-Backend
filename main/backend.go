package main

import (
	"encoding/json"
	"fmt"
	"github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	// "time"
	"github.com/gorilla/websocket"
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

var clients2 = make(map[*websocket.Conn]bool) // connected clients2
var broadcast = make(chan Message)           // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Message object
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
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
		log.Fatal(http.ListenAndServe(":6969", router))
	}()

	// go func() {
	// 	StartSocketServer()
	// }()

	// log.Print("This should run after listen finish")

	// // Dont end the program
	// for {
	// 	time.Sleep(9999999)
	// }

	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)

	// Start listening for incoming chat messages
	go handleMessages()

	// Start the server on localhost port 8000 and log any errors
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
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
		err := client.WriteJSON("hi Luke")
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(clients2, client)
		}
	}

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients2, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
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

func NewSocket(so socketio.Socket) {

	so.Join("chat")

	so.Emit("chat", "hello message")
	log.Println("on connection")
	so.On("chat", func(msg string) {
		log.Println("recieved message", msg)
		so.Emit("chat", msg)
	})
	// Socket.io acknowledgement example
	// The return type may vary depending on whether you will return
	// For this example it is "string" type
	so.On("chat message with ack", func(msg string) string {
		return msg
	})
	so.On("disconnection", func() {
		log.Println("disconnected from chat")
	})
}

func StartSocketServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"hello\": \"world\"}"))
	})

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation below for more options.
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(so socketio.Socket) {

		so.Join("chat")

		so.Emit("chat", "hello message")
		log.Println("on connection")
		so.On("chat", func(msg string) {
			log.Println("recieved message", msg)
			so.Emit("chat", msg)
		})
		// Socket.io acknowledgement example
		// The return type may vary depending on whether you will return
		// For this example it is "string" type
		so.On("chat message with ack", func(msg string) string {
			return msg
		})
		so.On("disconnection", func() {
			log.Println("disconnected from chat")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	mux.Handle("/socket.io/", server)
	mux.Handle("/assets", http.FileServer(http.Dir("./assets")))

	// provide default cors to the mux
	handler := cors.Default().Handler(mux)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: false,
	})

	// decorate existing handler with cors functionality set in c
	handler = c.Handler(handler)

	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", handler))

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
