package partition

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetPartitionMountPointPath(partition string) (string, error) {

	file, err := os.Open("/proc/mounts")
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if s, ret := lineGetPath(scanner.Text(), partition); ret {
			return s, nil
		}
	}

	//TODO implem Error
	return "", fmt.Errorf("Partition not found.\n Exemple: -p \"/dev/sda1\"\nFor more details see '/proc/mounts' to see mounted partitions")
}

func lineGetPath(s, partition string) (string, bool) {

	line := strings.Split(s, " ")

	if len(line) < 2 {
		return "", false
	}

	if line[0] == partition {
		return line[1], true
	}

	return "", false
}
