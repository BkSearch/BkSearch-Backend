package api

import (
	"fmt"
	"net/http"

	"github.com/BkSearch/BkSearch-Backend/api/parsers"
	"github.com/gin-gonic/gin"
)

func (a *API) SearchDocument(c *gin.Context) {
	query, err := parsers.ParserKeyWordFilter(c)
	if err != nil {
		fmt.Println(err)
		return
	}
  searchResult, err := a.ES.Search(&query.KeyWord)
  if err != nil {
    fmt.Println(err)
    return
  }
  data, err := a.ItemDB.GetListDocumentDetailAPI(searchResult)
  if err != nil {
    fmt.Println(err)
    return
  }

  c.JSON(http.StatusOK, data)
}
