package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/google/uuid"
)

var addrFlag = flag.String("addr", "0.0.0.0:8080", "server bind address in form IP:PORT")
var command = flag.String("command", "serve", "command to run: serve,uuid")

func main() {
	flag.Parse()
	switch *command {
	case "serve":
		runServer()
	case "uuid":
		generateUUID()
	default:
		log.Fatalf("unsupported command: %s", *command)
	}
}

func generateUUID() {
	fmt.Println(uuid.NewString())
}

func runServer() {
	blogStore := &BlogStore{}
	blogService := BlogService{
		Store: blogStore,
	}
	blogRouter := BlogRouter{
		Service: blogService,
	}
	srv := &Server{
		RouterPackages: []RouterPackage{
			{
				BasePath: "blog",
				Router:   blogRouter,
			},
		},
		Addr: *addrFlag,
	}
	err := srv.Run()
	if err != nil {
		log.Printf("server ended: %s", err)
	}
}
