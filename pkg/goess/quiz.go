package goess

import (
	"math/rand"
	"time"

	"github.com/enrichman/goess/pkg/types"
)

type Quiz struct {
	Questions []*types.Question
}

func NewQuiz(questions []*types.Question) *Quiz {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	random.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	questions = questions[:20]

	for _, q := range questions {
		random.Shuffle(len(q.AnswerGroup), func(i, j int) {
			q.AnswerGroup[i], q.AnswerGroup[j] = q.AnswerGroup[j], q.AnswerGroup[i]
		})
	}

	return &Quiz{
		Questions: questions,
	}
}
