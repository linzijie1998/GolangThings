package mock

import (
	"bufio"
	"os"
	"strings"
)

func ReadFirstLine() string {
	open, err := os.Open("test1.txt")
	defer open.Close()
	if err != nil {
		return ""
	}
	scanner := bufio.NewScanner(open)
	for scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

func ProcessFirstLine() string {
	line := ReadFirstLine()
	dstLine := strings.ReplaceAll(line, "11", "00")
	return dstLine
}
