package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"curls/ai"
	"curls/models"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	response := models.BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "",
		Color:      "#888888",
		Head:       "default",
		Tail:       "default",
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal(err)
	}
}

func HandleStart(w http.ResponseWriter, r *http.Request) {
	request := models.GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("START\n")
}

func HandleMove(w http.ResponseWriter, r *http.Request) {
	request := models.GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	response := ai.NextMove(request)

	fmt.Printf("MOVE: %s\n", response.Move)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal(err)
	}
}

func HandleEnd(w http.ResponseWriter, r *http.Request) {
	request := models.GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("END\n")
}

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/start", HandleStart)
	http.HandleFunc("/move", HandleMove)
	http.HandleFunc("/end", HandleEnd)

	fmt.Printf("Starting Battlesnake Server at http://0.0.0.0:%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
