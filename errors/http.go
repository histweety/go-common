package errors

import "errors"

var ErrNotFound = errors.New("error data not found")
var ErrInsert = errors.New("error inserting data")
var ErrUpdate = errors.New("error updating data")
var ErrDelete = errors.New("error deleting data")
var ErrValidateStruct = errors.New("error validating struct")
var ErrParsingBody = errors.New("error parsing body")
var ErrSomethingWentWrong = errors.New("error something went wrong")
