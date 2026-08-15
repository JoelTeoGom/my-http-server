package myhttp

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

// Define the HTTP methods as string constants
const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
	HEAD    = "HEAD"
)

type HttpResponse struct {
	StatusLine string            // Status line: includes version, status code and message.
	Headers    map[string]string // Headers: key-value pairs.
	Body       []byte            // Body: response content (can be HTML, JSON, etc.).
}
type HttpRequest struct {
	Method  string
	URI     string
	Version string
	Headers map[string]string
	Body    string
}

type HandleFunc func(req *HttpRequest, res *HttpResponse)

type Server struct {
	Routes map[string]HandleFunc
}

func NewServer() *Server {
	return &Server{
		Routes: make(map[string]HandleFunc),
	}
}

func (http *Server) HandleFunction(method, path string, handler HandleFunc) {
	key := fmt.Sprintf("%s:%s", method, path)
	http.Routes[key] = handler
}

func (http *Server) Serve(req *HttpRequest) *HttpResponse {
	path := req.URI
	method := req.Method
	key := fmt.Sprintf("%s:%s", method, path)

	handler, exists := http.Routes[key]
	if !exists {
		log.Printf("handler is not registered")
		return &HttpResponse{
			StatusLine: "HTTP/1.1 404 Not Found",
			Headers:    map[string]string{"Content-Type": "text/plain"},
			Body:       []byte("404 Not Found"),
		}
	}

	response := &HttpResponse{
		StatusLine: "HTTP/1.1 200 OK",
		Headers:    make(map[string]string),
	}
	handler(req, response)

	return response
}

func (f HandleFunc) ServeHTTP(req *HttpRequest, res *HttpResponse) {
	f(req, res)
}

func (http *Server) HttpServer(address string) {
	// Create the listener on the given address
	listener, err := net.Listen("tcp4", address)
	if err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}

	defer listener.Close()
	log.Printf("Server listening on %s\n", address)

	for {
		conn, err := listener.Accept()
		log.Printf("New connection from: %s", conn.RemoteAddr())

		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go handleConnection(conn, http)
	}
}

func handleConnection(conn net.Conn, server *Server) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	var head strings.Builder
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		head.WriteString(line)
		if line == "\r\n" || line == "\n" {
			break
		}
	}

	request, err := parseHttpRequest(head.String())
	if err != nil {
		log.Printf("Error parsing the data: %v", err)
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}
	if cl := request.Headers["Content-Length"]; cl != "" {
		n, err := strconv.Atoi(cl)
		if err != nil || n < 0 {
			conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
			return
		}
		body := make([]byte, n)
		if _, err := io.ReadFull(reader, body); err != nil {
			return
		}
		request.Body = string(body)
	}
	// here we serve the response
	response := server.Serve(request)

	res := formatHttpResponse(response)
	log.Printf("Processed response:\n%s", res)
	// Send the response to the client
	_, err = conn.Write(res)
	if err != nil {
		log.Printf("Error sending the response to the client: %v", err)
	}

	log.Println("Response sent to the client.")
}

func formatHttpResponse(res *HttpResponse) []byte {
	// Compute the body length
	bodyLength := len(res.Body)

	// Make sure the Content-Length header is included
	res.Headers["Content-Length"] = fmt.Sprintf("%d", bodyLength)

	// Build the headers
	headers := ""
	for key, value := range res.Headers {
		headers += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	// Build the full response
	return []byte(fmt.Sprintf(
		"%s\r\n%s\r\n%s",
		res.StatusLine,   // Status line
		headers,          // Headers
		string(res.Body), // Body
	))
}

func parseHttpRequest(request string) (*HttpRequest, error) {
	lines := strings.Split(request, "\r\n")
	if len(lines) < 1 {
		return nil, fmt.Errorf("malformed request")
	}

	// Parse the request line (Method, URI, Version)
	requestLine := strings.Fields(lines[0])
	if len(requestLine) < 3 {
		return nil, fmt.Errorf("malformed request line")
	}
	method, uri, version := requestLine[0], requestLine[1], requestLine[2]

	// Parse the headers
	headers := make(map[string]string)
	i := 1 // First line after the request line
	for ; i < len(lines); i++ {
		line := lines[i]
		if line == "" { // A blank line separates headers and body
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			headers[key] = value
		}
	}

	// Parse the body (if present)
	body := strings.Join(lines[i+1:], "\r\n")

	return &HttpRequest{
		Method:  method,
		URI:     uri,
		Version: version,
		Headers: headers,
		Body:    body,
	}, nil
}
