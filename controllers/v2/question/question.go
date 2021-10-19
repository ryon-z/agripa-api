package question

import (
	"agripa-api/common"
	modelsV2 "agripa-api/models/v2"
	"net/http"

	"github.com/gin-gonic/gin"
)

// QuestionPost 문의사항 삽입
// @Summary 문의사항 삽입
// @Description 문의사항 삽입합니다
// @Tags v2 문의사항
// @Accept  json
// @Produce  json
// @Router /v2/question [post]
// @Param question body v2.Question true "문의사항"
// @Success 200 {object} v2.Question
func QuestionPost(c *gin.Context) {
	var question modelsV2.Question

	// 인자 로드 실패
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, err.Error)
		return
	}

	// regiDate 날자 형태 체크
	if isDateString := common.IsDateString(question.RegiDate, "-"); !isDateString {
		c.JSON(http.StatusBadRequest, "regiDate가 날짜 형태(yyyy-mm-dd)가 아닙니다.")
		return
	}

	// email 형태 체크
	if isValidEmail := common.IsValidEmail(question.Email); !isValidEmail {
		c.JSON(http.StatusBadRequest, "email이 유효하지 않습니다.")
		return
	}

	// Insert 에러
	if err := modelsV2.InsertQuestion(question); err != nil {
		c.JSON(http.StatusBadRequest, err.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "success",
	})
}
