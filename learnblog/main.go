package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/google/uuid"
)

var addrFlag = flag.String("addr", "0.0.0.0:8080", "server bind address in form IP:PORT")
var command = flag.String("command", "serve", "command to run: serve,uuid,migrate")
var databaseURLParam = flag.String("database", "", "postgres database URL")
var migrationsPathParam = flag.String("migrations", "", "path to migration scripts")

func main() {
	flag.Parse()
	switch *command {
	case "serve":
		runServer()
	case "uuid":
		generateUUID()
	case "migrate-up":
		err := migrateDbUp(*databaseURLParam, *migrationsPathParam)
		if err != nil {
			log.Fatal(err)
		}
	case "migrate-down":
		err := migrateDbDown(*databaseURLParam, *migrationsPathParam)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unsupported command: %s", *command)
	}
}

func validateDatabaseParam() error {
	_, err := url.Parse(*databaseURLParam)
	if err != nil {
		return fmt.Errorf("failed to parse database URL: %w", err)
	}
	return nil
}

func generateUUID() {
	fmt.Println(uuid.NewString())
}

func runServer() {
	err := validateDatabaseParam()
	if err != nil {
		log.Fatal(err)
	}

	dsn, err := postgresURL2DSN(*databaseURLParam)
	if err != nil {
		log.Fatal(err)
	}
	db, err := openGorm(dsn)
	if err != nil {
		log.Fatal(err)
	}
	blogStore := BlogPostStoreImpl{DB: db}
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
	err = srv.Run()
	if err != nil {
		log.Printf("server ended: %s", err)
	}
}
