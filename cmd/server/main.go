package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/buffup/GolangTechTask/cmd/server/internal/handlers"
)

const PG_ENV =  "PG_CONNSTRING"

func main() {
	if _, ok := os.LookupEnv(PG_ENV); !ok {
		log.Println("PG_CONNSTRING environment variable not set")
		return
	}
	store, err := handlers.NewPGStore(handlers.PGConnString(os.Getenv(PG_ENV)))
	if err != nil {
		log.Println(fmt.Errorf("main#NewPGStore: %w", err))
		return
	}

	routes := handlers.Routes(store)
	if err := http.ListenAndServe(":8080", routes); err != nil {
		panic(err)
	}
}
