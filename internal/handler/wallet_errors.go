package handler

import "errors"

var (
	ErrInvalidFormatID = errors.New("invalid id format: not uuid")
	ErrPathParameterID = errors.New("path parameters: id not found")
)
