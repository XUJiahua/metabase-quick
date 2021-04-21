package util

import (
	"path/filepath"
	"strings"
)

func GetFilenameWithExt(filename string) string {
	filename = filepath.Base(filename)
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}
