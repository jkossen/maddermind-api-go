package main

import (
	"encoding/json"
	"fmt"
	"jkossen/maddermind-backend-go/sqlite"
	"net/http"
	"os"
	"time"
)

func handleTokenRequest(w http.ResponseWriter, _ *http.Request) {
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

	ip, err := getIP(r)
	fmt.Println(time.Now().Local(), ":: Check attempt from ::", ip)

	var guess Guess
	err = json.NewDecoder(r.Body).Decode(&guess)

	if err != nil {
		fmt.Println("ERROR 1: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	codeLen := len(guess.Attempt)
	switch codeLen {
	case 4:
	case 6:
	case 8:
		break
	default:
		http.Error(w, "Challenge length should be 4, 6 or 8", http.StatusBadRequest)
	}

	var challenge Challenge = challenge{}
	var cs = &sqlite.ChallengeStorage{}
	cs.DSN(os.Getenv("DSN"))
	curEpoch := StartOfDayEpoch(time.Now())
	_, hasCode := dc[codeLen]
	if dcDate != curEpoch || !hasCode {
		challenge = challenge.RetrieveOrGen(cs, curEpoch, codeLen)
		dcDate = curEpoch
		dc[codeLen] = challenge
	}

	res, err := dc[codeLen].Check(guess.Attempt)

	if err != nil {
		fmt.Println("ERROR 2: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := make(map[string]interface{})
	resp["attempt"] = guess.Attempt
	resp["result"] = res

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
	_, err := w.Write(r)
	checkErr(err)
}
