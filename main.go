package main

import "myhttpserver/myhttp"

func main() {

	http := myhttp.NewServer()

	http.HandleFunction(myhttp.GET, "/hello", func(req *myhttp.HttpRequest, res *myhttp.HttpResponse) {
		res.Headers["Content-Type"] = "text/plain"
		res.Body = []byte("Hola, mundo prueba handler!")
	})

	http.HandleFunction(myhttp.POST, "/hello", func(req *myhttp.HttpRequest, res *myhttp.HttpResponse) {
		res.Headers["Content-Type"] = "text/plain"
		res.Body = []byte("Hola, mundo prueba handler METODO POST!")
	})

	http.HandleFunction(myhttp.GET, "/adios", func(req *myhttp.HttpRequest, res *myhttp.HttpResponse) {
		res.Headers["Content-Type"] = "text/plain"
		res.Body = []byte("adios, mundo prueba handler!")
	})

	http.HandleFunction(myhttp.GET, "/paginaProtegida", Middleware(RutaProtegidaHandler()))

	http.HttpServer("0.0.0.0:6969")
}

func Middleware(next myhttp.HandleFunc) myhttp.HandleFunc {
	return func(req *myhttp.HttpRequest, res *myhttp.HttpResponse) {
		next.ServeHTTP(req, res)
	}

}

func RutaProtegidaHandler() myhttp.HandleFunc {
	return func(req *myhttp.HttpRequest, res *myhttp.HttpResponse) {
		res.Headers["Content-Type"] = "text/plain"
		res.Body = []byte("Ruta protegida!")
	}
}
