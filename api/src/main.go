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

func defaultRoute(w http.ResponseWriter, r *http.Request) {
	_, err := f.Fprintf(w, "%s", GetOutboundIP())

	if err != nil {
		log.Printf("%s for %s -- 500\n", r.Host, r.RequestURI)
		return
	}
	log.Printf("%s for %s -- 200\n", r.Host, r.RequestURI)
}

func main() {
	http.HandleFunc("/", defaultRoute)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
