package main

import (
	"fmt"
	"log"
	"net"
	"bufio"
	"time"
)

type client chan<- string // an outgoing message channel

type Timer struct {
	C <-chan Time
	// contains filtered or unexported fields
}
type Time struct {
	// contains filtered or unexported fields
}
var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		// 1. Send the message to all the clients
		case msg := <-messages:
			go func(){
			for cli := range clients{
				
					cli <-msg
			}
			}()		
		// 2. Update the clients map
		case cli := <-entering:
			clients[cli] = true
		// 3. Update the clients map and close the client channel
		case cli := <-leaving:
			delete(clients,cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages

	go clientWriter(conn, ch)
	// Client's IP address
	who := conn.RemoteAddr().String()

	// 1. Send a message to the new client confirming their identity
	// e.g. "You are 1.2.3"
	ch<-"You are " + who

	// 2. Broadcast a message to all the clients that a new client has joined
	// e.g. "1.2.3. has joined"

	messages<-who + " has joined"

	// 3. Send the client to the entering channel
	entering<-ch

	// 4. Use a scanner (e.g. bufio.NewScanner) to read from the client and broadcast the incoming messages to all the clients
	input := bufio.NewScanner(conn)
	// 5. Send the client to the leaving channel
	var timer1 = time.NewTimer(1000*time.Second)
	go func() {
		 <-timer1.C
			conn.Close()
	}()
	

	for input.Scan(){
		timer1.Reset(1000*time.Second)
		messages <- who + ": " + input.Text()
	}
	

	leaving<-ch
	messages<-who + " has left"
	// 6. Broadcast a message to all clients that the client has left
	// e.g. "1.2.3. has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}



func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}