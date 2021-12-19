package source

import (
	"context"
	"testing"
	"github.com/aaronland/fake-accession-number-apis/database"	
)

func TestNGAImport(t *testing.T){
	t.Skip()
}

func TestNGAObjectId(t *testing.T) {

	ctx := context.Background()
	
	s, err := NewSource(ctx, "nga://")

	if err != nil {
		t.Fatalf("Failed to create new nga:// source, %v", err)
	}

	tests := map[string]string {
		"89682": "https://www.nga.gov/collection/art-object-page.89682.html",
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
