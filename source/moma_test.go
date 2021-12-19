package source

import (
	"context"
	"testing"
	"github.com/aaronland/fake-accession-number-apis/database"	
)

func TestMoMAImport(t *testing.T){
	t.Skip()
}

func TestMoMAObjectId(t *testing.T) {

	ctx := context.Background()
	
	s, err := NewSource(ctx, "moma://")

	if err != nil {
		t.Fatalf("Failed to create new moma:// source, %v", err)
	}

	tests := map[string]string {
		"5": "https://www.moma.org/collection/works/5",
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
