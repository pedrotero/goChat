package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	addr = ":8080"
)

var (
	users    []string
	convos   []*Convo
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func getConvoByName(name string) *Convo {
	for _, convo := range convos {
		if convo.name == name {
			return convo
		}
	}
	return nil
}

func convosHandler(w http.ResponseWriter, _ *http.Request) {
	convosList := fmt.Sprintf("Total convos: %d", len(convos))
	for _, convo := range convos {
		convosList += fmt.Sprintf("\n/%s", convo.name)
	}
	_, err := fmt.Fprint(w, convosList)
	if err != nil {
		return
	}
}

func clientHandler(w http.ResponseWriter, r *http.Request) {

	room := mux.Vars(r)["room"]
	log.Println("new connection from", r.RemoteAddr, ", connecting to", room)
	//upgrade connection buffer size
	conn, err := upgrader.Upgrade(w, r, nil)
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("Error closing conn", err)
		}
	}(conn)
	if err != nil {
		log.Println(err)
		return
	}
	convo := getConvoByName(room)
	if convo == nil { //if conversation doesnt exist, create new one
		convo = newConvo(room)
	}
	writeMessage, err := json.Marshal(convo.messages)
	if err != nil {
		return
	}

	err = conn.WriteMessage(1, writeMessage)
	if err != nil {
		return
	}

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println("Unexpected close by", conn.RemoteAddr())
				return
			}
			log.Println("Error reading ws message", err)
			return
		}
		log.Printf("Received message type: %d, from websocket: %s", messageType, data)
		addMessage(convo, Message{"aaaa", string(data), time.Now()})

	}

}

func main() {
	router := mux.NewRouter()
	log.Println("starting server")
	convos = append(convos, newConvo("gen"))
	router.HandleFunc("/", convosHandler)
	router.HandleFunc("/{room}", clientHandler)
	err := http.ListenAndServe(addr, router)
	if err != nil {
		return
	}
}
