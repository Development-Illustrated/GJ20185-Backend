package main

import (
	// "encoding/json"
	// "fmt"
	"log"
)

// var clients = make(map[string]Client)

type Action struct {
	ClientId string
	Command  string
	RoomId   string

	/**
	Example:
	{"ClientID":"cl1234",
	"Command":"Up"
	}
	*/

}

func PerformAction(action Action) bool {
	log.Print("Performing action")
	// Only add client if the room has been preregistered
	client := GetClient(action.ClientId)
	if client != nil && action.RoomId == client.RoomId {
		log.Println("Client: " + action.ClientId + " performs action: " + action.Command)
		action_broadcast <- action
		return true
	} else {
		log.Println("Client: " + action.ClientId + " doesn't exist. Or client isn't in room " + action.RoomId)
		return false
	}
}
