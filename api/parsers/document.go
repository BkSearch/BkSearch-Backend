package parsers

import (

	"github.com/BkSearch/BkSearch-Backend/common"
	"github.com/gin-gonic/gin"
)

type SearchFilter struct {
	Keyword string `form:"keyword"`
}

func ParserKeyWordFilter(c *gin.Context) (common.QuerySearch, error) {
	var searchFilter SearchFilter

	if err := c.ShouldBind(&searchFilter); err != nil {
		return common.QuerySearch{}, err

	}

  return common.QuerySearch{KeyWord: searchFilter.Keyword}, nil
}
