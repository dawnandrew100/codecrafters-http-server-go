package main

import (
	"fmt"
	"net"
	"os"
    "strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("!Logs from your program will appear here!")

	 // Uncomment this block to pass the first stage
	 l, err := net.Listen("tcp", "0.0.0.0:4221")
	 defer l.Close()
     if err != nil {
	 	fmt.Println("Failed to bind to port 4221")
	 	os.Exit(1)
	 }
	
     conn, err := l.Accept()
	 if err != nil {
	 	fmt.Println("Error accepting connection: ", err.Error())
	 	os.Exit(1)
	 }
    handleConnection(conn)
 }

func handleConnection(conn net.Conn) {
    //frees memory by closing connection at end of function
    defer conn.Close()
    
    buf := make([]byte, 1024)
    conn.Read(buf)

    bufString := strings.Split(string(buf), "\n")
    fmt.Println(bufString[1]) 
    request := strings.Split(bufString[0], " ")
    headers := strings.Split(bufString[1], "\n")
    fmt.Println(headers)
    // method := request[0]
    path := request[1]
    // version := request[2]
    fmt.Printf("Path is %s:\n", path)

    var response string = "HTTP/1.1 404 Not Found\r\n\r\n"
    if path == "/" {
        response = "HTTP/1.1 200 OK\r\n\r\n"
    }

    if strings.Contains(path, "echo") {
        echostring := strings.Split(path, "/")
        response = "HTTP/1.1 200 OK\r\n"
        response += fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n\r\n", len(echostring[2]))
        response += echostring[2]

    }

    if path == "/user-agent" {
        user_agent := strings.Split(headers[1], " ")
        response = "HTTP/1.1 200 OK\r\n"
        response += fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n\r\n", len(user_agent[1]))
        response +=user_agent[1]
    }

    conn.Write([]byte(response))

}
