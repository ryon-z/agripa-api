package retail

import (
	"agripa-api/common"
	modelsV2 "agripa-api/models/v2"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RecentPriceList 최근 소매가격 정보 조회
// @Summary 최근 소매가격 정보 조회
// @Description 최근 소매가격 정보 조회 조회합니다
// @Tags v2 소매
// @Accept  json
// @Produce  json
// @Router /v2/retail/price/recent/{std-item-code} [get]
// @Param std-item-code path string true "표준품목코드"
// @Success 200 {object} v2.RetailPrice
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
		"data": modelsV2.GetRecentRetailPrice(examinItemCodes),
	})
}

// PreviousYearPriceList 전년 소매가격 정보 조회
// @Summary 전년 소매가격 정보 조회
// @Description 전년 소매가격 정보 조회합니다
// @Tags v2 소매
// @Accept  json
// @Produce  json
// @Router /v2/retail/price/previous-year/{std-item-code} [get]
// @Param std-item-code path string true "표준품목코드"
// @Success 200 {object} v2.RetailPrice
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
		"data": modelsV2.GetPreviousYearRetailPrice(examinItemCodes),
	})
}

// PriceLineGraphList 소매가격 꺾은선 그래프 정보 조회
// @Summary 소매가격 꺾은선 그래프 정보 조회
// @Description 소매가격 꺾은선 그래프 정보 조회합니다
// @Tags v2 소매
// @Accept  json
// @Produce  json
// @Router /v2/retail/price/line-graph/{std-item-code} [get]
// @Param std-item-code path string true "표준품목코드"
// @Success 200 {object} v2.RetailPriceGraph
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
		"data": modelsV2.GetRetailPriceLineGraph(stdItemCode, examinItemCodes),
	})
}
