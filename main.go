package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// isValidCsv checks to see if the file is a valid csv file
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

// parseQuestions extracts the questions from the csv file
func parseQuestions(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	if err != nil {
		return nil, err
	}

	if !isValidCsv(file) {
		panic(fmt.Sprintf("Could not open file: %s\n", filePath))
	}

	r := csv.NewReader(file)
	recs, err := r.ReadAll()

	// Handle error reading the csv
	// Probably empty file or improperly formatted values
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// askQuestion asks questions until end
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
		return askQuestion(i, question, r)
	}

	return ans, nil
}

func main() {
	// os.Args
	filePath := flag.String("file", "problems.csv", "Path to the location of the quiz file. (should be a .csv of two columns, first is question and second is answer.)")
	duration := flag.Int("time", 60, "Time to complete all questions")

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

	_ = reader // Take out later

	timer := time.NewTimer(time.Duration(*duration) * time.Second)
	done := make(chan bool) // Uses a semaphore

	go func() {
		for index, question := range questions {
			answer, err := askQuestion(index, question, reader)

			if err != nil {
				panic("Something really terrible happened!")
			}

			if question[1] == answer {
				count += 1
			}
		}
		done <- true
	}()

	select {
	case <-timer.C:
		fmt.Printf("You got %d/%d.\n", count, len(questions))
		close(done)
		return
	case <-done:
		fmt.Printf("You got %d/%d.\n", count, len(questions))
		close(done)
		return
	}

}
