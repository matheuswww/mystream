package router

import (
	"encoding/json"
	"io"
)

func BindJson(r io.Reader, s interface{}) error {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(s); err != nil {
		return err
	}
	return nil
}