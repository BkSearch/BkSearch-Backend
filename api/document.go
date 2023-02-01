package api

import (
	"fmt"
	"net/http"

	"github.com/BkSearch/BkSearch-Backend/api/parsers"
	"github.com/BkSearch/BkSearch-Backend/common"
	"github.com/gin-gonic/gin"
)

func (a *API) SearchDocument(c *gin.Context) {
	query, err := parsers.ParserKeyWordFilter(c)
	if err != nil {
		fmt.Println(err)
		return
	}
  fmt.Println(query)
	var searchResult []common.Document
	if query.SearchType == 0 {
		searchResult, err = a.ES.Search(&query.KeyWord)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		searchResult, err = a.ES.SearchMatchPhrase(&query.KeyWord)
		if err != nil {
			fmt.Println(err)
			return
		}

  }

	data, err := a.ItemDB.GetListDocumentDetailAPI(searchResult)
	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, data)
}
