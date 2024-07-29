// Copyright (C) Damien Dart, <damiendart@pobox.com>.
// This file is distributed under the MIT licence. For more information,
// please refer to the accompanying "LICENCE" file.

package editor

// Editor is implemented by anything that has a OpenFile method, which
// is used to open a file in a text editor.
type Editor interface {
	OpenFile(path string, cwd string) error
}
