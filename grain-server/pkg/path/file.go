package filex

import (
	"os"
)

// PathIsNotExist 路径不存在
func PathIsNotExist(path string) bool {
	f, err := os.Stat(path)
	if os.IsNotExist(err) {
		return true
	} else {
		if f.IsDir() {
			return false
		}
		return true
	}
}

// FileIsNotExist 文件不存在
func FileIsNotExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return true
	} else {
		return false
	}
}

func ReadFile(fileName string) string {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return ""
	}
	return string(file)
}
