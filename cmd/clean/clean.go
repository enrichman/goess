package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

func main() {
	f, err := os.Open("listatoquizkb2019.csv")
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	questions := []*Question{}
	var currentQuestion *Question

	for i, record := range records {
		mod := i % 4

		if mod == 0 {
			currentQuestion = &Question{
				ID:              record[0],
				Question:        capitalize(strings.ToLower(record[1])),
				PossibleAnswers: []string{},
			}
			questions = append(questions, currentQuestion)
			continue
		}

		possibleAnswer := record[1]
		possibleAnswer = strings.TrimPrefix(possibleAnswer, fmt.Sprintf("%d.", mod))
		possibleAnswer = strings.TrimSpace(possibleAnswer)

		currentQuestion.PossibleAnswers = append(currentQuestion.PossibleAnswers, possibleAnswer)

		if strings.TrimSpace(strings.ToLower(record[2])) == "x" {
			currentQuestion.Answer = mod
		}
	}

	// fmt.Printf("%s) %s\n%s\nCorrect: %d\n",
	// 	questions[0].ID,
	// 	questions[0].Question,
	// 	strings.Join(questions[0].PossibleAnswers, "\n"),
	// 	questions[0].Answer,
	// )

	b, _ := json.MarshalIndent(questions, "", "\t")
	fmt.Println(string(b))
}

type Question struct {
	ID              string
	Question        string
	Answer          int
	PossibleAnswers []string
}

func capitalize(str string) string {
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
