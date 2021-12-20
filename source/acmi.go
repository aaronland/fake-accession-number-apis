package source

// https://github.com/ACMILabs/acmi-api

import (
	"bufio"
	"context"
	"fmt"
	"github.com/aaronland/fake-accession-number-apis/database"
	"github.com/jtacoma/uritemplates"
	"io"
	_ "log"
	"os"
	"regexp"
	_ "strings"
)

const ACMI_ORGANIZATION_SCHEME string = "acmi"
const ACMI_OBJECT_TEMPLATE = "https://www.acmi.net.au/works/{objectid}/"

var re_identifiers *regexp.Regexp

type ACMISource struct {
	Source
	object_template *uritemplates.UriTemplate
}

func init() {
	ctx := context.Background()
	RegisterSource(ctx, ACMI_ORGANIZATION_SCHEME, NewACMISource)

	re_identifiers = regexp.MustCompile(`^(\d+)\t([0-9A-Z]+).*`)
}

func NewACMISource(ctx context.Context, uri string) (Source, error) {

	t, err := uritemplates.Parse(ACMI_OBJECT_TEMPLATE)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse object template, %w", err)
	}

	s := &ACMISource{
		object_template: t,
	}

	return s, nil
}

func (s *ACMISource) Import(ctx context.Context, db database.AccessionNumberDatabase, uris ...string) error {

	for _, u := range uris {

		err := s.importURI(ctx, db, u)

		if err != nil {
			return fmt.Errorf("Failed to import URI '%s', %w", u, err)
		}
	}

	return nil
}

func (s *ACMISource) ObjectURI(ctx context.Context, acc *database.AccessionNumber) (string, error) {

	values := map[string]interface{}{
		"objectid": acc.ObjectId,
	}

	return s.object_template.Expand(values)
}

func (s *ACMISource) importURI(ctx context.Context, db database.AccessionNumberDatabase, u string) error {

	fh, err := os.Open(u)

	if err != nil {
		return fmt.Errorf("Failed to open '%s', %w", u, err)
	}

	defer fh.Close()

	reader := bufio.NewReader(fh)
	lineno := 0

	for {

		select {
		case <-ctx.Done():
			break
		default:
			// pass
		}

		body, err := reader.ReadBytes('\n')

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("Failed to read bytes, %w", err)
		}

		lineno += 1

		if lineno == 1 {
			continue
		}

		if !re_identifiers.Match(body) {
			// fmt.Printf("NOPE: %d '%s'\n", lineno, string(body))
			continue
		}

		m := re_identifiers.FindStringSubmatch(string(body))

		object_id := m[1]
		accession_number := m[2]

		org_uri := fmt.Sprintf("%s://", ACMI_ORGANIZATION_SCHEME)

		a := &database.AccessionNumber{
			AccessionNumber: accession_number,
			ObjectId:        object_id,
			OrganizationURI: org_uri, // update to use https://github.com/sfomuseum/accession-numbers/blob/main/data/acmi.gov.json
		}

		err = db.AddAccessionNumber(ctx, a)

		if err != nil {
			return fmt.Errorf("Failed to add accession number '%s', %w", accession_number, err)
		}
	}

	return nil
}
