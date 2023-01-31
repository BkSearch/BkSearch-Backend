package common

type Question struct {
	ID           int       `db:"id, pk"`
	Content      string    `db:"content"`
	AmountAnswer int       `db:"amount_answer"`
	URL          string    `db:"url"`
	Vote         int       `db:"vote"`
}
