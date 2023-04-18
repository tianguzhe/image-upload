package dir

import (
	"fmt"
	"os"
	"strings"
)

func GetDirFile(path string) []string {

	result := make([]string, 0)

	dir, _ := os.ReadDir(path)

	for _, fi := range dir {
		if !fi.IsDir() && (strings.HasSuffix(fi.Name(), ".png") || strings.HasPrefix(fi.Name(), ".jpg")) {
			result = append(result, fmt.Sprintf("%s/%s", path, fi.Name()))
		} else {
			continue
		}
	}

	return result
}
