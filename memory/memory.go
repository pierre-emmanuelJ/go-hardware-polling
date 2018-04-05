package memory

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetMemoryInfoByKey(key string) (float64, error) {

	memline, err := getMemLine(key)
	if err != nil {
		return 0, err
	}

	memLineSplited := strings.Split(memline, " ")

	for _, v := range memLineSplited {
		value, err := strconv.Atoi(v)
		if err != nil {
			continue
		} else {
			return float64(value), nil
		}
	}
	//TODO implem Error
	return 0, fmt.Errorf("Value not found for key: %s", key)
}

func getMemLine(key string) (string, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), key) {
			return scanner.Text(), nil
		}
	}
	//TODO implem Error
	return "", fmt.Errorf("key '%s' not found in memory file or file not valid", key)
}
