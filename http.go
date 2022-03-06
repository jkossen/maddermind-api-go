package main

import (
	"encoding/json"
	"fmt"
	"jkossen/maddermind-backend-go/datetime"
	"jkossen/maddermind-backend-go/mastermind"
	"jkossen/maddermind-backend-go/sqlite"
	"jkossen/maddermind-backend-go/strutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func getIP(r *http.Request) (string, error) {
	// Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	// Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	// Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("No valid ip found")
}

func handleTokenRequest(w http.ResponseWriter, _ *http.Request) {
	token := strutil.SepEveryNth(strutil.RandString(16), 4, "-")

	resp := make(map[string]string)
	resp["token"] = token

	b, err := json.Marshal(resp)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "unable to serialize response", 500)
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
	log.Println("check attempt from ", ip)

	var guess Guess
	err = json.NewDecoder(r.Body).Decode(&guess)

	if err != nil {
		log.Println("err: " + err.Error())
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	codeLen := len(guess.Attempt)
	switch codeLen {
	case 4:
	case 6:
	case 8:
		break
	default:
		log.Println("http: received guess with invalid codelen: ", codeLen)
		http.Error(w, "Challenge length should be 4, 6 or 8", http.StatusBadRequest)
		return
	}

	var challenge = mastermind.Challenge{}
	var cs = &sqlite.ChallengeStorage{}
	cs.DSN(os.Getenv("DSN"))

	err = cs.Open()
	if err != nil {
		log.Println("http: error opening storage:", err)
	}

	defer func(cs *sqlite.ChallengeStorage) {
		err := cs.Close()
		if err != nil {
			log.Println("http: error closing storage:", err)
		}
	}(cs)

	curEpoch := datetime.StartOfDayEpoch(time.Now())
	_, hasCode := dc[codeLen]
	if dcDate != curEpoch || !hasCode {
		log.Println("http: getting code for timestamp", curEpoch, "with codeLen", codeLen)
		challenge, err = challenge.GetOrCreate(cs, curEpoch, codeLen)
		if err != nil {
			log.Println("http: error from RetrieveOrGen: ", err)
		}
		dcDate = curEpoch
		dc[codeLen] = challenge
	}

	res, err := dc[codeLen].Check(guess.Attempt)

	if err != nil {
		log.Println("ERROR 2: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := make(map[string]interface{})
	resp["attempt"] = guess.Attempt
	resp["result"] = res

	jsonResp, err := json.Marshal(resp)

	if err != nil {
		log.Println("http: unable to JSON encode response")
		http.Error(w, "Unable to encode response", http.StatusBadRequest)
		return
	}

	okResponse(w, jsonResp)
}

func okResponse(w http.ResponseWriter, r []byte) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Max-Age", "86400")
	w.Header().Set("Content-Type", "application/json")

	_, err := w.Write(r)

	if err != nil {
		log.Println(err)
	}
}
