package cli

import "fmt"

// An UnknownOptionError is returned if an option is found that is not
// in the option specification.
type UnknownOptionError string

// A MissingOptionValueError is returned if an option requiring a value
// is found without a value.
type MissingOptionValueError string

// An UnexpectedOptionValueError is returned if an option with a value
// is found but does not require a value.
type UnexpectedOptionValueError string

func (e UnknownOptionError) Error() string {
	return fmt.Sprintf("unknown option %q", string(e))
}

func (e MissingOptionValueError) Error() string {
	return fmt.Sprintf("missing value for %q", string(e))
}

func (e UnexpectedOptionValueError) Error() string {
	return fmt.Sprintf("unexpected value for %q", string(e))
}
