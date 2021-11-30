package source

// https://github.com/NationalGalleryOfArt/opendata

import (
	"context"
	"fmt"
	"github.com/aaronland/fake-accession-number-apis/database"
	"github.com/sfomuseum/go-csvdict"
	"io"
	"os"
)

const NGA_ORGANIZATION_URI string = "https://www.nga.gov/"

type NGASource struct {
	Source
}

func init() {
	ctx := context.Background()
	RegisterSource(ctx, "nga", NewNGASource)
}

func NewNGASource(ctx context.Context, uri string) (Source, error) {

	s := &NGASource{}
	return s, nil
}

func (s *NGASource) Import(ctx context.Context, db database.AccessionNumberDatabase, uris ...string) error {

	for _, u := range uris {

		err := s.importURI(ctx, db, u)

		if err != nil {
			return fmt.Errorf("Failed to import URI '%s', %w", u, err)
		}
	}

	return nil
}

func (s *NGASource) importURI(ctx context.Context, db database.AccessionNumberDatabase, u string) error {

	fh, err := os.Open(u)

	if err != nil {
		return fmt.Errorf("Failed to open '%s', %w", u, err)
	}

	defer fh.Close()

	csv_r, err := csvdict.NewReader(fh)

	if err != nil {
		return fmt.Errorf("Failed to create CSV reader for '%s', %w", u, err)
	}

	for {

		row, err := csv_r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("Failed to read row, %w", err)
		}

		object_id, ok := row["objectid"]

		if !ok {
			return fmt.Errorf("Row is missing objectid column")
		}

		accession_number, ok := row["accessionnum"]

		if !ok {
			return fmt.Errorf("Row is missing accessionnum column")
		}

		a := &database.AccessionNumber{
			AccessionNumber: accession_number,
			ObjectId:        object_id,
			OrganizationURI: NGA_ORGANIZATION_URI, // update to use https://github.com/sfomuseum/accession-numbers/blob/main/data/nga.gov.json
		}

		err = db.AddAccessionNumber(ctx, a)

		if err != nil {
			return fmt.Errorf("Failed to add accession number '%s', %w", accession_number, err)
		}
	}

	return nil
}
