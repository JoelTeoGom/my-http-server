package main

import (
	"fmt"
	"log"
	"net"
	"strings"
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

func main() {
	// Crear el listener en el puerto 6969
	listener, err := net.Listen("tcp4", "0.0.0.0:6969")
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}

	defer listener.Close()
	log.Println("Servidor escuchando en 0.0.0.0:6969")

	for {
		// Aceptar conexiones entrantes
		conn, err := listener.Accept()
		log.Printf("Nueva conexión de: %s", conn.RemoteAddr())

		if err != nil {
			log.Printf("Error al aceptar conexión: %v", err)
			continue
		}

		// Manejar la conexión en una goroutine separada
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
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

	fmt.Printf("Received request: %+v\n", request)
	// Preparar una respuesta HTTP válida
	responseBody := "hola mundo!"
	response := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\n"+ // Línea de estado
			"Content-Type: text/html\r\n"+ // Encabezado: tipo de contenido
			"Content-Length: %d\r\n"+ // Encabezado: longitud del contenido
			"\r\n%s", // Separador y cuerpo
		len(responseBody), responseBody,
	)

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
