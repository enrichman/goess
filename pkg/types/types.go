package types

type Question struct {
	ID          int
	Text        string
	AnswerGroup []*AnswerGroup
}

type AnswerGroup struct {
	ID      string
	Answers []*Answer
}

type Answer struct {
	Correct bool
	Text    string
}
