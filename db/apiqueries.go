package db

import (
	"fmt"
	"strconv"

	"github.com/BkSearch/BkSearch-Backend/common"
)

func (db *ItemDB) GetAnswerDetailAPI(query common.QueryAnswer) (*AnswerAPIView, error) {
	var answer AnswerAPIView
	row := db.dbWrite.QueryRowx(`SELECT * FROM "Answer" WHERE id = $1 LIMIT 1`, query.AnswerIndex)
	err := row.StructScan(&answer)
	return &answer, err
}

func (db *ItemDB) GetQuestionDetailAPI(query common.QueryQuestion) (*QuestionAPIView, error) {
	var question QuestionAPIView
	row := db.dbWrite.QueryRowx(`SELECT * FROM "Question" WHERE id = $1 LIMIT 1`, query.QuestionIndex)
	err := row.StructScan(&question)
	return &question, err
}

func (db *ItemDB) GetListDocumentDetailAPI(documents []common.Document) ([]DocumentAPIView, error) {
	var documentAPIs []DocumentAPIView
	for _, document := range documents {
		documentIndex, err := strconv.Atoi(document.ID)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		if document.Type == common.QuestionType {
			question, _ := db.getAQuestion(documentIndex)
			tmpDocument := DocumentAPIView{
				ID:                  documentIndex,
				Content:             document.Content,
				Type:                int(document.Type),
				Vote:                question.Vote,
				QuestionID:          question.ID,
				QuestionURL:         question.URL,
				QuestionVote:        question.Vote,
				QuesionContent:      question.Content,
				QuestionCountAwnser: question.AmountAnswer,
			}

			documentAPIs = append(documentAPIs, tmpDocument)
		} else {
			answer, _ := db.getAAnswer(documentIndex)
			question, _ := db.getAQuestion(answer.Question_ID)
			tmpDocument := DocumentAPIView{
				ID:                  documentIndex,
				Content:             document.Content,
				Type:                int(document.Type),
				Vote:                answer.Vote,
				QuestionID:          question.ID,
				QuestionURL:         question.URL,
				QuestionVote:        question.Vote,
				QuesionContent:      question.Content,
				QuestionCountAwnser: question.AmountAnswer,
			}
			documentAPIs = append(documentAPIs, tmpDocument)
		}
	}

	return documentAPIs, nil
}
