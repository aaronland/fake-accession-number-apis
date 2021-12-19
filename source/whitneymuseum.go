package source

// https://github.com/whitneymuseum/open-access

import (
	"context"
	"fmt"
	"github.com/aaronland/fake-accession-number-apis/database"
	"github.com/jtacoma/uritemplates"
	"github.com/sfomuseum/go-csvdict"
	"io"
	"os"
)

const WHITNEYMUSEUM_ORGANIZATION_SCHEME string = "whitneymuseum"
const WHITNEYMUSEUM_OBJECT_TEMPLATE = "https://whitney.org/collection/works/{objectid}"

type WhitneyMuseumSource struct {
	Source
	object_template *uritemplates.UriTemplate
}

func init() {
	ctx := context.Background()
	RegisterSource(ctx, WHITNEYMUSEUM_ORGANIZATION_SCHEME, NewWhitneyMuseumSource)
}

func NewWhitneyMuseumSource(ctx context.Context, uri string) (Source, error) {

	t, err := uritemplates.Parse(WHITNEYMUSEUM_OBJECT_TEMPLATE)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse object template, %w", err)
	}

	s := &WhitneyMuseumSource{
		object_template: t,
	}

	return s, nil
}

func (s *WhitneyMuseumSource) Import(ctx context.Context, db database.AccessionNumberDatabase, uris ...string) error {

	for _, u := range uris {

		err := s.importURI(ctx, db, u)

		if err != nil {
			return fmt.Errorf("Failed to import URI '%s', %w", u, err)
		}
	}

	return nil
}

func (s *WhitneyMuseumSource) ObjectURI(ctx context.Context, acc database.AccessionNumber) (string, error) {

	values := map[string]interface{}{
		"objectid": acc.ObjectId,
	}

	return s.object_template.Expand(values)
}

func (s *WhitneyMuseumSource) importURI(ctx context.Context, db database.AccessionNumberDatabase, u string) error {

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

		object_id, ok := row["id"]

		if !ok {
			return fmt.Errorf("Row is missing objectid column")
		}

		accession_number, ok := row["accession_number"]

		if !ok {
			return fmt.Errorf("Row is missing accessionnum column")
		}

		org_uri := fmt.Sprintf("%s://", WHITNEYMUSEUM_ORGANIZATION_SCHEME)

		a := &database.AccessionNumber{
			AccessionNumber: accession_number,
			ObjectId:        object_id,
			OrganizationURI: org_uri, // update to use https://github.com/sfomuseum/accession-numbers/blob/main/data/nga.gov.json
		}

		err = db.AddAccessionNumber(ctx, a)

		if err != nil {
			return fmt.Errorf("Failed to add accession number '%s', %w", accession_number, err)
		}
	}

	return nil
}
