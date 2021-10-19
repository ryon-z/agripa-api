package auction

import (
	modelsV2 "agripa-api/models/v2"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdjAuctionQuantityBarGraphList 정산 경락 거래량 바 차트 조회
// @Summary 정산 경락 거래량 바 차트 조회
// @Description 정산 경락 거래량 바 차트 조회합니다
// @Tags v2 정산 경락 거래량
// @Accept  json
// @Produce  json
// @Router /v2/auction/adj/quantity/bar-graph/{std-item-code} [get]
// @Param std-item-code path string true "표준품목코드"
// @Success 200 {object} v2.AdjAuctionQuantityBarGraph
func AdjAuctionQuantityBarGraphList(c *gin.Context) {
	stdItemCode := c.Param("std-item-code")

	c.JSON(http.StatusOK, gin.H{
		"data": modelsV2.GetAdjAuctionQuantityBarGraph(stdItemCode),
	})
}

// AdjAuctionQuantityRecentTop3List 정산 경락 거래량 최근 TOP 3 조회
// @Summary 정산 경락 거래량 최근 TOP 3 조회
// @Description 정산 경락 거래량 최근 TOP 3 조회합니다
// @Tags v2 정산 경락 거래량
// @Accept  json
// @Produce  json
// @Router /v2/auction/adj/quantity/recent-top-3 [get]
// @Success 200 {object} v2.AdjAuctionQuantityRecentTop3
func AdjAuctionQuantityRecentTop3List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": modelsV2.GetAdjAuctionQuantityRecentTop3(),
	})
}
