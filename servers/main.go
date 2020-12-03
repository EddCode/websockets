package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world a simple response from go server")
}

func ReaderWriter(conn *websocket.Conn) {
	msg := ""
	for msg != "end" {
		messageType, stream, err := conn.ReadMessage()

		if err != nil {
			fmt.Println("Error: ", err)
		}

		msg = string(stream)
		log.Println(msg)

		write := []byte("Recived")

		if err := conn.WriteMessage(messageType, write); err != nil {
			log.Println(err)
			return
		}
	}
	conn.Close()
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("New connection from %s\n", r.RemoteAddr)

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Error on ws", err)
	}

	ReaderWriter(ws)
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	addrArg := os.Args[1:]

	if len(addrArg) == 0 {
		log.Panic("Missing add your port")
	}

	var addr string = fmt.Sprintf(":%s", addrArg[0])

	setupRoutes()

	fmt.Printf("Server listening on port %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))

}
