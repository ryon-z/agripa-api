package trade

import (
	"agripa-api/common"
	modelsV2 "agripa-api/models/v2"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ExportationInfoList 수출 정보
// @Summary 수출 정보 조회
// @Description 수출 정보 조회합니다
// @Tags v2 수출입
// @Accept  json
// @Produce  json
// @Router /v2/trade/exportation/line-graph/{std-item-code} [get]
// @Param std-item-code path string true "표준품목코드"
// @Success 200 {object} v2.TradeInfoGraph
func ExportationInfoList(c *gin.Context) {
	stdItemCode := c.Param("std-item-code")

	var hskPrdlstCodes []string
	itemMap := modelsV2.GetItemCodeMap(stdItemCode)
	for _, row := range itemMap {
		if !common.IsStringInArray(row.HskPrdlstCode, hskPrdlstCodes) {
			hskPrdlstCodes = append(hskPrdlstCodes, row.HskPrdlstCode)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": modelsV2.GetExportationInfo(stdItemCode, hskPrdlstCodes),
	})
}

// ExportationRecentTop3List 수출 최근 중량 TOP 3 정보
// @Summary 수출 최근 중량 TOP 3 정보 조회
// @Description 수출 최근 중량 TOP 3 정보 조회합니다
// @Tags v2 수출입
// @Accept  json
// @Produce  json
// @Router /v2/trade/exportation/recent-top-3 [get]
// @Success 200 {object} v2.TradeRecentTop3
func ExportationRecentTop3List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": modelsV2.GetTradeRecentTop3("수출"),
	})
}
