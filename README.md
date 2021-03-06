# fake-accession-number-apis

Store and retrieve public-facing object IDs and URIs for accession numbers derived from cultural heriage open data releases.

## Background

The package provides an adjacent service to the `sfomuseum/accession-numbers` package whereby public-facing object IDs and URLs for individual online object records can be derived from an accession number. Many online collections allow you to search for an accession number but do not allow an object (web page) to be retrieved using only an accession number. This package provides services to store and retrieve the public object ID associated with an accession number.

_This is work in progress. Although the basic scaffolding is complete things may still change, specifically whether and how this package can be updated to use the [data defintion files](https://github.com/sfomuseum/accession-numbers/tree/main/data) in the `sfomuseum/accession-numbers` package._

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
#### Automatic redirects

It is also possible to have the server automatically redirect a matching accession number to its institution-specific object URL.

For example this request for accession number `2017.59` in the Whitney Museum of American Art collection will issue an HTTP redirect pointing to the webpage for Paul Mpagi Sepuya's [Self-Portrait Study with Roses at Night](https://whitney.org/collection/works/55448).

```
$> curl -s -I 'http://localhost:8080/redirect/?source-uri=whitneymuseum://&accession-number=2017.59'

HTTP/1.1 302 Found
Content-Type: text/html; charset=utf-8
Location: https://whitney.org/collection/works/55448
Date: Sun, 19 Dec 2021 21:41:19 GMT
```

## Models

### Databases

There is a common interface for storing accession number data, defined in the [database/database.go](database/database.go) file.

As of this writing there is a implementation of that interface for any package that supports the `database/sql` interface although only the `mattn/go-sqlite3` package is imported in the code.

### Sources

There is a common interface for accession number data sources, defined in the [source/source.go](source/source.go) file.

As of this writing the schemes and URIs used to define sources are different from those defined in the [sfomuseum/accession-numbers schema](https://github.com/sfomuseum/accession-numbers/blob/main/schema/definition.schema.json). These should be reconciled, if possible.

## Sources

The following data sources are supported:

### ACMI (acmi://)

* https://github.com/ACMILabs/acmi-api
* [source/acmi.go](source/acmi.go)

#### Importing data

```
$> ./bin/import \
	-database-uri 'sql://sqlite3?dsn=accessionumbers.db' \
	-source-uri 'acmi://' \
	/usr/local/data/acmi-api/app/tsv/works.tsv
```

Note: This package uses the `tsv/works.tsv` data file which _appears_ to have some encoding issues so not all records are able to be imported at this time.

#### Resolving object URLs

```
$> curl -s -I 'http://localhost:8080/redirect/?source-uri=acmi://&accession-number=X001564'

HTTP/1.1 302 Found
Content-Type: text/html; charset=utf-8
Location: https://www.acmi.net.au/works/116498/
Date: Mon, 20 Dec 2021 01:36:29 GMT
```

### Art Institute of Chicago (artic://)

* https://github.com/art-institute-of-chicago/api-data
* https://artic-api-data.s3.amazonaws.com/artic-api-data.tar.bz2
* [source/artic.go](source/artic.go)

#### Importing data

```
$> ./bin/import \
	-database-uri 'sql://sqlite3?dsn=accessionumbers.db' \
	-source-uri artic:// \
	/usr/local/data/artic-api-data/json/artworks/
```

#### Resolving object URLs

```
$> curl -s -I 'http://localhost:8080/redirect/?source-uri=artic://&accession-number=1982.2072'

HTTP/1.1 302 Found
Content-Type: text/html; charset=utf-8
Location: https://www.artic.edu/artworks/99990/
Date: Tue, 21 Dec 2021 06:43:08 GMT
```

### Metropolitan Museum of Art (metmuseum://)

* https://github.com/metmuseum/openaccess
* [source/metmuseum.go](source/metmuseum.go)

#### Importing data

```
$> ./bin/import \
	-database-uri 'sql://sqlite3?dsn=accessionumbers.db' \
	-source-uri 'metmuseum://?remove_bom=1' \
	/usr/local/data/openaccess/MetObjects.csv
```

Note: As of December, 2021 the Metropolitan Museum of Art openaccess CSV file contains a leading [byte order mark](https://en.wikipedia.org/wiki/Byte_order_mark) (BOM). In order to account for this you will need to explicitly pass a `?remove_bom=1` parameter when defining a `metmuseum://` source URI. The hope is that eventually the BOM will be removed from the published data making the flag unnecessary.

### Minneapolis Institute of Art (artsmia://)

* https://github.com/artsmia/collection
* [source/artmia.go](source/artsmia.go)

#### Importing data

```
./bin/import \
	-database-uri 'sql://sqlite3?dsn=accessionumbers.db' \
	-source-uri artsmia:// \
	/usr/local/data/collection/objects/
```

#### Resolving object URLs

```
$> curl -s -I 'http://localhost:8080/redirect/?source-uri=artsmia://&accession-number=85.34'

HTTP/1.1 302 Found
Content-Type: text/html; charset=utf-8
Location: https://collections.artsmia.org/art/3344/
Date: Mon, 20 Dec 2021 06:56:23 GMT
```

### Museum of Modern Art (moma://)

* https://github.com/MuseumofModernArt/collection
* [source/moma.go](source/moma.go)

#### Importing data

```
$> ./bin/import \
	-database-uri 'sql://sqlite3?dsn=accessionumbers.db' \
	-source-uri 'moma://' \
	/usr/local/data/collection/Artworks.csv
```

### National Gallery of Art (nga://)

* https://github.com/NationalGalleryOfArt/opendata
* [source/nga.go](source/nga.go)

#### Importing data

```
$> ./bin/import \
	-database-uri 'sql://sqlite3?dsn=accessionumbers.db' \
	-source-uri 'nga://' \
	/usr/local/data/opendata/data/objects.csv
```

### Whitney Museum of American Art

* https://github.com/whitneymuseum/open-access
* [source/whitneymuseum.go](source/whitneymuseum.go)

#### Importing data

```
$> ./bin/import \
	-database-uri 'sql://sqlite3?dsn=accessionumbers.db' \
	-source-uri 'whitneymuseum://' \
	/usr/local/data/open-access/artworks.csv
```

## See also

* https://github.com/sfomuseum/accession-numbers