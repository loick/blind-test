package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/matryer/goblueprints/chapter1/trace"
)

type roomNumber struct {
	RoomNumber string `json:"roomNumber"`
}

type token struct {
	Token      string `json:"token"`
	RoomNumber string `json:"roomNumber"`
}

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

type tokens struct {
	tokens map[string]string
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

func newTokens() *tokens {
	return &tokens{
		tokens: make(map[string]string),
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

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (r *rooms) handleHttp(w http.ResponseWriter, req *http.Request) {
	spew.Dump("la")
	params := mux.Vars(req)
	roomNumber := params["roomNumber"]
	isMaster := params["isMaster"]

	var isMasterBool bool
	if isMaster == "1" {
		isMasterBool = true
	} else {
		isMasterBool = false
	}

	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("ServeHTTP:", err)
		return
	}

	room, ok := r.rooms[roomNumber]
	if !ok {
		log.Printf("sortie")
		return
	}

	client := &client{
		socket:   socket,
		send:     make(chan []byte, messageBufferSize),
		room:     room,
		isMaster: isMasterBool,
	}
	room.join <- client
	defer func() { room.leave <- client }()
	go client.write()
	client.read()
}

func (rooms *rooms) roomNumber(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		var roomNmbr roomNumber
		json.Unmarshal(body, &roomNmbr)
		if _, ok := rooms.rooms[roomNmbr.RoomNumber]; ok {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		roomNum := createRoomNumber()
		a := make(map[string]string)
		a["roomNumber"] = roomNum
		res, _ := json.Marshal(a)
		room := newRoom()
		room.tracer = trace.New(os.Stdout)
		rooms.rooms[roomNum] = room

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(res))
	}
}

func (tokens *tokens) addTokens(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		var token token
		json.Unmarshal(body, &token)
		tokens.tokens[token.RoomNumber] = token.Token
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusBadRequest)

}
