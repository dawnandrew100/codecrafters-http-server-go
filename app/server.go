package main

import (
	"fmt"
	"net"
	"os"
    "strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	 // Uncomment this block to pass the first stage
	 l, err := net.Listen("tcp", "0.0.0.0:4221")
	 defer l.Close()
     if err != nil {
	 	fmt.Println("Failed to bind to port 4221")
	 	os.Exit(1)
	 }
   
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

     go handleConnection(conn)
   }
 }

 type Request struct {
     Method      string
     Target      string
     HTTPVersion string
     Headers     map[string]string
 }

 func parseRequest(request []string) *Request {
    clientrequest := strings.Split(request[0], " ")
	
    headers := make(map[string]string)
	for _, element := range request[1:] {
		if strings.Contains(element, "\x00") {
			break
		}
		if element != "" {
			headerSplit := strings.Split(element, ":")
			headers[headerSplit[0]] = strings.TrimSpace(headerSplit[1])
			continue
		}
	}

    return &Request{
        Method: request[0],
        Target: request[1],
        HTTPVersion: request[2],
        Headers: headers,
    }
 }

func handleConnection(conn net.Conn) {
    //frees memory by closing connection at end of function
    defer conn.Close()
    
    buf := make([]byte, 1024)
    conn.Read(buf)
    bufString := strings.Split(string(buf), "\n")

    parseRequest(bufString)

    var response string

    //host := bufString[1]
    user_agent := bufString[2]
    switch {
    case Request.Target == "/":
        response = "HTTP/1.1 200 OK\r\n\r\n"

    case strings.Contains(Request.Target, "echo"):
        echostring := strings.Split(Request.Target, "/")
        response = "HTTP/1.1 200 OK\r\n"
        response += fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n\r\n", len(echostring[2]))
        response += echostring[2]

    case Request.Target == "/user-agent":
        response = "HTTP/1.1 200 OK\r\n"
        // must subtract one becuase length also counts carriage return as character
        response += fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n\r\n", len(request.Headers["User-Agent"]))
        response += request.Headers["User-Agent"]

    default:
        response = "HTTP/1.1 404 Not Found\r\n\r\n"
    }

    conn.Write([]byte(response))

}
