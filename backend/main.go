package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse() // parse the flags

	rooms := newRooms()

	http.Handle("/room", rooms)
	http.HandleFunc("/roomnumber", func(w http.ResponseWriter, r *http.Request) {
		rooms.roomNumber(w, r)
	})

	// get the room going
	go rooms.run()

	// start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
