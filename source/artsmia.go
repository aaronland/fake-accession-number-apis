package source

// https://github.com/artsmia/collection

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aaronland/fake-accession-number-apis/database"
	"github.com/jtacoma/uritemplates"
	"io"
	"io/fs"
	_ "log"
	"os"
	"path/filepath"
)

const ARTSMIA_ORGANIZATION_SCHEME string = "artsmia"
const ARTSMIA_OBJECT_TEMPLATE = "https://collections.artsmia.org/art/{objectid}/"

type ArtsMIASource struct {
	Source
	object_template *uritemplates.UriTemplate
}

func init() {
	ctx := context.Background()
	RegisterSource(ctx, ARTSMIA_ORGANIZATION_SCHEME, NewArtsMIASource)
}

func NewArtsMIASource(ctx context.Context, uri string) (Source, error) {

	t, err := uritemplates.Parse(ARTSMIA_OBJECT_TEMPLATE)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse object template, %w", err)
	}

	s := &ArtsMIASource{
		object_template: t,
	}

	return s, nil
}

func (s *ArtsMIASource) Import(ctx context.Context, db database.AccessionNumberDatabase, uris ...string) error {

	for _, u := range uris {

		err := s.importURI(ctx, db, u)

		if err != nil {
			return fmt.Errorf("Failed to import URI '%s', %w", u, err)
		}
	}

	return nil
}

func (s *ArtsMIASource) ObjectURI(ctx context.Context, acc *database.AccessionNumber) (string, error) {

	values := map[string]interface{}{
		"objectid": acc.ObjectId,
	}

	return s.object_template.Expand(values)
}

func (s *ArtsMIASource) importURI(ctx context.Context, db database.AccessionNumberDatabase, u string) error {

	root := os.DirFS(u)

	err := fs.WalkDir(root, ".", func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return fmt.Errorf("Failed to walk dir, %w", err)
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()

		if err != nil {
			return fmt.Errorf("Failed to derive info for '%s', %w", path, err)
		}

		if info.Size() == 0 {
			return nil
		}

		fq_path := filepath.Join(u, path)

		fh, err := os.Open(fq_path)

		if err != nil {
			return fmt.Errorf("Failed to open '%s', %w", fq_path, err)
		}

		defer fh.Close()

		err = s.importReader(ctx, db, fh)

		if err != nil {
			return fmt.Errorf("Failed to import '%s', %w", fq_path, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Failed to walk root '%s', %w", u, err)
	}

	return nil
}

func (s *ArtsMIASource) importReader(ctx context.Context, db database.AccessionNumberDatabase, r io.Reader) error {

	var row map[string]interface{}

	dec := json.NewDecoder(r)
	err := dec.Decode(&row)

	if err != nil {
		return fmt.Errorf("Failed to decode row, %w", err)
	}

	object_url, ok := row["id"]

	if !ok {
		return fmt.Errorf("Row is missing id column")
	}

	accession_number, ok := row["accession_number"]

	if !ok {
		return fmt.Errorf("Row is missing accession_number column")
	}

	object_id := filepath.Base(object_url.(string))

	org_uri := fmt.Sprintf("%s://", ARTSMIA_ORGANIZATION_SCHEME)

	a := &database.AccessionNumber{
		AccessionNumber: accession_number.(string),
		ObjectId:        object_id,
		OrganizationURI: org_uri, // update to use https://github.com/sfomuseum/accession-numbers/blob/main/data/nga.gov.json
	}

	err = db.AddAccessionNumber(ctx, a)

	if err != nil {
		return fmt.Errorf("Failed to add accession number '%s', %w", accession_number, err)
	}

	return nil
}
