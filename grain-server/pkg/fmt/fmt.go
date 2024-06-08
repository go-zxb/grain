package fmtx

import (
	"fmt"
	"os/exec"
	"strings"
)

func FmtCode(filePath string) error {
	cmd := exec.Command("gofmt", "-w", filePath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to goutil Go code: %v", err)
	}
	return nil
}

func FormatGoCode(filePath string) error {
	cmd := exec.Command("gofmt", "-w", filePath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to goutil Go code: %v", err)
	}
	return nil
}

// FileExtension 获取文件后缀
func FileExtension(filename string) string {
	index := strings.LastIndex(filename, ".")
	if index == -1 || index == len(filename)-1 {
		return ""
	}
	return strings.TrimSpace(filename[index+1:])
}

// FileNameWithoutExtension 获取文件名
func FileNameWithoutExtension(fileName, suffix string) string {
	fileName = strings.TrimSuffix(fileName, suffix)
	return fileName
}
