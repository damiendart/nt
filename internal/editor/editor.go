// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package editor

import "errors"

// ErrPathNotAbsolute is returned if the path provided is not absolute.
var ErrPathNotAbsolute = errors.New("path is not absolute")

// Opener is the interface that wraps the basic Open method.
type Opener interface {
	// Open opens the file specified in path in a text editor. The path
	// must be an absolute path, otherwise Open must return
	// ErrPathNotAbsolute. If dir is provided, it may be used to change
	// the working directory before launching the editor.
	Open(path string, dir string) error
}
