package node

import (
	"fmt"
	"strconv"

	"github.com/BkSearch/BkSearch-Backend/common"
	"github.com/BkSearch/BkSearch-Backend/db"
	"github.com/BkSearch/BkSearch-Backend/elasticsearch"
	handleerror "github.com/BkSearch/BkSearch-Backend/handle-error"
)

type Node struct {
	es     *elasticsearch.StackOverflow
	itemDB *db.ItemDB
}

func NewNode(es *elasticsearch.StackOverflow, itemDB *db.ItemDB) *Node {
	return &Node{
		es:     es,
		itemDB: itemDB,
	}
}

func (node *Node) SynchronizeData() {
	questions, err := node.itemDB.GetListQuestion(10000, 1)
	if err != nil {
		handleerror.WrapErrorf(err, handleerror.ErrorCodeUnknown, "json.SynchronizeData")
	}

	for _, question := range questions {
    fmt.Println(question)
		var tmpDocument common.Document
		tmpDocument.ID = strconv.Itoa(question.ID)
		tmpDocument.Content = question.Content
		tmpDocument.Type = common.QuestionType

		node.es.Index(tmpDocument)
	}
}

func (node *Node) SynchronizeAnswerData() {
	answers, err := node.itemDB.GetListAnswer(40000, 5)
	if err != nil {
		handleerror.WrapErrorf(err, handleerror.ErrorCodeUnknown, "json.SynchronizeData")
	}

	for i, awnswer := range answers {
    fmt.Println(i)
    fmt.Println(awnswer)
		var tmpDocument common.Document
		tmpDocument.ID = strconv.Itoa(awnswer.ID)
		tmpDocument.Content = awnswer.Content
		tmpDocument.Type = common.AnswerType

		node.es.Index(tmpDocument)
	}
}

func (node *Node) Search(content *string) ([]common.Document, error){
  data, err := node.es.Search(content)
  return data, err
}
