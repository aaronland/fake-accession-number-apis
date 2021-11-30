package source

import (
	"context"
	"github.com/aaronland/fake-accession-number-apis/database"
)

type Source interface {
	Import(context.Context, database.AccessionNumberDatabase, ...string) error
}
