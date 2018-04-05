package net

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pierre-emmanuelJ/go-exercises/utils"
)

const NetInterfaceName = 0

type NetInfos struct {
	PrevUp   int64
	PrevDown int64
}

func GetInterfaceLine(inet string) ([]string, error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if line, isOk := getLineInfos(scanner.Text(), inet); isOk {
			return line, nil
		}
	}
	//TODO implem Error
	return nil, fmt.Errorf("Interface not found in file '/proc/net/dev'\n Exemple: -n eth0")
}

func getLineInfos(line, inet string) ([]string, bool) {

	s := utils.StringRemoveAllWhiteSpace(line)

	lineInfo := strings.Split(s, " ")

	if len(lineInfo) < 1 {
		return nil, false
	}

	inet += ":"
	if lineInfo[NetInterfaceName] == inet {
		return lineInfo, true
	}

	return nil, false
}
