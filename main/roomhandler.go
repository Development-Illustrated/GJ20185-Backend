package main

import (
	"encoding/json"
	"fmt"
	"log"
)

var rooms = make(map[string]*Room)

type Room struct {
	RoomId    string
	ClientIds []string

	/**
	Example:
	{
	"RoomId":"rm1234",
	"Clients":["cl1234","cl4233"]
	}
	*/

}

func AddRoom(room Room) {
	formattedStruct, _ := json.Marshal(room)
	log.Println("Adding new room: " + string(formattedStruct))
	rooms[room.RoomId] = &room

}

func AddClientToRoom(c Client) bool {

	for _, b := range rooms[c.RoomId].ClientIds {
		if b == c.ClientId {
			log.Print("Player is being added to a room they are already in. Dickwads writing the front end fucked up.")
			return false
		}
	}

	rooms[c.RoomId].ClientIds = append(rooms[c.RoomId].ClientIds, c.ClientId)
	client_broadcast <- c
	log.Print("Client: " + c.ClientId + " added to room")
	return true
}

// Attempt to get a room from the cache, if it doesnt exist return nil
func GetRoom(RoomId string) *Room {
	Room, ok := rooms[RoomId]
	if ok {
		return Room
	} else {
		return nil
	}
}

func GetRooms() map[string]*Room {
	return rooms
}

// Print all clients
func PrintRooms() {
	fmt.Println(GetRooms())
}

// Print individual client
func PrintRoom(RoomId string) {
	fmt.Println(GetRoom(RoomId))
}
