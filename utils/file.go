package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// get path file list
func GetFileList(path string) ([]string, error) {
	var rs, files []string
	PthSep := string(os.PathSeparator)
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		if f.IsDir() {
			files = append(files, path+PthSep)
		} else {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return rs, err
	}

	if len(files) > 0 {
		for _, v := range files {
			if len(strings.Split(v, PthSep)) == 3 && !strings.HasSuffix(v, PthSep) {
				rs = append(rs, v)
			}
		}
	}
	return rs, err
}
