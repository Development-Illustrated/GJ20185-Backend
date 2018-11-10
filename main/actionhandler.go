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

	/**
	Example:
	{"ClientID":"cl1234",
	"Command":"Up"
	}
	*/

}

func PerformAction(action Action) bool {
	// Only add client if the room has been preregistered
	if GetClient(action.ClientId) != nil {
		log.Println("Client: " + action.ClientId + " performs action: " + action.Command)
		action_broadcast <- action
		return true
	} else {
		log.Println("Client: " + action.ClientId + " doesn't exist.")
		return false
	}
}
