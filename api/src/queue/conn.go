package queue

import (
	"log"
	"os"
)

type Conn struct {
	host string
	port string
	user string
	pass string
}

func MakeConn() Conn {
	return Conn{
		envCheck("HOST"),
		envCheck("PORT"),
		envCheck("USER"),
		envCheck("PASS")}
}

func (e Conn) Host() string {
	return e.host
}

func (e Conn) Port() string {
	return e.port
}

func (e Conn) User() string {
	return e.user
}

func (e Conn) Pass() string {
	return e.pass
}

// Log if required env not set
func envCheck(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("%s not set\n", key)
	}
	return val
}

func (e Conn) UpdateFromEnv() {
	e.host = envCheck("HOST")
	e.port = envCheck("PORT")
	e.user = envCheck("USER")
	e.pass = envCheck("PASS")
}
