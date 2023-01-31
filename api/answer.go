package api

import (
	"fmt"
	"net/http"

	"github.com/BkSearch/BkSearch-Backend/api/parsers"
	"github.com/gin-gonic/gin"
)

func (a *API) GetListAnswer(c *gin.Context) {
	answers, err := a.ItemDB.GetListAnswer(100, 1)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, answers)
}

func (a *API) GetAnswerDetail(c *gin.Context) {
	query, err := parsers.ParseAnswerFilter(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := a.ItemDB.GetAnswerDetailAPI(query)
	if err != nil {
		fmt.Println(err)
		return
	}
  c.JSON(http.StatusOK, data)
}
