# MyHTTPServer

## Description
This project implements a basic HTTP server written in Go. The server architecture allows handling specific routes with different HTTP methods (GET, POST, etc.) and responds to requests with custom logic defined through handlers. It is a modular and extensible foundation for building HTTP servers in Go.

## Key Features
- **HTTP Methods Support**: Supports common HTTP methods like GET, POST, PUT, DELETE, PATCH, OPTIONS, and HEAD.
- **Custom Handlers**: Define custom logic for each route using handlers.
- **Routing**: Match HTTP requests to specific handlers based on method and path.
- **Middleware**: Add functionality like logging, authentication, or validation by wrapping handlers.
- **Concurrency**: Handles multiple client connections using Go’s lightweight goroutines.

## Project Structure
```
myhttp/
├── myhttp.go    # Core HTTP server logic and routing
main.go          # Entry point defining routes and starting the server
```

### `myhttp.go`
- **`HttpRequest`**: Represents an HTTP request, including method, URI, version, headers, and body.
- **`HttpResponse`**: Represents an HTTP response, including status line, headers, and body.
- **`HandleFunc`**: Type alias for request handlers (`func(req *HttpRequest, res *HttpResponse)`).
- **`Server`**:
  - `HandleFunction(method, path, handler)`: Registers a handler for a specific method and path.
  - `Serve(req)`: Routes the request to the appropriate handler.
  - `HttpServer(address)`: Starts the server on the specified address (e.g., `0.0.0.0:6969`).

### `main.go`
Defines routes and handlers using the `myhttp` server implementation.

## Example Usage
### Starting the Server
The server listens on `0.0.0.0:6969`. Example routes are defined in `main.go`:
```go
http := myhttp.NewServer()

http.HandleFunction(myhttp.GET, "/hello", func(req *myhttp.HttpRequest, res *myhttp.HttpResponse) {
    res.Headers["Content-Type"] = "text/plain"
    res.Body = []byte("Hello, world from GET handler!")
})

http.HandleFunction(myhttp.POST, "/hello", func(req *myhttp.HttpRequest, res *myhttp.HttpResponse) {
    res.Headers["Content-Type"] = "text/plain"
    res.Body = []byte("Hello, world from POST handler!")
})

http.HttpServer("0.0.0.0:6969")
```

### Example Requests
#### GET `/hello`
```bash
curl -X GET http://localhost:6969/hello
# Response:
# Hello, world from GET handler!
```

#### POST `/hello`
```bash
curl -X POST http://localhost:6969/hello
# Response:
# Hello, world from POST handler!
```

### Middleware Example
Middleware can wrap handlers to add functionality like authentication or logging:
```go
func Middleware(next myhttp.HandleFunc) myhttp.HandleFunc {
    return func(req *myhttp.HttpRequest, res *myhttp.HttpResponse) {
        log.Println("Request received:", req.URI)
        next(req, res)
    }
}

http.HandleFunction(myhttp.GET, "/protected", Middleware(func(req *myhttp.HttpRequest, res *myhttp.HttpResponse) {
    res.Headers["Content-Type"] = "text/plain"
    res.Body = []byte("Protected route accessed!")
}))
```

### Protected Route Example
#### GET `/protected`
```bash
curl -X GET http://localhost:6969/protected
# Response:
# Protected route accessed!
```

## Future Enhancements
- **Cookie Handling**: Add functionality to set and retrieve cookies.
- **Session Management**: Integrate session handling for user authentication.
- **Enhanced Middleware**: Add support for common middlewares like authentication, logging, and rate limiting.
- **HTTPS Support**: Add support for serving over TLS.
- **Error Handling**: Implement robust error handling for edge cases.

## Running the Project
1. Clone the repository.
2. Build and run the server:
   ```bash
   go run main.go
   ```
3. Test the server using tools like `curl` or a web browser.

---

