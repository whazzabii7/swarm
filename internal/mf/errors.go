package mf

import "errors"

var ErrExecFailed = errors.New("failed to execute")
var ErrWrongArguments = errors.New("wrong arguments")
var ErrCorruptedData = errors.New("coruppted data from database")
