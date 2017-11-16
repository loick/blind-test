package main

import (
	"encoding/json"
	"log"
	"net/http"

	"os"

	"github.com/gorilla/websocket"
	"github.com/matryer/goblueprints/chapter1/trace"
)

/*type rooms struct {
	rooms []room
}*/

type room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clients.
	forward chan []byte

	// join is a channel for clients wishing to join the room.
	join chan *client

	// leave is a channel for clients wishing to leave the room.
	leave chan *client

	// clients holds all current clients in this room.
	clients map[*client]bool

	// tracer will receive trace information of activity
	// in the room.
	tracer trace.Tracer
}

type rooms struct {
	rooms map[string]*room
}

// newRoom makes a new room that is ready to
// go.
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

func newRooms() *rooms {
	return &rooms{
		make(map[string]*room),
	}
}

func (rms *rooms) run() {
	for {
		for _, r := range rms.rooms {
			select {
			case client := <-r.join:
				// joining
				r.clients[client] = true
				r.tracer.Trace("New client joined")
			case client := <-r.leave:
				// leaving
				delete(r.clients, client)
				close(client.send)
				r.tracer.Trace("Client left")
			case msg := <-r.forward:
				r.tracer.Trace("Message received: ", string(msg))
				// forward message to all clients
				for client := range r.clients {
					client.send <- msg
					r.tracer.Trace(" -- sent to client")
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *rooms) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("ServeHTTP:", err)
		return
	}

	room := r.rooms["toto"]

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   room,
	}
	room.join <- client
	defer func() { room.leave <- client }()
	go client.write()
	client.read()
}

func (rooms *rooms) roomNumber(w http.ResponseWriter, r *http.Request) {
	roomNum := createRoomNumber()
	a := make(map[string]string)
	a["roomNumber"] = roomNum
	res, _ := json.Marshal(a)
	room := newRoom()
	room.tracer = trace.New(os.Stdout)
	rooms.rooms[roomNum] = room

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}