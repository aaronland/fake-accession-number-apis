package api

import (
	"encoding/json"
	"github.com/aaronland/fake-accession-number-apis/database"
	"github.com/aaronland/go-http-sanitize"
	"net/http"
)

func NewLookupHandler(db database.AccessionNumberDatabase) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		num, err := sanitize.GetString(req, "accession-number")

		if err != nil {
			http.Error(rsp, "Invalid accession number", http.StatusBadRequest)
			return
		}

		source, err := sanitize.GetString(req, "source-uri")

		if err != nil {
			http.Error(rsp, "Invalid source URI", http.StatusBadRequest)
			return
		}

		a, err := db.GetByAccessionNumber(ctx, source, num)

		if err != nil {
			http.Error(rsp, "Failed to retrieve accession number", http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-Type", "application/json")

		enc := json.NewEncoder(rsp)
		err = enc.Encode(a)

		if err != nil {
			http.Error(rsp, "Failed to encode response", http.StatusInternalServerError)
			return
		}

		return
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
