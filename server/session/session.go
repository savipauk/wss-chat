package session

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func CheckIfExists(token string) bool {
	file, err := os.OpenFile("sessions", os.O_RDONLY, 0644)
	defer file.Close()
	if err != nil {
		log.Printf("Couldn't check if token %s exists: %v", token, err)
		return false
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		checkingToken := strings.Split(scanner.Text(), "\t")
		if checkingToken[1] == token {
			log.Println("Found user: username " + checkingToken[0])
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error: %v", err)
		return false
	}

	return false
}

func CreateSession(username string, token string) {
	line := username + "\t" + token + "\n"

	file, err := os.OpenFile("sessions", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		log.Printf("Couldn't create session for user %s (%s): %v", username, token, err)
		return
	}

	_, err = file.WriteString(line)
	if err != nil {
		log.Printf("Couldn't create session for user %s (%s): %v", username, token, err)
		return
	}
}
