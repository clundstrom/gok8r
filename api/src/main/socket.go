package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gok8r/src/queue"
	"log"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Job struct {
	Id      string `json:"id"`
	Seconds int    `json:"seconds"`
}
type SocketChannel struct {
	SocketChan chan string `json:"socketChan"`
	RemoteAddr string      `json:"remoteAddr"`
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

func getUserChannel(socketChannel map[string]SocketChannel, r *http.Request) chan string {
	sock := socketChannel[r.URL.Query().Get("id")]
	return sock.SocketChan
}

// Handles queuing of jobs.
// Responses are sent via WS
func queueJob(socketChannel map[string]SocketChannel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var jobDetails Job
		err := json.NewDecoder(r.Body).Decode(&jobDetails)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		messageChan = getUserChannel(socketChannel, r)

		if messageChan == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Illegal socket id."))
			fmt.Println("Illegal socket id.")
			return
		}

		scheduled := queue.ScheduleWork(r.URL.Query().Get("id"), strconv.Itoa(jobDetails.Seconds), messageChan)

		if !scheduled {
			messageChan <- "Error: Could not queue job."
			w.WriteHeader(http.StatusBadRequest)
			if err != nil {
				return
			}
		} else {
			w.WriteHeader(http.StatusOK)
			messageChan <- "Job queued."
			if err != nil {
				return
			}
		}
	}
}

// openSocket creates a websocket to the client, and opens
// a string channel from which messages transmitted
// in and endless loop or until the socket is terminated.
// An initial 'connected' message is sent
// on opening the socket.
func openSocket(connections map[string]SocketChannel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		c, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer func(c *websocket.Conn) {
			err := c.Close()
			if err != nil {
				delete(connections, c.RemoteAddr().String())
			}
		}(c)

		messageChan = make(chan string)

		// Assign a unique id to each socket, and store in a map
		// When socket is closed, delete channel reference from list
		// On subsequent post request, include the response socket uuid
		connections[r.URL.Query().Get("id")] = SocketChannel{
			SocketChan: messageChan,
			RemoteAddr: r.RemoteAddr,
		}

		// close channel on exit
		defer func() {
			delete(connections, r.RemoteAddr)
			messageChan = nil
			close(messageChan)
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
}
