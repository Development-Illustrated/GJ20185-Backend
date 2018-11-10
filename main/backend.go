package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/register", Register).Methods("POST")

	log.Fatal(http.ListenAndServe(":6969", router))
}
func Register(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Welcome to registration biznatch")
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}
