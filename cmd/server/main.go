package main

import (
	"context"
	"github.com/aaronland/fake-accession-number-apis/api"
	"github.com/aaronland/fake-accession-number-apis/database"
	"github.com/aaronland/go-http-server"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"github.com/sfomuseum/go-flags/flagset"
	"log"
	"net/http"
)

func main() {

	fs := flagset.NewFlagSet("server")

	database_uri := fs.String("database-uri", "", "...")
	server_uri := fs.String("server-uri", "http://localhost:8080", "...")

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "SERVER")

	if err != nil {
		log.Fatalf("Failed to set flags from environment variables, %v", err)
	}

	ctx := context.Background()

	db, err := database.NewDatabase(ctx, *database_uri)

	if err != nil {
		log.Fatalf("Failed to create database, %v", err)
	}

	mux := http.NewServeMux()

	redirect_handler, err := api.NewRedirectHandler(db)

	if err != nil {
		log.Fatalf("Failed to create new redirect handler, %v", err)
	}

	mux.Handle("/redirect/", redirect_handler)

	lookup_handler, err := api.NewLookupHandler(db)

	if err != nil {
		log.Fatalf("Failed to create new lookup handler, %v", err)
	}

	lookup_handler = cors.Default().Handler(lookup_handler)

	mux.Handle("/", lookup_handler)

	s, err := server.NewServer(ctx, *server_uri)

	if err != nil {
		log.Fatalf("Failed to create new server, %v", err)
	}

	log.Printf("Listening for requests on %s\n", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		log.Fatalf("Failed to serve requests, %v", err)
	}
}
