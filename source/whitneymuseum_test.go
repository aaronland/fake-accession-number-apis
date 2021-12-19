package source

import (
	"context"
	"testing"
	"github.com/aaronland/fake-accession-number-apis/database"	
)

func TestWhitneyMuseumImport(t *testing.T){
	t.Skip()
}

func TestWhitneyMuseumObjectId(t *testing.T) {

	ctx := context.Background()
	
	s, err := NewSource(ctx, "whitneymuseum://")

	if err != nil {
		t.Fatalf("Failed to create new whitneymuseum:// source, %v", err)
	}

	tests := map[string]string {
		"55448": "https://whitney.org/collection/works/55448",
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
