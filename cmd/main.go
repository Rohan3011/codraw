package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"text/template"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	connections = make(map[*websocket.Conn]bool)
	broadcast   = make(chan Message)
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type DrawLine struct {
	PrevPoint    *Point `json:"prevPoint"`
	CurrentPoint Point  `json:"currentPoint"`
	Color        string `json:"color"`
}

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", home)
	http.HandleFunc("/ws", wsHandler)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		log.Println("Shutting down the server...")
		// Add cleanup logic if needed
		os.Exit(0)
	}()

	go handleMessages()

	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe("localhost:8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade failed: ", err)
		return
	}
	defer conn.Close()

	connections[conn] = true

	handleWebSocketConnection(conn)
}

func handleWebSocketConnection(conn *websocket.Conn) {
	defer func() {
		delete(connections, conn)
		conn.Close()
	}()

	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("Error during message reading: ", err)
			break
		}

		switch msg.Type {
		case "client-ready":
			log.Println(msg.Type)
			if err := conn.WriteJSON(Message{
				Type: "get-canvas-state",
				Data: msg.Data,
			}); err != nil {
				log.Println("Error during message writing:", err)
			}
		case "canvas-state":
			log.Println(msg.Type)
			if err := conn.WriteJSON(Message{
				Type: "canvas-state-from-server",
				Data: msg.Data,
			}); err != nil {
				log.Println("Error during message writing:", err)
			}
			// Handle other-type message
		case "draw-line":
			log.Println(msg.Type)

			broadcast <- msg // Broadcast the message to all connected clients
		case "clear":
			log.Println(msg.Type)
			broadcast <- msg // Broadcast the message to all connected clients
		default:
			log.Println(msg.Type)
			// Handle unknown message type
		}

	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for conn := range connections {
			if err := conn.WriteJSON(msg); err != nil {
				log.Println("Error during message writing:", err)
			}
		}
	}
}
