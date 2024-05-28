package main

import (
	"log"
	"net"
	"os"
    "strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	log.Println("Logs from your program will appear here!")

    // Creates TCP server listening on port 4221	
	 l, err := net.Listen("tcp", "0.0.0.0:4221")
	 if err != nil {
	 	log.Println("Failed to bind to port 4221")
	 	os.Exit(1)
	 }
     defer l.Close()
	
     conn, err := l.Accept()
	 if err != nil {
	 	log.Println("Error accepting connection: ", err.Error())
	 	os.Exit(1)
	 }

     handleConnection(conn)
}

func handleConnection(conn net.Conn) {
    //frees memory by closing connection at end of function
    defer conn.Close()
    
    buf := make([]byte, 0, 1024)
    conn.Read(buf)

    bufString := strings.Split(string(buf), "\n")

    request := strings.Split(bufString[0], " ")
    // headers := strings.Split(bufString[1], " ")

    // method := request[0]
    path := request[1]
    // version := request[2]
    log.Printf("Path is %s:\n", path)

    var response string = "HTTP/1.1 404 Not Found\r\n\r\n"
    if path == "/" {
        response == "HTTP/1.1 200 OK\r\n\r\n"
    }

    _, err = conn.Write([]byte(response))
    if err != nil {
        log.Println("Failed to write response")
        os.Exit(1)
    }

}
