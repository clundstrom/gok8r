package main

import (
	"bufio"
	f "fmt"
	"gok8r/src/queue"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

var messageChan chan string

// Open a server-sent-event stream.
// On a successful connection, a Connected to ${IP} is sent.
func getSSE(w http.ResponseWriter, r *http.Request) {

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

func sendSSE(message string) http.HandlerFunc {
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

func input() {
	reader := bufio.NewReader(os.Stdin)
	f.Println("---------------------")
	f.Println("Gok8r")
	f.Println("Usage: say <message>")
	f.Println("---------------------")
	for {
		f.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\r\n", "", -1)

		if strings.HasPrefix(text, "say") == true {
			stripped := strings.TrimPrefix(text, "say")
			stripped = strings.TrimSpace(stripped)

			if messageChan == nil {
				log.Printf("No client connected to channel.")

			} else {
				messageChan <- stripped
				log.Printf("Send: %s", stripped)
			}

			text = ""
		}
	}
}

func main() {
	go input()
	http.HandleFunc("/api/v1/socket", openSocket)
	http.HandleFunc("/api/v1/echo", echo)
	http.HandleFunc("/api/v1/stream", getSSE)
	http.HandleFunc("/api/v1/sendsse", sendSSE("hello through sse"))
	http.HandleFunc("/api/v1/sendws", sendWs("hello through ws"))
	http.HandleFunc("/api/v1/queue", queueResponse(5))
	http.HandleFunc("/", defaultRoute)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
