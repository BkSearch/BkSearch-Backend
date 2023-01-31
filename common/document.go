package common

type DocumentType int

const (
  QuestionType DocumentType = iota
  AnswerType 
)

type Document struct {
	ID      string
	Content string
	Type    DocumentType
}

