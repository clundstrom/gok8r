package main

import (
	f "fmt"
	"log"
	"net"
	"net/http"
)

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func testRoute(w http.ResponseWriter, r *http.Request) {
	f.Fprintf(w, "%s", GetOutboundIP())
}

func main() {
	http.HandleFunc("/", testRoute)
	log.Fatal(http.ListenAndServe(":8000", nil))
}