package api

import (
	"fmt"
	"net/http"

	"github.com/BkSearch/BkSearch-Backend/api/parsers"
	// "github.com/BkSearch/BkSearch-Backend/common"
	"github.com/gin-gonic/gin"
  "github.com/BkSearch/BkSearch-Backend/elasticsearch"
  "github.com/BkSearch/BkSearch-Backend/utils"
)

func (a *API) SearchDocument(c *gin.Context) {
	query, err := parsers.ParserKeyWordFilter(c)
	if err != nil {
		fmt.Println(err)
		return
	}
  standardKeyword := utils.DecodeURLString(query.KeyWord)
	var searchResult []elasticsearch.DocumentElasticSearchView
	if query.SearchType == 0 {
    fmt.Println(query.Page)
		searchResult, err = a.ES.Search(&standardKeyword, query.Page)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		searchResult, err = a.ES.SearchMatchPhrase(&standardKeyword, query.Page)
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
