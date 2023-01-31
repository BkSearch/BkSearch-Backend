package parsers

import (
	"strconv"

	"github.com/BkSearch/BkSearch-Backend/common"
	"github.com/gin-gonic/gin"
)

type QuestionFilter struct {
	QuestionIndex string `uri:"questionIndex" binding:"required"`
}

func ParseQuesitonFilter(c *gin.Context) (common.QueryQuestion, error) {
	var questionFilter QuestionFilter

	if err := c.ShouldBindUri(&questionFilter); err != nil {
		return common.QueryQuestion{}, err

	}
  index, err := strconv.Atoi(questionFilter.QuestionIndex) 
  if err != nil {
    return common.QueryQuestion{}, err
  }
  return common.QueryQuestion{QuestionIndex: index}, nil
}
