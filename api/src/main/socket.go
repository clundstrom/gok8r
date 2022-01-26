package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func echo(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func sendWs(message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "%d", 200)

		if messageChan != nil {
			messageChan <- message
		}

		if err != nil {
			log.Printf("%s for %s -- 500\n", r.Host, r.RequestURI)
			return
		}
		log.Printf("%s for %s -- 200\n", r.Host, r.RequestURI)
	}
}

func openSocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	messageChan = make(chan string)

	// close channel on exit
	defer func() {
		close(messageChan)
		messageChan = nil
		log.Printf("client connection closed")
	}()

	err = c.WriteMessage(websocket.TextMessage, []byte(`{"message":"Connected"}`))

	for {
		select {
		case message := <-messageChan:
			sendJson := []byte(fmt.Sprintf(`{"message":"%s"}`, message))
			err = c.WriteMessage(websocket.TextMessage, sendJson)
		case <-r.Context().Done():
			return
		}
	}
}
