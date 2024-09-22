package common

import (
	"bufio"
	"os"
)

// 读取文件的每一行，并返回一个字符串数组
func ReadLinesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		// 返回错误，如果文件无法打开
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		// 如果在扫描过程中出现错误，返回错误
		return nil, err
	}

	return lines, nil
}
