package main

import (
	"log"

	"github.com/laiboonh/proglog/internal/server"
)

func main() {
	s := server.NewHttpServer(":8080")
	log.Fatal(s.ListenAndServe())
}
