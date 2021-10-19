package whole

import (
	"agripa-api/common"
	modelsV2 "agripa-api/models/v2"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RecentPriceList 최근 도매가격 정보 조회
// @Summary 최근 도매가격 정보 조회
// @Description 최근 도매가격 정보 조회 조회합니다
// @Tags v2 도매
// @Accept  json
// @Produce  json
// @Router /v2/whole/price/recent/{std-item-code} [get]
// @Param std-item-code path string true "표준품목코드"
// @Success 200 {object} v2.WholePrice
func RecentPriceList(c *gin.Context) {
	stdItemCode := c.Param("std-item-code")

	var examinItemCodes []string
	itemMap := modelsV2.GetItemCodeMap(stdItemCode)
	for _, row := range itemMap {
		if row.ExaminItemCode != "non" && !common.IsStringInArray(row.ExaminItemCode, examinItemCodes) {
			examinItemCodes = append(examinItemCodes, row.ExaminItemCode)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": modelsV2.GetRecentWholePrice(examinItemCodes),
	})
}

// PreviousYearPriceList 전년 도매가격 정보 조회
// @Summary 전년 도매가격 정보 조회
// @Description 전년 도매가격 정보 조회합니다
// @Tags v2 도매
// @Accept  json
// @Produce  json
// @Router /v2/whole/price/previous-year/{std-item-code} [get]
// @Param std-item-code path string true "표준품목코드"
// @Success 200 {object} v2.WholePrice
func PreviousYearPriceList(c *gin.Context) {
	stdItemCode := c.Param("std-item-code")

	var examinItemCodes []string
	itemMap := modelsV2.GetItemCodeMap(stdItemCode)
	for _, row := range itemMap {
		if row.ExaminItemCode != "non" && !common.IsStringInArray(row.ExaminItemCode, examinItemCodes) {
			examinItemCodes = append(examinItemCodes, row.ExaminItemCode)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": modelsV2.GetPreviousYearWholePrice(examinItemCodes),
	})
}

// PriceLineGraphList 도매가격 꺾은선 그래프 정보 조회
// @Summary 도매가격 꺾은선 그래프 정보 조회
// @Description 도매가격 꺾은선 그래프 정보 조회합니다
// @Tags v2 도매
// @Accept  json
// @Produce  json
// @Router /v2/whole/price/line-graph/{std-item-code} [get]
// @Param std-item-code path string true "표준품목코드"
// @Success 200 {object} v2.WholePriceGraph
func PriceLineGraphList(c *gin.Context) {
	stdItemCode := c.Param("std-item-code")

	var examinItemCodes []string
	itemMap := modelsV2.GetItemCodeMap(stdItemCode)
	for _, row := range itemMap {
		if row.ExaminItemCode != "non" && !common.IsStringInArray(row.ExaminItemCode, examinItemCodes) {
			examinItemCodes = append(examinItemCodes, row.ExaminItemCode)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": modelsV2.GetWholePriceLineGraph(stdItemCode, examinItemCodes),
	})
}
