package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getText(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func checkNotificationComandFormat(input string) (string, int, error) {
	parts := strings.Fields(input)
	if len(parts) != 3 {
		return "", 0, fmt.Errorf("неверный формат строки")
	}
	message := parts[1]
	duration, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", 0, fmt.Errorf("неверный формат числа")
	}
	return message, duration, nil
}

func swap(x int) int {
	if x == 1 {
		return 2
	} else {
		return 1
	}
}

func sliceStringWithFirstSpace(input string) string {
	index := strings.Index(input, " ")
	if index == -1 {
		return input
	}
	return input[:index]
}
