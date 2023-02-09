package node

import (
	"fmt"
	"strconv"
	"strings"

	"sync"

	"github.com/BkSearch/BkSearch-Backend/common"
	"github.com/BkSearch/BkSearch-Backend/db"
	"github.com/BkSearch/BkSearch-Backend/elasticsearch"
	handleerror "github.com/BkSearch/BkSearch-Backend/handle-error"
)

type Node struct {
	es        *elasticsearch.StackOverflow
	itemDB    *db.ItemDB
	wg        *sync.WaitGroup
	Semaphore chan struct{}
}

func NewNode(es *elasticsearch.StackOverflow, itemDB *db.ItemDB, amountRountine int) *Node {
	return &Node{
		es:     es,
		itemDB: itemDB,
    Semaphore: make(chan struct{}, amountRountine),
    wg: &sync.WaitGroup{},
	}
}

func (node *Node) SynchronizeData() {
	questions, err := node.itemDB.GetListQuestion(10000, 1)
	if err != nil {
		handleerror.WrapErrorf(err, handleerror.ErrorCodeUnknown, "json.SynchronizeData")
	}
  node.wg.Add(len(questions))
	for _, question := range questions {
		fmt.Println(question)
		var tmpDocument common.Document
		tmpDocument.ID = strconv.Itoa(question.ID)
		tmpDocument.Title = question.Content
		tmpDocument.Type = common.QuestionType

    go node.SynchronizeDocumentData(tmpDocument)
	}
  node.wg.Wait()
}
func (node *Node) SynchronizeDocumentData(document common.Document) {
  defer func() {
    <- node.Semaphore
    node.wg.Done()
  }()

  node.Semaphore <- struct{}{}
  node.es.Index(document)
}

func (node *Node) SynchronizeAnswerData() {
	answers, err := node.itemDB.GetListAnswer(100000, 2)
	if err != nil {
		handleerror.WrapErrorf(err, handleerror.ErrorCodeUnknown, "json.SynchronizeData")
	}

  node.wg.Add(len(answers))
	for _, awnswer := range answers {
		var tmpDocument common.Document
		tmpDocument.ID = strconv.Itoa(awnswer.ID)
		tmpDocument.Content = trimEmptyLines(awnswer.Content)
		tmpDocument.Type = common.AnswerType
		tmpDocument.Accepted = awnswer.Accepted

    go node.SynchronizeDocumentData(tmpDocument)
	}

  node.wg.Wait()
}

func (node *Node) StandardData() {
	_, err := node.itemDB.GetListQuestion(10000, 1)
	if err != nil {
		handleerror.WrapErrorf(err, handleerror.ErrorCodeUnknown, "json.StandardData")
	}

	// for _, question := range questions {
	// }
}
func trimEmptyLines(b string) string {
	strs := strings.Split(string(b), "\n")
	str := ""
	for _, s := range strs {
		if len(strings.TrimSpace(s)) == 0 {
			continue
		}
		str += s + " "
	}
	str = strings.TrimSuffix(str, "\n")

	return str
}

func (node *Node) GetAnswerForTest() {
	answers, err := node.itemDB.GetListAnswer(1, 1)
	if err != nil {
		handleerror.WrapErrorf(err, handleerror.ErrorCodeUnknown, "json.SynchronizeData")
	}
	data := answers[0].Content
	fmt.Println(trimEmptyLines(data))
}
