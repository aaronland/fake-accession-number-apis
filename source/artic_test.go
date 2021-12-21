package source

import (
	"context"
	"github.com/aaronland/fake-accession-number-apis/database"
	"testing"
)

func TestArtICImport(t *testing.T) {
	t.Skip()
}

func TestArtICObjectId(t *testing.T) {

	ctx := context.Background()

	s, err := NewSource(ctx, "artic://")

	if err != nil {
		t.Fatalf("Failed to create new artic:// source, %v", err)
	}

	tests := map[string]string{
		"99990": "https://www.artic.edu/artworks/99990/",
	}

	for id, expected_uri := range tests {

		acc := &database.AccessionNumber{
			ObjectId: id,
		}

		uri, err := s.ObjectURI(ctx, acc)

		if err != nil {
			t.Fatalf("Failed to derive object URI for %s, %v", id, err)
		}

		if uri != expected_uri {
			t.Fatalf("Invalid object URI for %s. Got '%s' but expected '%s'", id, uri, expected_uri)
		}
	}
}
