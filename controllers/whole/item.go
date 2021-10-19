package whole

import (
	models "agripa-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var saleType = "whole"

// CategoryList 카테고리 정보 조회
// @Summary 카테고리 정보 조회
// @Description 카테고리 정보를 조회합니다
// @Tags 도매
// @Accept  json
// @Produce  json
// @Router /whole/category [get]
// @Success 200 {object} models.Category
func CategoryList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": models.GetCategoryList(saleType),
	})
}

// ItemList 품목 정보 조회
// @Summary 품목 정보 조회
// @Description 품목 정보를 조회합니다
// @Tags 도매
// @Accept  json
// @Produce  json
// @Router /whole/item/{catecode} [get]
// @Param catecode path int false "부류코드"
// @Success 200 {object} models.Item
func ItemList(c *gin.Context) {
	catecode := c.Param("catecode")

	c.JSON(http.StatusOK, gin.H{
		"data": models.GetItemList(saleType, catecode),
	})
}

// ItemKindList 품종 정보 조회
// @Summary 품종 정보 조회
// @Description 품종 정보를 조회합니다
// @Tags 도매
// @Accept  json
// @Produce  json
// @Router /whole/kind/{itemcode} [get]
// @Param itemcode path int true "품목코드"
// @Success 200 {object} models.ItemKind
func ItemKindList(c *gin.Context) {
	itemcode := c.Param("itemcode")

	c.JSON(http.StatusOK, gin.H{
		"data": models.GetItemKindList(saleType, itemcode),
	})
}

// ItemGradeList 등급 정보 조회
// @Summary 등급 정보 조회
// @Description 등급 정보를 조회합니다
// @Tags 도매
// @Accept  json
// @Produce  json
// @Router /whole/grade/{itemcode}/{kindcode} [get]
// @Param itemcode path int true "품목코드"
// @Param kindcode path string true "품종코드"
// @Success 200 {object} models.ItemGrade
func ItemGradeList(c *gin.Context) {
	itemCode := c.Param("itemcode")
	kindCode := c.Param("kindcode")

	c.JSON(http.StatusOK, gin.H{
		"data": models.GetItemGradeList(saleType, itemCode, kindCode),
	})
}
