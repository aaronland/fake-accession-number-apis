package source

// https://github.com/NationalGalleryOfArt/opendata

import (
	"context"
	_ "encoding/json"
	"fmt"
	"github.com/aaronland/fake-accession-number-apis/database"
	"github.com/jtacoma/uritemplates"
	"github.com/sfomuseum/go-csvdict"
	"io"
	"os"
)

const WCMA_ORGANIZATION_SCHEME string = "wcma"
const WCMA_OBJECT_TEMPLATE = "http://egallery.williams.edu/objects/{objectid}/"

type WCMASource struct {
	Source
	object_template *uritemplates.UriTemplate
}

func init() {
	ctx := context.Background()
	RegisterSource(ctx, WCMA_ORGANIZATION_SCHEME, NewWCMASource)
}

func NewWCMASource(ctx context.Context, uri string) (Source, error) {

	t, err := uritemplates.Parse(WCMA_OBJECT_TEMPLATE)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse object template, %w", err)
	}

	s := &WCMASource{
		object_template: t,
	}

	return s, nil
}

func (s *WCMASource) Import(ctx context.Context, db database.AccessionNumberDatabase, uris ...string) error {

	for _, u := range uris {

		err := s.importURI(ctx, db, u)

		if err != nil {
			return fmt.Errorf("Failed to import URI '%s', %w", u, err)
		}
	}

	return nil
}

func (s *WCMASource) ObjectURI(ctx context.Context, acc *database.AccessionNumber) (string, error) {

	values := map[string]interface{}{
		"objectid": acc.ObjectId,
	}

	return s.object_template.Expand(values)
}

func (s *WCMASource) importURI(ctx context.Context, db database.AccessionNumberDatabase, u string) error {

	fh, err := os.Open(u)

	if err != nil {
		return fmt.Errorf("Failed to open '%s', %w", u, err)
	}

	defer fh.Close()

	/*
		var collection []map[string]string

		dec := json.NewDecoder(fh)
		err = dec.Decode(&collection)

		if err != nil {
			return fmt.Errorf("Failed to decode '%s', %w", u, err)
		}

		for _, row := range collection {
	*/

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
			return fmt.Errorf("Row is missing id column")
		}

		accession_number, ok := row["accession_number"]

		if !ok {
			return fmt.Errorf("Row is missing accession_number column")
		}

		org_uri := fmt.Sprintf("%s://", WCMA_ORGANIZATION_SCHEME)

		a := &database.AccessionNumber{
			AccessionNumber: accession_number,
			ObjectId:        object_id,
			OrganizationURI: org_uri, // update to use https://github.com/sfomuseum/accession-numbers/blob/main/data/wcma.gov.json
		}

		err = db.AddAccessionNumber(ctx, a)

		if err != nil {
			return fmt.Errorf("Failed to add accession number '%s', %w", accession_number, err)
		}
	}

	return nil
}
