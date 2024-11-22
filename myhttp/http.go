package myhttp

import (
	"fmt"
	"log"
	"net"
	"strings"
)

// Definir los métodos HTTP como constantes string
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
	StatusLine string            // Línea de estado: incluye versión, código de estado y mensaje.
	Headers    map[string]string // Encabezados: clave-valor.
	Body       []byte            // Cuerpo: contenido de la respuesta (puede ser HTML, JSON, etc.).
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
		return &HttpResponse{
			StatusLine: "HTTP/1.1 404 Not Found",
			Headers:    map[string]string{"Content-Type": "text/plain"},
			Body:       []byte("404 Not Found"),
		}
	}

	response := &HttpResponse{
		Headers: make(map[string]string),
	}
	handler(req, response)

	return response
}

func (http *Server) HttpServer(address string) {
	// Crear el listener en el puerto 6969
	listener, err := net.Listen("tcp4", address)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}

	defer listener.Close()
	log.Println("Servidor escuchando en %s", address)

	for {
		// Aceptar conexiones entrantes
		conn, err := listener.Accept()
		log.Printf("Nueva conexión de: %s", conn.RemoteAddr())

		if err != nil {
			log.Printf("Error al aceptar conexión: %v", err)
			continue
		}

		// Manejar la conexión en una goroutine separada
		go handleConnection(conn, http)
	}
}

func handleConnection(conn net.Conn, server *Server) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	// Leer datos enviados por el cliente
	n, err := conn.Read(buffer)
	if err != nil {
		log.Printf("Error al leer datos del cliente: %v", err)
		return
	}

	// Registrar la solicitud recibida
	clientRequest := string(buffer[:n])
	log.Printf("Solicitud recibida:\n%s", clientRequest)

	request, err := parseHttpRequest(clientRequest)

	if err != nil {
		log.Printf("Error al parsear los datos: %v", err)
	}

	response := server.Serve(request)

	// Enviar la respuesta al cliente
	_, err = conn.Write([]byte(response))
	if err != nil {
		log.Printf("Error al enviar respuesta al cliente: %v", err)
	}

	log.Println("Respuesta enviada al cliente.")
}

func parseHttpRequest(request string) (*HttpRequest, error) {
	lines := strings.Split(request, "\r\n")
	if len(lines) < 1 {
		return nil, fmt.Errorf("solicitud mal formada")
	}

	// Parsear la línea inicial (Método, URI, Versión)
	requestLine := strings.Fields(lines[0])
	if len(requestLine) < 3 {
		return nil, fmt.Errorf("línea de solicitud mal formada")
	}
	method, uri, version := requestLine[0], requestLine[1], requestLine[2]

	// Parsear los encabezados
	headers := make(map[string]string)
	i := 1 // Primera línea después de la línea de solicitud
	for ; i < len(lines); i++ {
		line := lines[i]
		if line == "" { // Línea en blanco separa encabezados y cuerpo
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			headers[key] = value
		}
	}

	// Parsear el cuerpo (si existe)
	body := strings.Join(lines[i+1:], "\r\n")

	return &HttpRequest{
		Method:  method,
		URI:     uri,
		Version: version,
		Headers: headers,
		Body:    body,
	}, nil
}
