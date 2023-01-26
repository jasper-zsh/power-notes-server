package util

import "fmt"

func FileKey(projectName, filePath string) string {
	return fmt.Sprintf("%s-%s", projectName, filePath)
}
