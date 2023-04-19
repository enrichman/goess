package main

import (
	"fmt"
	"os"

	"github.com/enrichman/goess/pkg/csv"
	"github.com/enrichman/goess/pkg/goess"
	"golang.org/x/exp/slog"
)

type Question struct {
	ID              string
	Question        string
	Answer          int
	PossibleAnswers []string
}

func main() {
	questions, err := csv.LoadFile("quiz.csv")
	if err != nil {
		slog.Error("loading CSV", "error", err)
		os.Exit(1)
	}

	quiz := goess.NewQuiz(questions)

	for _, q := range quiz.Questions {
		answerGroup := q.AnswerGroup[0]

		fmt.Printf("%s) %s\n\n", answerGroup.ID, q.Text)
		for _, a := range answerGroup.Answers {
			correct := " "
			if a.Correct {
				correct = "X"
			}
			fmt.Printf("[%s] %s\n", correct, a.Text)
		}
		fmt.Println()
	}
}
