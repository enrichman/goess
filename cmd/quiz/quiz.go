package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

type Question struct {
	ID              string
	Question        string
	Answer          int
	PossibleAnswers []string
}

func main() {
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
