# fake-accession-number-apis

Store and retrieve public-facing object IDs and URIs for accession numbers derived from cultural heriage open data releases.

## Background

The package provides an adjacent service to the `sfomuseum/accession-numbers` package whereby public-facing object IDs and URLs for individual online object records can be derived from an accession number. Many online collections allow you to search for an accession number but do not allow an object (web page) to be retrieved using only an accession number. This package provides services to store and retrieve the public object ID associated with an accession number.

_This is work in progress. Although the basic scaffolding is complete things may still change, specifically whether and how this package can be updated to use the [data defintion files](https://github.com/sfomuseum/accession-numbers/tree/main/data) in the `sfomuseum/accession-numbers` package.

## Documentation

Documentation is incomplete at this time.

## Example

### Building the tools

```
$> make cli
go build -mod vendor -o bin/import cmd/import/main.go
go build -mod vendor -o bin/lookup cmd/lookup/main.go
go build -mod vendor -o bin/server cmd/server/main.go
```

### Importing data

Here is how you would import object IDs and accession numbers from the National Gallery of Art's (NGA) [opendata release](https://github.com/NationalGalleryOfArt/opendata):

```
$> bin/import \
	-database-uri 'sql://sqlite3?dsn=accessionumbers.db' \
	-source-uri nga:// \
	/usr/local/data/nga/opendata/data/objects.csv

# Time passes...

$> sqlite3 accessionumbers.db 
SQLite version 3.36.0 2021-06-18 18:58:49
Enter ".help" for usage hints.
sqlite> SELECT COUNT(object_id) FROM accession_numbers;
136612
```

### Looking up an accession number (from the command line)

Here is how you would look up the corresponding object ID, in the NGA's collection, for the accession number `1994.59.10`:

```
$> bin/lookup \
	-database-uri 'sql://sqlite3?dsn=accessionumbers.db' \
	-source-uri nga:// \
	1994.59.10
	
89682
```

As in: https://www.nga.gov/collection/art-object-page.89682.html

### Looking up an accession number (via an HTTP API)

First start the `server` tool:

```
$> ./bin/server -database-uri 'sql://sqlite3?dsn==accessionumbers.db'
2021/11/30 12:20:36 Listening for requests on http://localhost:8080
```

And then query the server with a source and accession number:

```
$> curl -s 'http://localhost:8080/?source-uri=nga://&accession-number=1994.59.10' | jq
{
  "accession_number": "1994.59.10",
  "object_id": "89682",
  "organization_uri": "nga://"
}
```

## Models

### Databases

There is an common interface for storing accession number data, defined in the [database/database.go](database/database.go) file.

As of this writing there is a implementation of that interface for any package that supports the `database/sql` interface although only the `mattn/go-sqlite3` package is imported in the code.

### Sources

There is an common interface for accession number data sources, defined in the [source/source.go](source/source.go) file.

As of this writing the schemes and URIs used to define sources are different from those defined in the [sfomuseum/accession-numbers schema](https://github.com/sfomuseum/accession-numbers/blob/main/schema/definition.schema.json). These should be reconciled, if possible.

## Sources

The following data sources are supported:

### Metropolitan Museum of Art (metmuseum://)

* [source/metmuseum.go](source/metmuseum.go)

For example:

```
$> ./bin/import \
	-database-uri 'sql://sqlite3?dsn=accessionumbers.db' \
	-source-uri 'metmuseum://?remove_bom=1' \
	/usr/local/data/openaccess/MetObjects.csv
```

Note: As of December, 2021 the Metropolitan Museum of Art openaccess CSV file contains a leading [byte order mark](https://en.wikipedia.org/wiki/Byte_order_mark) (BOM). In order to account for this you will need to explicitly pass a `?remove_bom=1` parameter when defining a `metmuseum://` source URI. The hope is that eventually the BOM will be removed from the published data making the flag unnecessary.

### National Gallery of Art (nga://)

* [source/nga.go](source/nga.go)

```
$> ./bin/import \
	-database-uri 'sql://sqlite3?dsn=accessionumbers.db' \
	-source-uri 'nga://' \
	/usr/local/data/opendata/data/objects.csv
```

### Whitney Museum of American Art

* [source/whitneymuseum.go](source/whitneymuseum.go)

For example:

```
$> ./bin/import \
	-database-uri 'sql://sqlite3?dsn=accessionumbers.db' \
	-source-uri 'whitneymuseum://' \
	/usr/local/data/open-access/artworks.csv
```

## See also

* https://github.com/sfomuseum/accession-numbers