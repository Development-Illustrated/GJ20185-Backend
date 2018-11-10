package main

import "fmt"

var clients = make(map[string]Client)

type Client struct {
	ClientId   string
	ClientType string
	RoomId     string

	/**
	Example:
	{"ClientID":"cl1234",
	"ClientType":"controller",
	"RoomId":"rm1234"
	}
	*/

}

func AddClient(client Client) {
	clients[client.ClientId] = client

}

func GetClient(ClientId string) Client {
	return clients[ClientId]
}

// Print all clients
func PrintClients() {
	fmt.Println(clients)
}

// Print individual client
func PrintClient(ClientId string) {
	fmt.Println(GetClient(ClientId))
}
