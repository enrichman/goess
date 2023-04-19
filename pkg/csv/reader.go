package csv

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func LoadFile() {
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

	b, err := os.ReadFile("out.json")
	if err != nil {
		log.Fatal(err)
	}

	var questions []Question
	err = json.Unmarshal(b, &questions)
	if err != nil {
		log.Fatal(err)
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	random.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	quiz := questions[:20]

	for _, q := range quiz {
		fmt.Printf("%s) %s [%d]\n\n", q.ID, q.Question, q.Answer)
		fmt.Println(" - ", q.PossibleAnswers[0])
		fmt.Println(" - ", q.PossibleAnswers[1])
		fmt.Println(" - ", q.PossibleAnswers[2])
		fmt.Println()
	}
}
