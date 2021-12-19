package source

import (
	"context"
	"fmt"
	"github.com/aaronland/fake-accession-number-apis/database"
	"github.com/aaronland/go-roster"
	"net/url"
	"sort"
	"strings"
)

type Source interface {
	Import(context.Context, database.AccessionNumberDatabase, ...string) error
	ObjectURI(context.Context, database.AccessionNumber) (string, error)
}

type SourceInitializeFunc func(ctx context.Context, uri string) (Source, error)

var sources_roster roster.Roster

func ensureRoster() error {

	if sources_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		sources_roster = r
	}

	return nil
}

func RegisterSource(ctx context.Context, scheme string, f SourceInitializeFunc) error {

	err := ensureRoster()

	if err != nil {
		return err
	}

	return sources_roster.Register(ctx, scheme, f)
}

func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range sources_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}

func NewSource(ctx context.Context, uri string) (Source, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := sources_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(SourceInitializeFunc)
	return f(ctx, uri)
}
