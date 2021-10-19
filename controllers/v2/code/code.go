package code

import (
	"agripa-api/common"
	modelsV2 "agripa-api/models/v2"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StdItemKeywordList 표준품목코드 키워드 조회
// @Summary 표준품목코드 키워드 조회
// @Description 표준품목코드 키워드 조회합니다
// @Tags v2 코드
// @Accept  json
// @Produce  json
// @Router /v2/code/std-item-codes [get]
// @Param query query string true "표준품목코드 검색어"
// @Param limit query int false "row 수(허용값: 최소 0, 최대 50)" mininum(0) maxinum(50)
// @Success 200 {object} v2.StdItemKeword
func StdItemKeywordList(c *gin.Context) {
	const maxLimit int = 50
	query := c.Query("query")
	limitQuery := c.DefaultQuery("limit", "5")

	results, errors := common.ConvertStringToInt(
		map[string]string{
			"limit": limitQuery})
	errorInfo := common.IsError(errors)
	if errorInfo.Status {
		c.JSON(http.StatusBadRequest, errorInfo.Message)
		return
	}

	if results["limit"] > maxLimit {
		results["limit"] = maxLimit
	}

	c.JSON(http.StatusOK, gin.H{
		"data": modelsV2.GetStdItemKeyword(
			query, results["limit"]),
	})
}

// ItemCodeMapList 품목코드 맵 조회
// @Summary 품목코드 맵 조회
// @Description 품목코드 맵 조회합니다
// @Tags v2 코드
// @Accept  json
// @Produce  json
// @Router /v2/code/item-code-map [get]
// @Param stdItemCode query string true "표준품목코드"
// @Success 200 {object} v2.ItemCodeMap
func ItemCodeMapList(c *gin.Context) {
	stdItemCode := c.Query("stdItemCode")

	c.JSON(http.StatusOK, gin.H{
		"data": modelsV2.GetItemCodeMap(stdItemCode),
	})
}
