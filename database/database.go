package database

import (
	"context"
)

type AccessionNumber struct {
	AccessionNumber string `json:"accession_number"`
	ObjectId        string `json:"object_id"`
	OrganizationURI string `json:"organization_uri"`
}

type AccessionNumberDatabase interface {
	GetByAccessionNumber(context.Context, string, string) (*AccessionNumber, error)
	GetByObjectId(context.Context, string, string) (*AccessionNumber, error)
	AddAccessionNumber(context.Context, *AccessionNumber) error
	RemoveAccessionNumber(context.Context, *AccessionNumber) error
}
