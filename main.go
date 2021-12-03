package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"go-pad/controllers"
	"go-pad/models"
	"go-pad/websocket"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	// Parse values from command line
	flag.Parse()

	// Initiate redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Initiate hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initiate models
	padModel := models.NewPadModel(rdb)

	// Initiate controllers
	padController := controllers.NewPadController(padModel)

	// Init routes
	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/favicon.ico", http.StripPrefix("/", http.FileServer(http.Dir("./public/"))))
	router.Handle("/js/{name}", http.StripPrefix("/js/", http.FileServer(http.Dir("./public/js/"))))
	router.Handle("/css/{name}", http.StripPrefix("/css/", http.FileServer(http.Dir("./public/css/"))))
	router.HandleFunc("/ws/{name}", func(w http.ResponseWriter, r *http.Request) {websocket.ServeWs(hub, w, r)})
	router.HandleFunc("/", padController.ServeHome).Methods(http.MethodGet)
	router.HandleFunc("/{name}", padController.ServePad).Methods(http.MethodGet)

	// Init webserver
	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      router,
		Addr:         *addr,
	}

	log.Println("Starting backend application on " + *addr)
	log.Fatal(server.ListenAndServe())
}
