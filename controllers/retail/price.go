package retail

import (
	models "agripa-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RecentPriceList 소매 최근 가격 정보
// @Summary 소매 최근 가격 정보 조회
// @Description 소매 최근 가격 정보를 조회합니다
// @Tags 소매
// @Accept  json
// @Produce  json
// @Router /retail/price/recent/{itemcode}/{kindcode}/{gradecode}/{areaname} [get]
// @Param itemcode path int true "품목코드"
// @Param kindcode path string true "품종코드"
// @Param gradecode path string true "등급코드"
// @Param areaname path string true "지역명"
// @Success 200 {object} models.RecentPrice
func RecentPriceList(c *gin.Context) {
	itemCode := c.Param("itemcode")
	kindCode := c.Param("kindcode")
	gradeCode := c.Param("gradecode")
	areaName := c.Param("areaname")

	c.JSON(http.StatusOK, gin.H{
		"data": models.GetRetailRecentPrice(itemCode, kindCode, gradeCode, areaName),
	})
}

// MonthPriceList 소매 월 가격 정보
// @Summary 소매 월 가격 정보 조회
// @Description 소매 월 가격 정보를 조회합니다
// @Tags 소매
// @Accept  json
// @Produce  json
// @Router /retail/price/month/{itemcode}/{kindcode}/{gradecode}/{areaname} [get]
// @Param itemcode path int true "품목코드"
// @Param kindcode path string true "품종코드"
// @Param gradecode path string true "등급코드"
// @Param areaname path string true "지역명"
// @Success 200 {object} models.RecentPrice
func MonthPriceList(c *gin.Context) {
	itemCode := c.Param("itemcode")
	kindCode := c.Param("kindcode")
	gradeCode := c.Param("gradecode")
	areaName := c.Param("areaname")

	c.JSON(http.StatusOK, gin.H{
		"data": models.GetRetailMonth(itemCode, kindCode, gradeCode, areaName),
	})
}

// GraphPriceList 소매 가격 정보(그래프용)
// @Summary 소매 가격 정보(그래프용)
// @Description 소매 가격 정보 - 그래프용
// @Tags 소매
// @Accept  json
// @Produce  json
// @Router /retail/price/graph/{itemcode}/{kindcode}/{gradecode}/{areaname}/{graphType} [get]
// @Param itemcode path int true "품목코드"
// @Param kindcode path string true "품종코드"
// @Param gradecode path string true "등급코드"
// @Param areaname path string true "지역명"
// @Param graphType path string true "그래프타입(day,month,year)"
// @Success 200 {object} models.Graph
func GraphPriceList(c *gin.Context) {
	itemCode := c.Param("itemcode")
	kindCode := c.Param("kindcode")
	gradeCode := c.Param("gradecode")
	areaName := c.Param("areaname")
	graphType := c.Param("graphtype")

	c.JSON(http.StatusOK, gin.H{
		"data": models.GetRetailGraphPrice(itemCode, kindCode, gradeCode, areaName, graphType),
	})
}
