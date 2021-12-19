package source

// https://github.com/metmuseum/openaccess

import (
	"context"
	"fmt"
	"github.com/aaronland/fake-accession-number-apis/database"
	"github.com/jtacoma/uritemplates"
	"github.com/sfomuseum/go-csvdict"
	"io"
	"net/url"
	"os"
	"strconv"
)

const METMUSEUM_ORGANIZATION_SCHEME string = "metmuseum"
const METMUSEUM_OBJECT_TEMPLATE string = "http://www.metmuseum.org/art/collection/search/{objectid}"

type MetMuseumSource struct {
	Source
	remove_bom      bool
	object_template *uritemplates.UriTemplate
}

func init() {
	ctx := context.Background()
	RegisterSource(ctx, METMUSEUM_ORGANIZATION_SCHEME, NewMetMuseumSource)
}

func NewMetMuseumSource(ctx context.Context, uri string) (Source, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	q := u.Query()

	// As of December, 2021 Metmuseum openaccess CSV contains a BOM
	// so we need to strip it in order to read column names correctly

	str_remove := q.Get("remove_bom")

	var remove_bom bool

	if str_remove != "" {

		remove, err := strconv.ParseBool(str_remove)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse ?remove_bom parameter, %w", err)
		}

		remove_bom = remove
	}

	t, err := uritemplates.Parse(METMUSEUM_OBJECT_TEMPLATE)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse object template, %w", err)
	}

	s := &MetMuseumSource{
		remove_bom:      remove_bom,
		object_template: t,
	}

	return s, nil
}

func (s *MetMuseumSource) Import(ctx context.Context, db database.AccessionNumberDatabase, uris ...string) error {

	for _, u := range uris {

		err := s.importURI(ctx, db, u)

		if err != nil {
			return fmt.Errorf("Failed to import URI '%s', %w", u, err)
		}
	}

	return nil
}

func (s *MetMuseumSource) ObjectURI(ctx context.Context, acc *database.AccessionNumber) (string, error) {

	values := map[string]interface{}{
		"objectid": acc.ObjectId,
	}

	return s.object_template.Expand(values)
}

func (s *MetMuseumSource) importURI(ctx context.Context, db database.AccessionNumberDatabase, u string) error {

	fh, err := os.Open(u)

	if err != nil {
		return fmt.Errorf("Failed to open '%s', %w", u, err)
	}

	defer fh.Close()

	// See notes above

	if s.remove_bom {
		fh.Seek(3, 0)
	}

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

		object_id, ok := row["Object ID"]

		if !ok {
			return fmt.Errorf("Row is missing Object ID column")
		}

		accession_number, ok := row["Object Number"]

		if !ok {
			return fmt.Errorf("Row is missing Object Number column")
		}

		org_uri := fmt.Sprintf("%s://", METMUSEUM_ORGANIZATION_SCHEME)

		a := &database.AccessionNumber{
			AccessionNumber: accession_number,
			ObjectId:        object_id,
			OrganizationURI: org_uri, // update to use https://github.com/sfomuseum/accession-numbers/blob/main/data/metmuseum.org.json
		}

		err = db.AddAccessionNumber(ctx, a)

		if err != nil {
			return fmt.Errorf("Failed to add accession number '%s', %w", accession_number, err)
		}
	}

	return nil
}
