package main

import (
	"fmt"
	"os"

	"github.com/enrichman/goess/pkg/csv"
	"github.com/enrichman/goess/pkg/goess"
	"github.com/manifoldco/promptui"
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

	var totalCorrect int
	for _, q := range quiz.Questions {
		answerGroup := q.AnswerGroup[0]

		iconizer := func(correct bool) string {
			if correct {
				return promptui.IconGood
			}
			return promptui.IconBad
		}

		promptui.FuncMap["iconizer"] = iconizer

		prompt := promptui.Select{
			HideHelp: true,
			Templates: &promptui.SelectTemplates{
				Active:   fmt.Sprintf("%s {{ .Text | underline }}", promptui.IconSelect),
				Inactive: `  {{ .Text }}`,
				Selected: fmt.Sprintf(`{{ .Correct | iconizer }} (%s) {{ .Text | faint }}`, answerGroup.ID),
				//Selected: `{{ .Correct | iconizer }} (` + answerGroup.ID + `) {{ .Text | faint }}`,
			},
			Label: fmt.Sprintf("%s) %s", answerGroup.ID, q.Text),
			Items: answerGroup.Answers,
		}

		i, _, err := prompt.Run()
		if err != nil {
			slog.Error("prompt", "error", err, "i", i)
			os.Exit(1)
		}

		if answerGroup.Answers[i].Correct {
			totalCorrect += 1
		} else {
			for _, a := range answerGroup.Answers {
				if a.Correct {
					fmt.Println("  " + a.Text)
				}
			}
		}
	}

	fmt.Printf("\nTotal: %d/%d\n", totalCorrect, len(quiz.Questions))

}
