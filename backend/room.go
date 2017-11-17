package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"fmt"

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

type spotify struct {
	Album struct {
		AlbumType string `json:"album_type"`
		Artists   []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"artists"`
		AvailableMarkets []string `json:"available_markets"`
		ExternalUrls     struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href   string `json:"href"`
		ID     string `json:"id"`
		Images []struct {
			Height int    `json:"height"`
			URL    string `json:"url"`
			Width  int    `json:"width"`
		} `json:"images"`
		Name string `json:"name"`
		Type string `json:"type"`
		URI  string `json:"uri"`
	} `json:"album"`
	Artists []struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href string `json:"href"`
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
		URI  string `json:"uri"`
	} `json:"artists"`
	AvailableMarkets []string `json:"available_markets"`
	DiscNumber       int      `json:"disc_number"`
	DurationMs       int      `json:"duration_ms"`
	Explicit         bool     `json:"explicit"`
	ExternalIds      struct {
		Isrc string `json:"isrc"`
	} `json:"external_ids"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href        string `json:"href"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Popularity  int    `json:"popularity"`
	PreviewURL  string `json:"preview_url"`
	TrackNumber int    `json:"track_number"`
	Type        string `json:"type"`
	URI         string `json:"uri"`
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

	ladderBoard map[string]int

	artist string
}

type tokens struct {
	tokens map[string]string
}

type rooms struct {
	rooms map[string]*room
	song  map[string]string
}

// newRoom makes a new room that is ready to
// go.
func newRoom() *room {
	return &room{
		forward:     make(chan []byte),
		join:        make(chan *client),
		leave:       make(chan *client),
		clients:     make(map[*client]bool),
		ladderBoard: make(map[string]int),
		tracer:      trace.Off(),
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
		make(map[string]string),
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
				if client.isMaster {
					preview, artist := sendTrack()
					r.artist = artist
					client.send <- []byte(preview)
					r.tracer.Trace(fmt.Sprintf("Artist: %s, preview: %s", preview, artist))
				}

				r.ladderBoard[client.nickname] = 0
			case client := <-r.leave:
				// leaving
				delete(r.clients, client)
				close(client.send)
				r.tracer.Trace("Client left")
			case msg := <-r.forward:
				r.tracer.Trace("Message received: ", string(msg))
				// forward message to all clients
				sendLadderBoard := false
				for client := range r.clients {
					client.send <- msg
					r.tracer.Trace(" -- sent to client")
					if r.artist == string(msg) {
						r.ladderBoard[r.artist] = r.ladderBoard[r.artist] + 1
						spew.Dump(r.ladderBoard)
						sendLadderBoard = true
					}
				}

				if sendLadderBoard == true {
					for client := range r.clients {
						if client.isMaster == true {
							jsonify, _ := json.Marshal(r.ladderBoard)
							client.send <- jsonify
						}
					}
				}
			}
		}
	}
}

func sendTrack() (string, string) {
	trackId := "2Fa5PbnEZixzN910CloiiS"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", trackId), nil)
	req.Header.Set("Authorization", "Bearer BQC3QXAsKeL9I8_CIokazrkgYauGtM7bbf0SHNfwzls7saCTtGFSjZ-eAiTOf8H3Zb7duyHheWUp2lZknOBZjMpFt7okb__d1NbV-CqJsC0ZqvQlC0-vqxApwoMdQniPBi1JKmlWJs4Auh9WrQldF34")
	res, err := client.Do(req)
	if err != nil {
		fmt.Sprintf("err %s", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Sprintf("err %s", err)
	}
	var spotify spotify
	json.Unmarshal(body, &spotify)

	return spotify.PreviewURL, spotify.Artists[0].Name
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
	params := mux.Vars(req)
	roomNumber := params["roomNumber"]
	isMaster := params["isMaster"]
	nickname := params["nickname"]
	log.Printf(fmt.Sprintf("roomNumber: %s, isMaster: %s, nickname: %s", roomNumber, isMaster, nickname))
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
		socket:     socket,
		send:       make(chan []byte, messageBufferSize),
		room:       room,
		isMaster:   isMasterBool,
		roomNumber: roomNumber,
		nickname:   nickname,
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
