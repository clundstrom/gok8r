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

// Handles queuing of jobs.
// Responses are sent via WS
func queueJob() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var jobDetails Job
		err := json.NewDecoder(r.Body).Decode(&jobDetails)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		if !queue.ScheduleWork(strconv.Itoa(jobDetails.Seconds)) {
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
