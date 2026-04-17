// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package main

// ReturnCoder attaches a command-line exit code to errors.
type ReturnCoder interface {
	error
	ReturnCode() int
}

// A NoResultsError is returned when a command returns no results.
type NoResultsError string

func (e NoResultsError) Error() string {
	return string(e)
}

// ReturnCode implements the ReturnCoder interface.
func (e NoResultsError) ReturnCode() int {
	return 2
}
