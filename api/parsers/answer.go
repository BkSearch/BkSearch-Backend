package parsers

import (
	"strconv"

	"github.com/BkSearch/BkSearch-Backend/common"
	"github.com/gin-gonic/gin"
)

type AnswerFilter struct {
	AnswerIndex string `uri:"answerIndex" binding:"required"`
}

func ParseAnswerFilter(c *gin.Context) (common.QueryAnswer, error) {
	var answerFilter AnswerFilter

	if err := c.ShouldBindUri(&answerFilter); err != nil {
		return common.QueryAnswer{}, err

	}
  index, err := strconv.Atoi(answerFilter.AnswerIndex) 
  if err != nil {
    return common.QueryAnswer{}, err
  }
  return common.QueryAnswer{AnswerIndex: index}, nil
}
