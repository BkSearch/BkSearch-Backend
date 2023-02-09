package db

type QuestionAPIView struct {
	ID           int    `json:"id"            db:"id"`
	Content      string `json:"content"       db:"content"`
	AmountAnswer int    `json:"amount_answer" db:"amount_answer"`
	URL          string `json:"url"           db:"url"`
	Vote         int    `json:"vote"          db:"vote"`
}

type AnswerAPIView struct {
	ID          int    `json:"id"           db:"id, pk"`
	Content     string `json:"content"      db:"content"`
	Vote        int    `json:"vote"         db:"vote"`
	Question_ID int    `json:"question_id"  db:"question_id"`
	Accepted    bool   `json:"accepted"  db:"accepted"`
}

type DocumentAPIView struct {
	ID                  int     `json:"id"`
	Content             string  `json:"content"`
	Type                int     `json:"type"`
	Vote                int     `json:"vote"`
	Score               float64 `json:"score"`
	QuestionID          int     `json:"question_id"`
	QuesionContent      string  `json:"quesion_content"`
	QuestionURL         string  `json:"quesion_url"`
	QuestionVote        int     `json:"quesion_vote"`
	QuestionCountAwnser int     `json:"question_count_answer"`
}
