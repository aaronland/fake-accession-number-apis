package database

import (
	"context"
	"fmt"
	"github.com/aaronland/go-roster"
	"net/url"
	"sort"
	"strings"
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

type DatabaseInitializeFunc func(ctx context.Context, uri string) (AccessionNumberDatabase, error)

var databases_roster roster.Roster

func ensureRoster() error {

	if databases_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		databases_roster = r
	}

	return nil
}

func RegisterDatabase(ctx context.Context, scheme string, f DatabaseInitializeFunc) error {

	err := ensureRoster()

	if err != nil {
		return err
	}

	return databases_roster.Register(ctx, scheme, f)
}

func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range databases_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}

func NewDatabase(ctx context.Context, uri string) (AccessionNumberDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := databases_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(DatabaseInitializeFunc)
	return f(ctx, uri)
}
