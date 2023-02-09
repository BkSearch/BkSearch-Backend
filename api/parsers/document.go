package parsers

import (
	"fmt"

	"github.com/BkSearch/BkSearch-Backend/common"
	"github.com/gin-gonic/gin"
)

type SearchFilter struct {
	Keyword    string `form:"keyword"`
	SearchType int    `form:"searchtype"`
	Page       int    `form:"page"`
}

func ParserKeyWordFilter(c *gin.Context) (common.QuerySearch, error) {
	var searchFilter SearchFilter

	if err := c.ShouldBind(&searchFilter); err != nil {
		return common.QuerySearch{}, err

	}
	fmt.Println(searchFilter)

	return common.QuerySearch{KeyWord: searchFilter.Keyword, SearchType: searchFilter.SearchType, Page: searchFilter.Page}, nil
}
