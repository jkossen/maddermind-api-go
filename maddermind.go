package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type Guess struct {
	Attempt []int
}

var dc = make(map[int][]int)
var dcDate int64

func getIP(r *http.Request) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
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

	ip, err := getIP(r)
	fmt.Println(time.Now().Local(), ":: Check attempt from ::", ip)

	var g Guess
	err = json.NewDecoder(r.Body).Decode(&g)

	if err != nil {
		fmt.Println("ERROR 1: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	codeLength := len(g.Attempt)
	switch codeLength {
	case 4:
	case 6:
	case 8:
		break
	default:
		http.Error(w, "Code length should be 4, 6 or 8", http.StatusBadRequest)
	}

	var c []int
	curEpoch := StartOfDayEpoch()
	if dcDate != curEpoch || dc[codeLength] == nil {
		fmt.Println("New day, new dawn. Trying to retrieve today's code")
		db := OpenDb()

		cStr, err := SelectTodaysChallenge(db, codeLength)
		var code []int

		switch err {
		case sql.ErrNoRows:
			fmt.Println("No code found for today. Generating new code")
			code = GenCode(codeLength)
			cStr = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(code)), ","), "")
			CreateTodaysChallenge(db, codeLength, cStr)
		default:
			checkErr(err)
		}

		err = json.Unmarshal([]byte(cStr), (&c))
		checkErr(err)

		dcDate = curEpoch
		dc[codeLength] = c

		CloseDb(db)
	}

	res, err := ChkAttempt(g.Attempt, dc[codeLength])

	if err != nil {
		fmt.Println("ERROR 2: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := make(map[string]interface{})
	resp["attempt"] = g.Attempt
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
	w.Write(r)
}

func loadEnv() {
	env := os.Getenv("MADDERMIND_ENV")
	if "" == env {
		env = "development"
	}

	godotenv.Load(".env." + env + ".local")
	if "test" != env {
		godotenv.Load(".env.local")
	}

	godotenv.Load(".env." + env)
	godotenv.Load() // The Original .env
}

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
