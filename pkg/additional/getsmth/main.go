package getsmth

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const apiFilePath = "./secret/APItokens.txt"

func GetAPIkey(key string) (string, error) {
	data, err := os.ReadFile(apiFilePath)
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 2 && strings.TrimSpace(parts[0]) == key {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return "", fmt.Errorf("getAPIkey: value not found for key %s", key)
}
