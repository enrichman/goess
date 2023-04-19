package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/enrichman/goess/pkg/types"
	"golang.org/x/exp/slices"
)

func LoadFile(filename string) ([]*types.Question, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file '%s': %w", filename, err)
	}

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("reading CSV file '%s': %w", filename, err)
	}

	questionMap := map[string]*types.Question{}

	var currentQuestion *types.Question

	for _, record := range records {
		answerGroupID := strings.TrimSpace(record[0])
		text := strings.TrimSpace(record[1])

		// if not empty this line is a new question
		if answerGroupID != "" {
			// check if the question was already there
			if q, found := questionMap[text]; !found {
				currentQuestion = &types.Question{
					ID:          len(questionMap) + 1,
					Text:        text,
					AnswerGroup: []*types.AnswerGroup{},
				}
			} else {
				currentQuestion = q
			}

			currentQuestion.AnswerGroup = append(currentQuestion.AnswerGroup, &types.AnswerGroup{
				ID:      answerGroupID,
				Answers: []*types.Answer{},
			})

			questionMap[text] = currentQuestion
			continue
		}

		// this is an answer
		correct := strings.ToLower(strings.TrimSpace(record[2]))
		answer := &types.Answer{
			Correct: (correct == "x"),
			Text:    text,
		}

		// get the last answerGroup and append the answer
		currentAnswerGroup := currentQuestion.AnswerGroup[len(currentQuestion.AnswerGroup)-1]
		currentAnswerGroup.Answers = append(currentAnswerGroup.Answers, answer)
	}

	questions := []*types.Question{}
	for _, v := range questionMap {
		questions = append(questions, v)
	}
	slices.SortFunc(questions, func(a, b *types.Question) bool {
		return a.ID < b.ID
	})

	return questions, nil
}
