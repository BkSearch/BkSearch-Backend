package api

import (
	"fmt"
	"net/http"

	"github.com/BkSearch/BkSearch-Backend/api/parsers"
	"github.com/gin-gonic/gin"
)

func (a *API) GetListQuestion(c *gin.Context) {
	questions, err := a.ItemDB.GetListQuestion(100, 1)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, questions)
}

func (a *API) GetQuestionDetail(c *gin.Context) {
	query, err := parsers.ParseQuesitonFilter(c)

	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := a.ItemDB.GetQuestionDetailAPI(query)
  if err != nil {
    fmt.Println(err)
    return
  }
  c.JSON(http.StatusOK, data)
}
