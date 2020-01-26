package web

import (
	"encoding/json"
	"net/http"
)

func Decoder(r *http.Request, val interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(val); err != nil {
		return err
	}

	return nil
}
