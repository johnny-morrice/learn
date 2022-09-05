package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/golang-migrate/migrate/v4"
	migratedb "github.com/golang-migrate/migrate/v4/database"
	pgmigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const dbDriverName = "postgres"

var addrFlag = flag.String("addr", "0.0.0.0:8080", "server bind address in form IP:PORT")
var command = flag.String("command", "serve", "command to run: serve,uuid,migrate")
var databaseURL = flag.String("database", "", "postgres database URL")
var migrationsPath = flag.String("migrations", "", "path to migration scripts")

func main() {
	flag.Parse()
	switch *command {
	case "serve":
		runServer()
	case "uuid":
		generateUUID()
	case "migrate":
		migrateDbUp()
	default:
		log.Fatalf("unsupported command: %s", *command)
	}
}

func validateDatabaseParam() error {
	_, err := url.Parse(*databaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse database URL: %w", err)
	}
	return nil
}

func openDb() (migratedb.Driver, error) {
	err := validateDatabaseParam()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(dbDriverName, *databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error opening postgres connection: %w", err)
	}
	driver, err := pgmigrate.WithInstance(db, &pgmigrate.Config{})
	if err != nil {
		return nil, fmt.Errorf("error creating migration driver: %w", err)
	}
	return driver, nil
}

func migrateDbUp() {
	driver, err := openDb()
	if err != nil {
		log.Fatal(err)
		return
	}
	m, err := migrate.NewWithDatabaseInstance(*migrationsPath, dbDriverName, driver)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = m.Up()
	if err != nil {
		log.Fatal(err)
		return
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
