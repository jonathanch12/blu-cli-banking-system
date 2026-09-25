package services

import "path/filepath"

func filePathFor(kind, name string) string {
	return filepath.Join("filebase", kind, name+".txt")
}
