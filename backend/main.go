package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	var addr = flag.String("addr", ":8181", "The addr of the application.")
	flag.Parse() // parse the flags

	rooms := newRooms()
	tokens := newTokens()
	rtr := mux.NewRouter()
	//
	rtr.HandleFunc("/room/{roomNumber}/{isMaster}/{nickname}", rooms.handleHttp)
	http.HandleFunc("/roomnumber", func(w http.ResponseWriter, r *http.Request) {
		rooms.roomNumber(w, r)
	})

	http.HandleFunc("/tokens", func(w http.ResponseWriter, r *http.Request) {
		tokens.addTokens(w, r)
	})

	http.Handle("/", rtr)

	// get the room going
	go rooms.run()

	// start the web server
	log.Println("Starting web server on", *addr)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
