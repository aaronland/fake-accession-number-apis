package source

import (
	"context"
	"testing"
	"github.com/aaronland/fake-accession-number-apis/database"	
)

func TestMetMuseumImport(t *testing.T){
	t.Skip()
}

func TestMetMuseumObjectId(t *testing.T) {

	ctx := context.Background()
	
	s, err := NewSource(ctx, "metmuseum://")

	if err != nil {
		t.Fatalf("Failed to create new metmuseum:// source, %v", err)
	}

	tests := map[string]string {
		"5": "http://www.metmuseum.org/art/collection/search/5",
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
