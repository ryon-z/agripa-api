package whole

import (
	models "agripa-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RecentPriceList 도매 최근 가격 정보
// @Summary 도매 최근 가격 정보 조회
// @Description 도매 최근 가격 정보를 조회합니다
// @Tags 도매
// @Accept  json
// @Produce  json
// @Router /whole/price/recent/{itemcode}/{kindcode}/{gradecode} [get]
// @Param itemcode path int true "품목코드"
// @Param kindcode path string true "품종코드"
// @Param gradecode path string true "등급코드"
// @Success 200 {object} models.RecentPrice
func RecentPriceList(c *gin.Context) {
	itemCode := c.Param("itemcode")
	kindCode := c.Param("kindcode")
	gradeCode := c.Param("gradecode")

	c.JSON(http.StatusOK, gin.H{
		"data": models.GetWholeRecentPrice(itemCode, kindCode, gradeCode),
	})
}

// MonthPriceList 도매 월 가격 정보
// @Summary 도매 월 가격 정보 조회
// @Description 도매 전월 기준 3년간 가격 정보를 조회합니다
// @Tags 도매
// @Accept  json
// @Produce  json
// @Router /whole/price/month/{itemcode}/{kindcode}/{gradecode} [get]
// @Param itemcode path int true "품목코드"
// @Param kindcode path string true "품종코드"
// @Param gradecode path string true "등급코드"
// @Success 200 {object} models.MonthInfo
func MonthPriceList(c *gin.Context) {
	itemCode := c.Param("itemcode")
	kindCode := c.Param("kindcode")
	gradeCode := c.Param("gradecode")

	c.JSON(http.StatusOK, gin.H{
		"data": models.GetWholeMonth(itemCode, kindCode, gradeCode),
	})
}

// GraphPriceList 도매 가격 정보(그래프용)
// @Summary 도매 가격 정보(그래프용)
// @Description 도매 가격 정보 - 그래프용
// @Tags 도매
// @Accept  json
// @Produce  json
// @Router /whole/price/graph/{itemcode}/{kindcode}/{gradecode}/{graphType} [get]
// @Param itemcode path int true "품목코드"
// @Param kindcode path string true "품종코드"
// @Param gradecode path string true "등급코드"
// @Param graphType path string true "그래프타입(day,month,year)"
// @Success 200 {object} models.Graph
func GraphPriceList(c *gin.Context) {
	itemCode := c.Param("itemcode")
	kindCode := c.Param("kindcode")
	gradeCode := c.Param("gradecode")
	graphType := c.Param("graphtype")

	c.JSON(http.StatusOK, gin.H{
		"data": models.GetWholeGraphPrice(itemCode, kindCode, gradeCode, graphType),
	})
}
