package main

import (
	"encoding/json"
	"fmt"
	"log"
)

var clients = make(map[string]Client)

type Client struct {
	ClientId string
	RoomId   string

	/**
	Example:
	{"ClientID":"cl1234",
	"ClientType":"controller",
	"RoomId":"rm1234"
	}
	*/

}

func AddClient(client Client) bool {
	// Only add client if the room has been preregistered
	if GetRoom(client.RoomId) != nil {
		clients[client.ClientId] = client
		formattedStruct, _ := json.Marshal(client)
		log.Println("Adding new client: " + string(formattedStruct))
		return true
	} else {
		return false
	}
}

// Attempt to get a Client from the cache, if it doesnt exist return nil
func GetClient(ClientId string) *Client {
	Client, ok := clients[ClientId]
	if ok {
		return &Client
	} else {
		return nil
	}
}

// Print all clients
func PrintClients() {
	fmt.Println(clients)
}

// Print individual client
func PrintClient(ClientId string) {
	fmt.Println(GetClient(ClientId))
}
