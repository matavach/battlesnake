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
	http.HandleFunc("/move", handleMove)
	http.HandleFunc("/start", HandleStart)
	http.HandleFunc("/end", HandleEnd)
	http.HandleFunc("/", HandleRoot)
	err := http.ListenAndServe(":"+c.Port, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	response := info()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ERROR: Failed to encode info response, %s", err)
	}
}

func HandleStart(w http.ResponseWriter, r *http.Request) {
	state := GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		log.Printf("ERROR: Failed to decode start json, %s", err)
		return
	}

	start(&state)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}
func HandleEnd(w http.ResponseWriter, r *http.Request) {
	state := GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		log.Printf("ERROR: Failed to decode end json, %s", err)
		return
	}

	end(&state)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func handleMove(w http.ResponseWriter, r *http.Request) {
	request := GameState{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Failed to parse JSON for /move. Error: %s", err)
	}
	response := move(&request)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ERROR: Failed to encode move response, %s", err)
		return
	}

}
