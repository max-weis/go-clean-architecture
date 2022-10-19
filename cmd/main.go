package main

import (
	"github.com/rs/cors"
	"log"
	"net/http"
	"webshop"
	"webshop/shop/boundary"
)

func main() {
	appContext := webshop.Initialize()

	log.Println("started webshop app")

	handler := boundary.Handler(appContext)
	handler = cors.Default().Handler(handler)
	log.Fatal(http.ListenAndServe(":8000", handler))
}
