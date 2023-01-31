package api

import (
	"github.com/BkSearch/BkSearch-Backend/db"
	"github.com/BkSearch/BkSearch-Backend/elasticsearch"
	handleerror "github.com/BkSearch/BkSearch-Backend/handle-error"
	"github.com/gin-gonic/gin"
)

type API struct {
	ItemDB *db.ItemDB
	ES     *elasticsearch.StackOverflow
}

type Config struct {
	Version string
	Server  *gin.Engine
	ItemDB  *db.ItemDB
	ES      *elasticsearch.StackOverflow
}

func NewAPI(setup Config) (*API, error) {
  if setup.ItemDB == nil {
    return nil, handleerror.NewErrorf(handleerror.ErrorCodeInvalidArgument, "cannot serve Explorer endpoints without ItemDB")
  }

  if setup.ES == nil {
    return nil, handleerror.NewErrorf(handleerror.ErrorCodeInvalidArgument, "cannot serve Explorer endpoints without Elastic Search")
  }

  a := &API {
    ItemDB : setup.ItemDB,
    ES: setup.ES,
  }

  v1 := setup.Server.Group("/v1")

  v1.GET("/questions", a.GetListQuestion)
  v1.GET("/questions/:questionIndex", a.GetQuestionDetail)
  v1.GET("/answers", a.GetListAnswer)
  v1.GET("/answers/:answerIndex", a.GetAnswerDetail)
  v1.GET("/search", a.SearchDocument)
  // v1.GET("/search/", a.Search())

  return a, nil
}
