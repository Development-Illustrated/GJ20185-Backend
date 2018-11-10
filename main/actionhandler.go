package main

import (
	// "encoding/json"
	"fmt"
	// "log"
)

// var clients = make(map[string]Client)

type Action struct {
	clientId string
	actionKey string

	/**
	Example:
	{"ClientID":"cl1234",
	"actionKey":"Up"
	}
	*/

}

func PerformAction(action Action) {
	// Only add client if the room has been preregistered
	if GetClient(action.clientId) != nil {
		// clients[client.ClientId] = client
		// formattedStruct, _ := json.Marshal(client)
		// log.Println("Adding new client: " + string(formattedStruct))
		fmt.Println("perform an action" + action.actionKey) 
	// 	return true
	// } else {
	// 	return false
	}

}

