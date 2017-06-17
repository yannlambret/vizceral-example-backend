package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

// Listening address
var address = "localhost:9090"

// Clients connected to the backend
var clients = make(map[*websocket.Conn]net.Addr)

// Synchronize accesses to clients map
var clientsMutex sync.Mutex

// Brodacast channel
var broadcast = make(chan map[string]interface{})

// Upgrade an HTTP connection to a WebSocket connection
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func computeMessages() {
	for {
		data, err := ioutil.ReadFile("sample_data_simple.json")
		if err != nil {
			log.Println(err)
		}
		var traffic map[string]interface{}

		if err := json.Unmarshal(data, &traffic); err != nil {
			log.Println(err)
		}
		broadcast <- traffic
		time.Sleep(5 * time.Second)
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		clientsMutex.Lock()
		for client, addr := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println(err)
				log.Println("closing remote connection", addr)
				client.Close()
				delete(clients, client)
			}
		}
		clientsMutex.Unlock()
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Switch to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	clientsMutex.Lock()
	clients[ws] = ws.RemoteAddr()
	clientsMutex.Unlock()
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()
	go computeMessages()
	log.Fatal(http.ListenAndServe(address, nil))
}
