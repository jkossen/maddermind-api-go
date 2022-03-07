package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"jkossen/maddermind-backend-go/api"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	loadEnv()

	// disable date and datetime for logging, assume eg systemd takes care of those
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	log.Println("starting up")

	r := mux.NewRouter()

	// routing
	r.HandleFunc("/chk", api.Check).Methods("GET", "POST", "OPTIONS")
	r.HandleFunc("/new", api.Token).Methods("GET")

	c := cors.New(cors.Options{
		AllowedMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowedOrigins:     strings.Fields(os.Getenv("ALLOWED_ORIGINS")),
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
