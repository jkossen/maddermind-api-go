package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"strings"
)

type DailyChallenges struct {
	Date   int64
	CodeL1 []int
	CodeL2 []int
	CodeL3 []int
}

type Guess struct {
	Attempt []int
}

var dc DailyChallenges

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

	var c []int
	curEpoch := StartOfDayEpoch()
	if dc.Date != curEpoch {
		fmt.Println("Generating new code")
		db := OpenDb()
		var code []int
		cStr, err := SelectTodaysChallenge(db, 4)
		switch err {
		case sql.ErrNoRows:
			code = GenCode(4)
			cStr = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(code)), ","), "")
			CreateTodaysChallenge(db, 4, cStr)
		default:
			checkErr(err)
		}

		err = json.Unmarshal([]byte(cStr), (&c))
		checkErr(err)

		dc.Date = curEpoch
		dc.CodeL1 = c

		CloseDb(db)
	}

	var g Guess
	err := json.NewDecoder(r.Body).Decode(&g)

	if err != nil {
		fmt.Println("ERROR 1: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := ChkAttempt(g.Attempt, dc.CodeL1)

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
		AllowedOrigins:     []string{"http://localhost:3000", "https://madmuon.com"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Content-Type", "Bearer", "Bearer ", "content-type", "Origin", "Accept"},
		OptionsPassthrough: true,
	})

	handler := c.Handler(r)

	// start serving
	log.Fatal(http.ListenAndServe(
		"127.0.0.1:12001",
		handler))
}
