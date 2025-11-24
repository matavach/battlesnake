package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func Run(c *Config) {

	if _, err := strconv.Atoi(c.Port); err != nil {
		log.Fatalf("%s is not a valid port number", c.Port)
	}
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/test", handleMove)
	err := http.ListenAndServe(":"+c.Port, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {

	log.Println("Recieved!")
}

func handleMove(w http.ResponseWriter, r *http.Request) {
	request := GameState{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Failed to parse JSON for /move. Error: %s", err)
	}

}
