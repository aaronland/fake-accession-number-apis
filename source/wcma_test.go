package source

import (
	"context"
	"github.com/aaronland/fake-accession-number-apis/database"
	"testing"
)

func TestWCMAImport(t *testing.T) {
	t.Skip()
}

func TestWCMAObjectId(t *testing.T) {

	ctx := context.Background()

	s, err := NewSource(ctx, "wcma://")

	if err != nil {
		t.Fatalf("Failed to create new wcma:// source, %v", err)
	}

	tests := map[string]string{
		"25936": "http://egallery.williams.edu/objects/25936/",
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
