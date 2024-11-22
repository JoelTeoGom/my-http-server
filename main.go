package main

import "myhttpserver/myhttp"

func main() {

	http := myhttp.NewServer()

	http.HandleFunction(myhttp.GET, "/hello", func(req *myhttp.HttpRequest, res *myhttp.HttpResponse) {
		res.Headers["Content-Type"] = "text/plain"
		res.Body = []byte("Hola, mundo prueba handler!")
	})

	http.HttpServer("0.0.0.0:6969")
}
