package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

func isValidCsv(f *os.File) bool {
	info, err := f.Stat()

	if err != nil {
		return false
	}

	if !strings.HasSuffix(info.Name(), ".csv") {
		return false
	}

	return true
}

// Load questions
func parseQuestions(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)

	defer file.Close()

	if !isValidCsv(file) {
		panic("File is not a valid CSV File")
	}

	if err != nil {
		return nil, err
	}

	r := csv.NewReader(file)
	recs, err := r.ReadAll()

	// Handle error reading the csv
	// Probably empty file or imporoperly formatted values
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Ask questions until end
func askQuestion(i int, question []string, r *bufio.Reader) (string, error) {
	fmt.Printf("%d%-2s", i+1, ".")
	fmt.Printf("%v: ", question[0])

	ans, err := r.ReadString('\n')

	if err != nil {
		return "", err
	}

	// Validate string
	ans = strings.TrimSpace(ans)

	if len(ans) == 0 {
		fmt.Println("Invalid answer. Try again!")
		askQuestion(i, question, r)
	}

	return ans, nil
}

func main() {
	// os.Args
	filePath := flag.String("file", "", "Path to the location of the quiz file. (should be a .csv of two columns, first is question and second is answer.)")

	flag.Parse()

	if *filePath == "" {
		panic("Please, pass a valid filePath. Run './main -h' for more help.")
	}

	// Process CSV into questions
	questions, err := parseQuestions(*filePath)

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
