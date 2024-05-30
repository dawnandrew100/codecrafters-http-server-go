package main

import (
	"fmt"
	"net"
	"os"
    "strings"
    "bytes"

)
// I'm going to refactor this because I can forsee this becoming complicated

const (
    OK          = "HTTP/1.1 200 OK\r\n"
	BAD_REQUEST = "HTTP/1.1 400 Bad Request\r\n\r\n"
	NOT_FOUND   = "HTTP/1.1 404 Not Found\r\n\r\n"
    CREATED     = "HTTP/1.1 201 Created\r\n"
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

 type HttpRequest struct {
     method  string
     path    string
     version string
     headers map[string]string
     body    string
 }

func requestParser(buffer []byte) HttpRequest {
    // parse request buffer
	bufString := strings.Split(string(buffer), "\r\n")
	
    method := strings.Split(bufString[0], " ")[0]
	path := strings.Split(bufString[0], " ")[1]
	version := strings.Split(bufString[0], " ")[2]
	headers := make(map[string]string)
	for i := 1; i < len(bufString)-2; i++ {
		header := strings.Split(bufString[i], ": ")
		headers[header[0]] = header[1]
	}
	body := strings.Split(string(buffer), "\r\n\r\n")[1]
	
    return HttpRequest{
		method:  method,
		path:    path,
        version: version,
		headers: headers,
		body:    body,
	}
}

func handleConnection(conn net.Conn) {
    //frees memory by closing connection at end of function
    defer conn.Close()
    
    buf := make([]byte, 1024)
    conn.Read(buf)
    buf = bytes.Trim(buf, "\x00")

    req := requestParser(buf)

    method := req.method
    path := req.path
    
    fmt.Printf("Port: %s\nPath: %s\nHTTP version: %s\n", req.headers["Host"], path, req.version)

    var response string

    switch {
    case path == "/":
        response = OK + "\r\n"

    case strings.Contains(path, "echo"):
        acceptedEncoding := req.headers["Accept-Encoding"]
        if acceptedEncoding == "gzip" || acceptedEncoding == "brotli"{
            echostring := strings.Split(path, "/")
            response = compressedResponseBuilder(OK, acceptedEncoding, "text/plain", len(echostring[2]), echostring[2])
        } else {
            echostring := strings.Split(path, "/")
            response = responseBuilder(OK, "text/plain", len(echostring[2]), echostring[2])
        }
    
    case path == "/user-agent":
        user_agent := req.headers["User-Agent"]
        response = responseBuilder(OK, "text/plain", len(user_agent), user_agent)

    case method == "GET" && strings.Contains(path, "files"):
        directory := os.Args[2]
        fileName := strings.TrimPrefix(path, "/files/")
        data, err := os.ReadFile(directory + fileName)
				if err != nil {
                    response = NOT_FOUND 
				} else {
                    dataString := string(data)
                    response = responseBuilder(OK, "application/octet-stream", len(dataString), dataString)
				}
    
    case method == "POST" && strings.Contains(path, "files"):
        directory := os.Args[2]
        fileName := strings.TrimPrefix(path, "/files/")
        filepath := fmt.Sprintf("%s%s", directory, fileName)

        if directory == "" {
            fmt.Println("Error reading file directory")
            response = BAD_REQUEST
        }
	    f, err := os.Create(filepath)
	    if err != nil {
		    fmt.Printf("Unable to create file: %s\n", filepath)
	    }
	    _, err = f.WriteString(req.body)
	    if err != nil {
		    fmt.Println("Unable to write to file")
	    }
        response = responseBuilder(CREATED, "application/octet-stream", len(req.body), req.body)

    default:
        response = NOT_FOUND
    }

    conn.Write([]byte(response))

}

func responseBuilder(statusLine string, contentType string, contentLength int, body string) string {
    return fmt.Sprintf("%sContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", statusLine, contentType, contentLength, body)
}

func compressedResponseBuilder(statusLine string, contentEncoding string, contentType string, contentLength int, body string) string {
    return fmt.Sprintf("%sContent-Encoding: %s\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", statusLine, contentEncoding, contentType, contentLength, body)
}
