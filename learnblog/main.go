package main

import (
	"flag"
	"log"
)

var addrFlag = flag.String("addr", "0.0.0.0:8080", "server bind address in form IP:PORT")

func main() {
	flag.Parse()
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
