package common

type DocumentType int

const (
	QuestionType DocumentType = iota
	AnswerType
)

type Document struct {
	ID       string
	Title    string
	Content  string
	Type     DocumentType
	Accepted bool
}

type DocumentV2 struct {
	ID         int          `db:"id, pk"`
	Title      int          `db:"title"`
	Content    string       `db:"content"`
	Type       DocumentType `db:"type"`
	Vote       int          `db:"vote"`
	QuestionID int          `db:"question_id"`
}
