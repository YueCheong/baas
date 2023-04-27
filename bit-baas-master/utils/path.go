package utils

import (
	"flag"
	"os"
	"path/filepath"
	"strconv"
)

func ConfigPath() string {
	path, _ := os.Getwd()
	if flag.Lookup("test.v") != nil {
		path = filepath.Dir(path)
	}

	path = path + string(os.PathSeparator) + "artifacts"

	return path
}

func ConfigPathWithId(id int) string {
	path, _ := os.Getwd()
	if flag.Lookup("test.v") != nil {
		path = filepath.Dir(path)
	}

	path = path + string(os.PathSeparator) + "artifacts" +
		string(os.PathSeparator) + strconv.Itoa(id)

	return path
}
