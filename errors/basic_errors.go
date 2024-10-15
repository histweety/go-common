package errors

import "errors"

var ErrDataNotFound = errors.New("ERR_DATA_NOT_FOUND")
var ErrDataInsert = errors.New("ERR_DATA_INSERT")
var ErrDataUpdate = errors.New("ERR_DATA_UPDATE")
var ErrDataDelete = errors.New("ERR_DATA_DELETE")

var ErrStructValidation = errors.New("ERR_STRUCT_VALIDATION")
var ErrStructParsing = errors.New("ERR_STRUCT_PARSING")
