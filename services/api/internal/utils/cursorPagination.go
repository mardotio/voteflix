package utils

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/uptrace/bun"
	"net/http"
	"net/url"
	"strconv"
)

type CursorDataItem interface {
	CursorId() *string
}

type CursorPaginatedResponse[T CursorDataItem] struct {
	Data     []T     `json:"data"`
	Limit    int     `json:"limit"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
}

func (body CursorPaginatedResponse[T]) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type Cursor[T CursorDataItem] struct {
	limit     int
	direction string
	marker    *string

	hasNext         bool
	hasPrevious     bool
	hasMoreData     bool
	isStart         bool
	cursorDirection string
	data            []T
}

type cursorParamsValidator struct {
	Before    string `validate:"omitempty,uuid4"`
	After     string `validate:"omitempty,uuid4,excluded_unless=Before ''"`
	Limit     string `validate:"omitempty,number"`
	Direction string `validate:"omitempty,oneof=asc desc"`
}

func (c *Cursor[T]) IsStart() bool   { return c.isStart }
func (c *Cursor[T]) IsAfter() bool   { return c.marker == nil || c.cursorDirection == "after" }
func (c *Cursor[T]) IsBefore() bool  { return !c.IsAfter() }
func (c *Cursor[T]) Marker() *string { return c.marker }

func (c *Cursor[T]) Order() bun.Safe {
	if c.IsBefore() {
		if c.direction == "asc" {
			return "desc"
		}

		return "asc"
	}

	return bun.Safe(c.direction)
}

func (c *Cursor[T]) Comparator() bun.Safe {
	if c.IsBefore() {
		if c.direction == "asc" {
			return "<"
		}
		return ">"
	}

	if c.direction == "asc" {
		return ">"
	}

	return "<"
}

func (c *Cursor[T]) FetchLimit() int {
	return c.limit + 1
}

func (c *Cursor[T]) WithData(data []T) *Cursor[T] {
	c.withTrimmedData(data)
	dataLen := len(c.data)
	c.hasMoreData = len(data) > dataLen

	if dataLen <= 0 {
		c.hasNext = false
		c.hasPrevious = false
	} else if c.IsBefore() {
		c.hasNext = true
		c.hasPrevious = c.hasMoreData
	} else {
		c.hasNext = c.hasMoreData
		c.hasPrevious = !c.isStart
	}

	return c
}

func (c *Cursor[T]) withTrimmedData(data []T) []T {
	c.data = data

	if len(c.data) > c.limit {
		c.data = c.data[:c.limit]
	}

	dataLen := len(c.data)

	// Reverse data to keep correct order
	if c.IsBefore() {
		for i := dataLen/2 - 1; i >= 0; i-- {
			opp := dataLen - 1 - i
			data[i], data[opp] = data[opp], data[i]
		}
	}

	return data
}

func (c *Cursor[T]) ToResponse() CursorPaginatedResponse[T] {
	response := CursorPaginatedResponse[T]{
		Limit: c.limit,
		Data:  c.data,
	}

	if c.hasNext {
		response.Next = c.data[len(c.data)-1].CursorId()
	}

	if c.hasPrevious {
		response.Previous = c.data[0].CursorId()
	}

	return response
}

func NewCursorFromMap[T CursorDataItem](params url.Values, validate *validator.Validate) (Cursor[T], error) {
	dirtyParams := cursorParamsValidator{
		Before:    params.Get("before"),
		After:     params.Get("after"),
		Limit:     params.Get("limit"),
		Direction: params.Get("direction"),
	}
	cursor := Cursor[T]{limit: 10, isStart: true, direction: "asc"}

	if err := validate.Struct(&dirtyParams); err != nil {
		return cursor, err
	}

	if dirtyParams.Limit != "" {
		limit, limitErr := strconv.Atoi(dirtyParams.Limit)

		if limitErr != nil {
			return cursor, limitErr
		} else if limit < 1 || limit > 100 {
			return cursor, errors.New("limit must be between 1 and 100")
		}

		cursor.limit = limit
	}

	if dirtyParams.Before != "" {
		cursor.isStart = false
		cursor.marker = &dirtyParams.Before
		cursor.cursorDirection = "before"
	}

	if dirtyParams.After != "" {
		cursor.isStart = false
		cursor.marker = &dirtyParams.After
		cursor.cursorDirection = "after"
	}

	if dirtyParams.Direction != "" {
		cursor.direction = dirtyParams.Direction
	}

	return cursor, nil
}
