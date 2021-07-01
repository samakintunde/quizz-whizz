package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
)

func isValidCsv(file *os.File) bool {
	info, err := file.Stat()

	if isCsv := strings.HasSuffix(info.Name(), ".csv"); err != nil || !isCsv {
		return false
	}

	return true
}

// Load questions
func parseQuestions(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)

	if !isValidCsv(file) {
		err := errors.New("File is not a valid CSV File")
		return nil, err
	}

	// Handle error opening file
	if err != nil {
		return nil, err
	}

	// Read file as csv
	// We should probably check if it is a valid csv file first
	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()

	// Handle error reading the csv
	// Probably empty file or imporoperly formatted values
	if err != nil {
		return nil, err
	}

	return records, nil
}

// Ask questions until end
func askQuestion(index int, question []string, reader *bufio.Reader) (string, error) {
	fmt.Printf("%d%-2s", index+1, ".")
	fmt.Printf("%v: ", question[0])

	answer, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	// Validate string
	answer = strings.TrimSpace(answer)

	if length := len(answer); length == 0 {
		fmt.Println("Invalid answer. Try again!")
		askQuestion(index, question, reader)
	}

	return answer, nil
}

func main() {
	// Process CSV into questions
	questions, err := parseQuestions("problems.csv")

	if err != nil {
		fmt.Println(err)
	}

	var count int

	reader := bufio.NewReader(os.Stdin)

	for index, question := range questions {
		answer, err := askQuestion(index, question, reader)

		if err != nil {
			panic("Something really terrible happened!")
		}

		if question[1] == answer {
			count += 1
		}
	}

	fmt.Printf("You got %d/%d.", count, len(questions))
}
