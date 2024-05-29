package main

import (
	"fmt"
	"net"
	"os"
    "strings"
    "bytes"
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

func handleConnection(conn net.Conn) {
    //frees memory by closing connection at end of function
    defer conn.Close()
    
    buf := make([]byte, 1024)
    conn.Read(buf)
    buf = bytes.Trim(buf, "\x00")

    bufString := strings.Split(string(buf), "\n")
    request := strings.Split(bufString[0], " ")
    host := bufString[1]

    method := request[0]
    path := request[1]
    version := request[2]
    fmt.Printf("Port: %s\nPath: %s\nHTTP version: %s\n", host, path, version)

    var response string

    switch {
    case path == "/":
        response = "HTTP/1.1 200 OK\r\n\r\n"

    case strings.Contains(path, "echo"):
        encoding := strings.Split(bufString[2], ":")
        if strings.Contains(bufString[2], "Accept-Encoding") && encoding[1] == " gzip" {
            echostring := strings.Split(path, "/")
            response = "HTTP/1.1 200 OK\r\n"
            response += fmt.Sprintf("Content-Encoding: %s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n", encoding[1], len(echostring[2]))
            response += echostring[2]
        } else {
            echostring := strings.Split(path, "/")
            response = "HTTP/1.1 200 OK\r\n"
            response += fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n\r\n", len(echostring[2]))
            response += echostring[2]
        }
    
    case path == "/user-agent":
        user_agent := bufString[2]
        user_agent_echo := strings.Split(user_agent, " ")
        response = "HTTP/1.1 200 OK\r\n"
        // must subtract one becuase length also counts carriage return as character
        response += fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n\r\n", len(user_agent_echo[1])-1)
        response +=user_agent_echo[1]

    case method == "GET" && strings.Contains(path, "files"):
        directory := os.Args[2]
        fileName := strings.TrimPrefix(path, "/files/")
        data, err := os.ReadFile(directory + fileName)
				if err != nil {
                    response = "HTTP/1.1 404 Not Found\r\n\r\n"
				} else {
                    dataString := string(data)
					response = "HTTP/1.1 200 OK\r\n"
                    response += fmt.Sprintf("Content-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n", len(dataString))
                    response += dataString
				}
    
    case method == "POST" && strings.Contains(path, "files"):
        directory := os.Args[2]
        fileName := strings.TrimPrefix(path, "/files/")
        filepath := fmt.Sprintf("%s%s", directory, fileName)

	    f, err := os.Create(filepath)
	    if err != nil {
		    fmt.Printf("Unable to create file: %s\n", filepath)
	    }
	    _, err = f.WriteString(bufString[len(bufString)-1])
	    if err != nil {
		    fmt.Println("Unable to write to file")
	    }
		response = "HTTP/1.1 201 Created\r\n"
        response += fmt.Sprintf("Content-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n", len(bufString[len(bufString)-1]))
        response += bufString[len(bufString)-1]

    default:
        response = "HTTP/1.1 404 Not Found\r\n\r\n"
    }

    conn.Write([]byte(response))

}
