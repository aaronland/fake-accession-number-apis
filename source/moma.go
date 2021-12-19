package source

// https://github.com/MuseumofModernArt/collection

import (
	"context"
	"fmt"
	"github.com/aaronland/fake-accession-number-apis/database"
	"github.com/jtacoma/uritemplates"
	"github.com/sfomuseum/go-csvdict"
	"io"
	"os"
)

const MOMA_ORGANIZATION_SCHEME string = "moma"
const MOMA_OBJECT_TEMPLATE string = "https://www.moma.org/collection/works/{objectid}"

type MoMASource struct {
	Source
	object_template *uritemplates.UriTemplate
}

func init() {
	ctx := context.Background()
	RegisterSource(ctx, MOMA_ORGANIZATION_SCHEME, NewMoMASource)
}

func NewMoMASource(ctx context.Context, uri string) (Source, error) {

	t, err := uritemplates.Parse(MOMA_OBJECT_TEMPLATE)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse object template, %w", err)
	}

	s := &MoMASource{
		object_template: t,
	}

	return s, nil
}

func (s *MoMASource) Import(ctx context.Context, db database.AccessionNumberDatabase, uris ...string) error {

	for _, u := range uris {

		err := s.importURI(ctx, db, u)

		if err != nil {
			return fmt.Errorf("Failed to import URI '%s', %w", u, err)
		}
	}

	return nil
}

func (s *MoMASource) ObjectURI(ctx context.Context, acc *database.AccessionNumber) (string, error) {

	values := map[string]interface{}{
		"objectid": acc.ObjectId,
	}

	return s.object_template.Expand(values)
}

func (s *MoMASource) importURI(ctx context.Context, db database.AccessionNumberDatabase, u string) error {

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

		object_id, ok := row["ObjectID"]

		if !ok {
			return fmt.Errorf("Row is missing ObjectID column")
		}

		accession_number, ok := row["AccessionNumber"]

		if !ok {
			return fmt.Errorf("Row is missing AccessionNumber column")
		}

		org_uri := fmt.Sprintf("%s://", MOMA_ORGANIZATION_SCHEME)

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
