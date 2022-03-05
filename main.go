package main

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

// Container for the player's guess
type Guess struct {
	Attempt []int
}

// Daily Challenge, one per codeLength
var dc = make(map[int]Challenge)
var dcDate int64

func main() {
	loadEnv()

	r := mux.NewRouter()

	// routing
	r.HandleFunc("/chk", handleCheckAttemptRequest).Methods("GET", "POST", "OPTIONS")
	r.HandleFunc("/new", handleTokenRequest).Methods("GET")

	c := cors.New(cors.Options{
		AllowedMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowedOrigins:     []string{"http://localhost:3000", "https://madmuon.com"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Content-Type", "Bearer", "Bearer ", "content-type", "Origin", "Accept"},
		OptionsPassthrough: true,
	})

	handler := c.Handler(r)

	// start serving
	log.Fatal(http.ListenAndServe(
		os.Getenv("HOST")+":"+os.Getenv("PORT"),
		handler))
}
