package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type jsonEncodingErrorResponse struct {
	Message string `json:"message"`
}

func JsonError(w http.ResponseWriter, error interface{}, code int) {
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	err := enc.Encode(error)

	if err != nil {
		errBuf := &bytes.Buffer{}
		errEnc := json.NewEncoder(errBuf)
		errEnc.SetEscapeHTML(true)
		_ = errEnc.Encode(jsonEncodingErrorResponse{Message: err.Error()})

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(errBuf.Bytes())
		return
	}

	w.WriteHeader(code)
	_, _ = w.Write(buf.Bytes())
}
