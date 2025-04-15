// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package cli

import (
	"fmt"
	"strings"
)

// OptionMap represents parsed options, the key being the option and
// value being the parsed value. For options that do not take an
// value the value will be "".
type OptionMap map[string]string

// OptionType represents whether an option does not accept a value,
// requires a value, or optionally takes a value.
type OptionType int

// A Spec is the option specification that arguments are parsed against.
// The key is the option and the value is an OptionType.
type Spec map[string]OptionType

const (
	// NoValueAccepted is an OptionType denoting that an option does not
	// accept a value and will result in an error if one is provided.
	NoValueAccepted OptionType = iota

	// ValueOptional is an OptionType denoting that an option takes an
	// optional value.
	ValueOptional

	// ValueRequired is an OptionType denoting that an option requires
	// a value and will result in an error if one is omitted.
	ValueRequired
)

// An UnknownOptionError is returned if an option is found that is not
// in the option specification.
type UnknownOptionError string

// A MissingOptionValueError is returned if an option requiring a value
// is found without a value.
type MissingOptionValueError string

// An UnexpectedOptionValueError is returned if an option with a value
// is found but does not require a value.
type UnexpectedOptionValueError string

// ParseArgs parses command-line options in a similar style to the POSIX
// "getopt" option-parsing functionality, where command-line arguments
// are parsed against an option specification. Both long and short
// options are supported.
func ParseArgs(args []string, spec Spec) (OptionMap, []string, error) {
	var currentOption string
	var remainingArgsStart int
	options := make(OptionMap)

	for _, arg := range args {
		if arg == "--" {
			remainingArgsStart++
			break
		}

		if strings.HasPrefix(arg, "--") {
			opt, val, _ := strings.Cut(arg[2:], "=")
			if val != "" {
				if t, ok := spec[opt]; ok {
					if t == NoValueAccepted {
						return OptionMap{}, []string{}, UnexpectedOptionValueError(opt)
					}
					options[opt] = val
				} else {
					return OptionMap{}, []string{}, UnknownOptionError(opt)
				}
			} else {
				if t, ok := spec[opt]; ok {
					if t == ValueRequired {
						currentOption = opt
					} else {
						options[opt] = ""
					}
				} else {
					return OptionMap{}, []string{}, UnknownOptionError(opt)
				}
			}
		} else if strings.HasPrefix(arg, "-") {
			opts := arg[1:]
			for i, opt := range opts {
				if t, ok := spec[string(opt)]; ok {
					if t == ValueRequired {
						if i == len(opts)-1 {
							currentOption = string(opt)
						} else {
							options[string(opt)] = opts[i+1:]
							break
						}
					} else if t == ValueOptional {
						options[string(opt)] = opts[i+1:]
						break
					} else {
						options[string(opt)] = ""
					}
				} else {
					return OptionMap{}, []string{}, UnknownOptionError(opt)
				}
			}
		} else {
			if currentOption != "" {
				options[currentOption] = arg
				currentOption = ""
			} else {
				break
			}
		}
		remainingArgsStart++
	}

	if currentOption != "" {
		return OptionMap{}, []string{}, MissingOptionValueError(currentOption)
	}

	return options, args[remainingArgsStart:], nil
}

func (e UnknownOptionError) Error() string {
	return fmt.Sprintf("unknown option %q", string(e))
}

func (e MissingOptionValueError) Error() string {
	return fmt.Sprintf("missing value for %q", string(e))
}

func (e UnexpectedOptionValueError) Error() string {
	return fmt.Sprintf("unexpected value for %q", string(e))
}
