package notifications

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func NotifyCheckout() (int, string, error) {
	file, err := os.Open("./notifications.txt")
	if err != nil {
		return 0, "", fmt.Errorf("ошибка открытия файла: %s", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 4 {
			log.Print("Некорректный формат строки: ", line)
			continue
		}
		current_time := getCurrentTime()
		reach, err := strconv.Atoi(parts[2])
		if err != nil {
			log.Print("Ошибка преобразования : ", err)
			continue
		}
		interval, err := strconv.Atoi(parts[3])
		if err != nil {
			log.Print("Ошибка преобразования : ", err)
			continue
		}
		if current_time >= reach {
			file.Close()
			newline := fmt.Sprintf("%s:%s:%d:%s", parts[0], parts[1], reach+interval, parts[3])
			replaceStringInFile(line, newline)
			ID, City := notificationTime(parts[0], parts[1])
			return ID, City, nil
		}
	}
	return 0, "", nil
}

func getCurrentTime() int {
	currentTime := time.Now()
	epoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	duration := currentTime.Sub(epoch)
	return int(duration.Minutes())
}

func notificationTime(ID string, City string) (int, string) {
	ID_value, err := strconv.Atoi(ID)
	if err != nil {
		log.Print("Ошибка преобразования : ", err)
		return 0, ""
	}
	return ID_value, City
}

func replaceStringInFile(oldStr string, newStr string) error {
	data, err := os.ReadFile("./notifications.txt")
	if err != nil {
		return err
	}
	newContent := strings.Replace(string(data), oldStr, newStr, -1)
	err = os.WriteFile("./notifications.txt", []byte(newContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

func SetNotification(ID int, City string, interval int) error {
	file, err := os.Open("./notifications.txt")
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %s", err)
	}
	defer file.Close()
	current_time := getCurrentTime()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 4 {
			log.Print("Некорректный формат строки: ", line)
			continue
		}
		if strconv.Itoa(ID) == parts[0] {
			file.Close()
			newline := fmt.Sprintf("%d:%s:%d:%d", ID, City, current_time+interval, interval)
			replaceStringInFile(line, newline)
			return nil
		}
	}
	file.Close()
	file, err = os.OpenFile("./notifications.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.WriteString(fmt.Sprintf("%d:%s:%d:%d", ID, City, current_time+interval, interval) + "\n"); err != nil {
		return err
	}
	return nil
}
