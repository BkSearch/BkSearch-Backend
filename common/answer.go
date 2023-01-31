package common

type Answer struct {
	ID          int       `db:"id, pk"`
	Content     string    `db:"content"`
	Vote        int       `db:"vote"`
	Question_ID int       `db:"question_id"`
}


