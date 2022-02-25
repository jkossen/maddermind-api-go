package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

type Guess struct {
	Attempt []int
}

func handleTokenRequest(w http.ResponseWriter, r *http.Request) {
	token := SepEveryNth(RandString(16), 4, "-")

	resp := make(map[string]string)
	resp["token"] = token

	b, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	okResponse(w, b)
}

func handleCheckAttemptRequest(w http.ResponseWriter, r *http.Request) {
	// just return for preflight call
	if r.Method != "POST" {
		return
	}

	// code to break
	c := []int{1, 6, 4, 1}

	var g Guess
	err := json.NewDecoder(r.Body).Decode(&g)

	if err != nil {
		fmt.Println("ERROR 1: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := ChkAttempt(g.Attempt, c)

	resp := make(map[string]interface{})
	resp["attempt"] = g.Attempt
	resp["result"] = res

	if err != nil {
		fmt.Println("ERROR 2: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse, jsonError := json.Marshal(resp)

	if jsonError != nil {
		fmt.Println("Unable to encode JSON")
	}

	okResponse(w, jsonResponse)
}

func okResponse(w http.ResponseWriter, r []byte) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Max-Age", "86400")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(r)
}

func main() {
	r := mux.NewRouter()

	// routing
	r.HandleFunc("/chk", handleCheckAttemptRequest).Methods("GET", "POST", "OPTIONS")
	r.HandleFunc("/new", handleTokenRequest).Methods("GET")

	c := cors.New(cors.Options{
		AllowedMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowedOrigins:     []string{"http://localhost:3000"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Content-Type", "Bearer", "Bearer ", "content-type", "Origin", "Accept"},
		OptionsPassthrough: true,
	})

	handler := c.Handler(r)

	// start serving
	log.Fatal(http.ListenAndServe(
		"localhost:8080",
		handler))
}
