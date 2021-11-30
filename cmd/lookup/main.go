package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/fake-accession-number-apis/database"
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

	accession_numbers := flag.Args()

	for _, n := range accession_numbers {

		a, err := db.GetByAccessionNumber(ctx, *source_uri, n)

		if err != nil {
			log.Fatalf("Failed to retrieve accession number '%s', %v", n, err)
		}

		fmt.Println(a.ObjectId)
	}

}
