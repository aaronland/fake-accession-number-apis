package api

import (
	"github.com/aaronland/fake-accession-number-apis/database"
	"github.com/aaronland/fake-accession-number-apis/source"	
	"github.com/aaronland/go-http-sanitize"
	"net/http"
)

func NewRedirectHandler(db database.AccessionNumberDatabase) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		num, err := sanitize.GetString(req, "accession-number")

		if err != nil {
			http.Error(rsp, "Invalid accession number", http.StatusBadRequest)
			return
		}

		source_uri, err := sanitize.GetString(req, "source-uri")

		if err != nil {
			http.Error(rsp, "Invalid source URI", http.StatusBadRequest)
			return
		}

		s, err := source.NewSource(ctx, source_uri)

		if err != nil {
			http.Error(rsp, "Invalid source", http.StatusBadRequest)
			return
		}
		
		a, err := db.GetByAccessionNumber(ctx, source_uri, num)

		if err != nil {
			http.Error(rsp, "Failed to retrieve accession number", http.StatusInternalServerError)
			return
		}
		
		uri, err := s.ObjectURI(ctx, a)

		if err != nil {
			http.Error(rsp, "Failed to retrieve accession number", http.StatusInternalServerError)
			return
		}
		
		http.Redirect(rsp, req, uri, http.StatusFound)
		return
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
