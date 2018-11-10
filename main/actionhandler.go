package main

import (
	// "encoding/json"
	// "fmt"
	"log"
)

// var clients = make(map[string]Client)

type Action struct {
	ClientId  string
	ActionKey string

	/**
	Example:
	{"ClientID":"cl1234",
	"ActionKey":"Up"
	}
	*/

}

func PerformAction(action Action) {
	// Only add client if the room has been preregistered
	if GetClient(action.ClientId) != nil {
		log.Println("Client: " + action.ClientId + " performs action: " + action.ActionKey)
	} else {
		log.Println("Client: " + action.ClientId + " doesn't exist.")
	}

	broadcast <- action

}

//TO DO
// func SendActionToGameRoom( ){

// }
