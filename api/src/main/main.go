package main

import (
	f "fmt"
	"gok8r/src/queue"
	"log"
	"net"
	"net/http"
)

var messageChan chan string

// Open a server-sent-event stream.
// On a successful connection, a Connected to ${IP} is sent.
func getSSE() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		messageChan = make(chan string)
		log.Printf("open connection %s", r.Host)

		// close channel on exit
		defer func() {
			close(messageChan)
			messageChan = nil
			log.Printf("client connection closed")
		}()

		flusher, _ := w.(http.Flusher)
		_, _ = f.Fprintf(w, "%s %s\n\n", "data: Connected to ", GetOutboundIP())
		flusher.Flush()

		for {
			select {
			case message := <-messageChan:
				_, err := f.Fprintf(w, "%s\n\n", message)
				if err != nil {
					return
				}

				flusher.Flush()
			case <-r.Context().Done():
				return
			}
		}
	}
}

func sendMessage(message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if messageChan != nil {
			log.Printf("Write %b to %s", len(message), r.Host)
			messageChan <- "data: " + message
		}
	}
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func defaultRoute(w http.ResponseWriter, r *http.Request) {
	_, err := f.Fprintf(w, "%s", GetOutboundIP())
	if err != nil {
		log.Printf("%s for %s -- 500\n", r.Host, r.RequestURI)
		return
	}
	log.Printf("%s for %s -- 200\n", r.Host, r.RequestURI)
}

func queueResponse(seconds int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !queue.Work(seconds) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Message queued"))
		} else {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("Message queued"))
			if err != nil {
				return
			}
		}
	}
}

func main() {
	http.HandleFunc("/api/v1/socket", echo)
	http.HandleFunc("/api/v1/stream", getSSE())
	http.HandleFunc("/api/v1/sendmessage", sendMessage("hello client"))
	http.HandleFunc("/api/v1/queue", queueResponse(5))
	http.HandleFunc("/", defaultRoute)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
