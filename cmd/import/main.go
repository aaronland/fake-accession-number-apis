package main

import (
	"context"
	"flag"
	"github.com/aaronland/fake-accession-number-apis/database"
	"github.com/aaronland/fake-accession-number-apis/source"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {

	database_uri := flag.String("database-uri", "", "...")
	source_uri := flag.String("source-uri", "", "...")

	flag.Parse()

	ctx := context.Background()

	db, err := database.NewDatabase(ctx, *database_uri)

	if err != nil {
		log.Fatalf("Failed to create database, %v", err)
	}

	src, err := source.NewSource(ctx, *source_uri)

	if err != nil {
		log.Fatalf("Failed to create database, %v", err)
	}

	sources := flag.Args()

	err = src.Import(ctx, db, sources...)

	if err != nil {
		log.Fatalf("Failed to import, %v", err)
	}
}
