package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
	"reflect"
)

type JsonSender struct {
	writer  http.ResponseWriter
	request *http.Request
	code    int
}

type jsonAppError struct {
	Code  int    `json:"code"`
	Type  string `json:"type"`
	Error any    `json:"error"`
}

func NewJsonSender(w http.ResponseWriter, r *http.Request) *JsonSender {
	return &JsonSender{writer: w, request: r}
}
func sendError(j *JsonSender, e any) {
	w := j.writer
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	err := enc.Encode(jsonAppError{
		Code:  j.code,
		Error: e,
		Type:  reflect.TypeOf(e).Name(),
	})

	if err != nil {
		errBuf := &bytes.Buffer{}
		errEnc := json.NewEncoder(errBuf)
		errEnc.SetEscapeHTML(true)
		_ = errEnc.Encode(jsonAppError{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
			Type:  "encodingError",
		})

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(errBuf.Bytes())
		return
	}

	w.WriteHeader(j.code)
	_, _ = w.Write(buf.Bytes())
}

func send(j *JsonSender, data render.Renderer) {
	render.Status(j.request, j.code)
	if err := render.Render(j.writer, j.request, data); err != nil {
		j.InternalServerError(err)
	}
}

type badRequestError string

type validationError struct {
	Namespace string `json:"namespace"` // can differ when a custom TagNameFunc is registered or
	Field     string `json:"field"`     // by passing alt name to ReportError like below
	Tag       string `json:"tag"`
	Kind      string `json:"kind"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	Param     string `json:"param"`
	Message   string `json:"message"`
}

type validationBadRequestError []validationError

func (j *JsonSender) BadRequest(err error) {
	j.code = http.StatusBadRequest
	var validateErrs validator.ValidationErrors
	if errors.As(err, &validateErrs) {
		allErrs := make(validationBadRequestError, len(validateErrs))
		for i, err := range validateErrs {
			allErrs[i] = validationError{
				Namespace: err.Namespace(),
				Field:     err.Field(),
				Tag:       err.Tag(),
				Kind:      fmt.Sprintf("%v", err.Kind()),
				Type:      fmt.Sprintf("%v", err.Type()),
				Value:     fmt.Sprintf("%v", err.Value()),
				Param:     err.Param(),
				Message:   err.Error(),
			}
		}
		sendError(j, allErrs)
		return
	}

	sendError(j, badRequestError(err.Error()))
}

type internalServerError string

func (j *JsonSender) InternalServerError(err error) {
	j.code = http.StatusInternalServerError
	sendError(j, internalServerError(err.Error()))
}

type conflictError string

func (j *JsonSender) Conflict(err error) {
	j.code = http.StatusConflict
	sendError(j, conflictError(err.Error()))
}

type unauthorizedError string

func (j *JsonSender) Unauthorized(err error) {
	j.code = http.StatusUnauthorized
	sendError(j, unauthorizedError(err.Error()))
}

type notFoundError string

func (j *JsonSender) NotFound(err error) {
	j.code = http.StatusNotFound
	sendError(j, notFoundError(err.Error()))
}

type unprocessableEntityError string

func (j *JsonSender) UnprocessableEntity(err error) {
	j.code = http.StatusUnprocessableEntity
	sendError(j, unprocessableEntityError(err.Error()))
}

func (j *JsonSender) Created(res render.Renderer) {
	j.code = http.StatusCreated
	send(j, res)
}

func (j *JsonSender) Ok(res render.Renderer) {
	j.code = http.StatusOK
	send(j, res)
}
