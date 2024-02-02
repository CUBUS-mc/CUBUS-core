package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

func readFile(path string) string {
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		return ""
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	fileContent, err := io.ReadAll(file)
	return string(fileContent)
}

func writeFile(path string, content string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.WriteString(content)
	if err != nil {
		log.Fatal(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
}

func input(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	print(prompt)
	userInput, _ := reader.ReadString('\n')
	return strings.TrimSpace(userInput)
}

func askForString(prompt string, defaultValue string, acceptCondition func(string) bool) string {
	userInput := input(prompt)
	if userInput == "" {
		return defaultValue
	}
	if !(acceptCondition(userInput)) {
		return askForString(prompt, defaultValue, acceptCondition)
	}
	return userInput
}

func askForBoolean(prompt string, defaultValue bool) bool {
	userInput := strings.ToLower(askForString(
		prompt,
		"default",
		func(input string) bool {
			switch strings.ToLower(input) {
			case "y", "yes", "n", "no":
				return true
			default:
				return false
			}
		},
	))
	switch strings.ToLower(userInput) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		return defaultValue
	}
}
