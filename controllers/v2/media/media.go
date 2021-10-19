package media

import (
	"agripa-api/common"
	modelsV2 "agripa-api/models/v2"
	"net/http"

	"github.com/gin-gonic/gin"
)

// BlogList 품목별 블로그 조회
// @Summary 품목별 블로그 조회(기준시간으로부터 6개월까지만 조회 가능)
// @Description 품목별 블로그 조회(기준시간으로부터 6개월까지만 조회 가능)
// @Tags v2 미디어
// @Accept  json
// @Produce  json
// @Param std-item-code query int true "표준품목코드"
// @Param base-date-time query string true "기준시간(YYYY-MM-DD HH:mm:ss)"
// @Param start query int false "조회 시작 index(허용값: 최소 0)" mininum(0)
// @Param count-per-page query int false "페이지 당 자료 수(허용값: 최소 0, 최대 50)" mininum(0) maxinum(50)
// @Router /v2/media/blog [get]
// @Success 200 {object} v2.Blog
func BlogList(c *gin.Context) {
	const maxCountPerPage int = 50
	stdItemCode := c.Query("std-item-code")
	baseDateTimeQuery := c.Query("base-date-time") // YYYY-MM-DD HH:mm:ss
	startQuery := c.DefaultQuery("start", "0")
	countPerPageQuery := c.DefaultQuery("count-per-page", "10")

	results, errors := common.ConvertStringToInt(
		map[string]string{
			"start":        startQuery,
			"countPerPage": countPerPageQuery})

	errorInfo := common.IsError(errors)
	if errorInfo.Status {
		c.JSON(http.StatusBadRequest, errorInfo.Message)
		return
	}

	// 표준품목코드와 매핑된 조사가격품목코드 획득
	var examinItemCodes []string
	itemMap := modelsV2.GetItemCodeMap(stdItemCode)
	for _, row := range itemMap {
		if row.ExaminItemCode != "non" && !common.IsStringInArray(row.ExaminItemCode, examinItemCodes) {
			examinItemCodes = append(examinItemCodes, row.ExaminItemCode)
		}
	}

	// 조사가격품목코드에 해당하는 검색어(query) 획득(중복제거)
	var queries []string
	for _, examinItemCode := range examinItemCodes {
		for _, query := range modelsV2.GetMediaQuery("blog", examinItemCode) {
			if !common.IsStringInArray(query, queries) {
				queries = append(queries, query)
			}
		}
	}

	if results["countPerPage"] > maxCountPerPage {
		results["countPerPage"] = maxCountPerPage
	}

	var meta map[string]interface{}
	meta = make(map[string]interface{})
	meta["startIndex"] = results["start"]
	meta["countPerPage"] = results["countPerPage"]
	meta["totalCount"] = modelsV2.GetBlogCount(examinItemCodes, baseDateTimeQuery, queries)

	c.JSON(http.StatusOK, gin.H{
		"data": modelsV2.GetBlogList(
			examinItemCodes,
			results["start"],
			results["countPerPage"],
			baseDateTimeQuery,
			queries),
		"meta": meta,
	})
}

// NewsList 품목별 뉴스 조회
// @Summary 품목별 뉴스 조회(기준시간으로부터 6개월까지만 조회 가능)
// @Description 품목별 뉴스 조회(기준시간으로부터 6개월까지만 조회 가능)
// @Tags v2 미디어
// @Accept  json
// @Produce  json
// @Param std-item-code query int true "표준품목코드"
// @Param base-date-time query string true "기준시간(YYYY-MM-DD HH:mm:ss)"
// @Param start query int false "조회 시작 index(허용값: 최소 0)" mininum(0)
// @Param count-per-page query int false "페이지 당 자료 수(허용값: 최소 0, 최대 50)" mininum(0) maxinum(50)
// @Router /v2/media/news [get]
// @Success 200 {object} v2.News
func NewsList(c *gin.Context) {
	const maxCountPerPage int = 50
	stdItemCode := c.Query("std-item-code")
	baseDateTimeQuery := c.Query("base-date-time") // YYYY-MM-DD HH:mm:ss
	startQuery := c.DefaultQuery("start", "0")
	countPerPage := c.DefaultQuery("count-per-page", "10")

	results, errors := common.ConvertStringToInt(
		map[string]string{
			"start":        startQuery,
			"countPerPage": countPerPage})

	errorInfo := common.IsError(errors)
	if errorInfo.Status {
		c.JSON(http.StatusBadRequest, errorInfo.Message)
		return
	}

	// 표준품목코드와 매핑된 조사가격품목코드 획득
	var examinItemCodes []string
	itemMap := modelsV2.GetItemCodeMap(stdItemCode)
	for _, row := range itemMap {
		if row.ExaminItemCode != "non" && !common.IsStringInArray(row.ExaminItemCode, examinItemCodes) {
			examinItemCodes = append(examinItemCodes, row.ExaminItemCode)
		}
	}

	// 조사가격품목코드에 해당하는 검색어(query) 획득(중복제거)
	var queries []string
	for _, examinItemCode := range examinItemCodes {
		for _, query := range modelsV2.GetMediaQuery("news", examinItemCode) {
			if !common.IsStringInArray(query, queries) {
				queries = append(queries, query)
			}
		}
	}

	if results["countPerPage"] > maxCountPerPage {
		results["countPerPage"] = maxCountPerPage
	}

	var meta map[string]interface{}
	meta = make(map[string]interface{})
	meta["startIndex"] = results["start"]
	meta["countPerPage"] = results["countPerPage"]
	meta["totalCount"] = modelsV2.GetNewsCount(examinItemCodes, baseDateTimeQuery, queries)

	c.JSON(http.StatusOK, gin.H{
		"data": modelsV2.GetNewsList(
			examinItemCodes,
			results["start"],
			results["countPerPage"],
			baseDateTimeQuery,
			queries),
		"meta": meta,
	})
}

// YoutubeList 품목별 유튜브 조회
// @Summary 품목별 유튜브 조회(기준시간으로부터 6개월까지만 조회 가능)
// @Description 품목별 유튜브 조회(기준시간으로부터 6개월까지만 조회 가능)
// @Tags v2 미디어
// @Accept  json
// @Produce  json
// @Param std-item-code query int true "표준품목코드"
// @Param base-date-time query string true "기준시간(YYYY-MM-DD HH:mm:ss)"
// @Param start query int false "조회 시작 index(허용값: 최소 0)" mininum(0)
// @Param count-per-page query int false "페이지 당 자료 수(허용값: 최소 0, 최대 50)" mininum(0) maxinum(50)
// @Router /v2/media/youtube [get]
// @Success 200 {object} v2.Youtube
func YoutubeList(c *gin.Context) {
	const maxCountPerPage int = 50
	stdItemCode := c.Query("std-item-code")
	baseDateTimeQuery := c.Query("base-date-time") // YYYY-MM-DD HH:mm:ss
	startQuery := c.DefaultQuery("start", "0")
	countPerPage := c.DefaultQuery("count-per-page", "10")

	results, errors := common.ConvertStringToInt(
		map[string]string{
			"start":        startQuery,
			"countPerPage": countPerPage})

	errorInfo := common.IsError(errors)
	if errorInfo.Status {
		c.JSON(http.StatusBadRequest, errorInfo.Message)
		return
	}

	// 표준품목코드와 매핑된 조사가격품목코드 획득
	var examinItemCodes []string
	itemMap := modelsV2.GetItemCodeMap(stdItemCode)
	for _, row := range itemMap {
		if row.ExaminItemCode != "non" && !common.IsStringInArray(row.ExaminItemCode, examinItemCodes) {
			examinItemCodes = append(examinItemCodes, row.ExaminItemCode)
		}
	}

	// 조사가격품목코드에 해당하는 검색어(query) 획득(중복제거)
	var queries []string
	for _, examinItemCode := range examinItemCodes {
		for _, query := range modelsV2.GetMediaQuery("youtube", examinItemCode) {
			if !common.IsStringInArray(query, queries) {
				queries = append(queries, query)
			}
		}
	}

	if results["countPerPage"] > maxCountPerPage {
		results["countPerPage"] = maxCountPerPage
	}

	var meta map[string]interface{}
	meta = make(map[string]interface{})
	meta["startIndex"] = results["start"]
	meta["countPerPage"] = results["countPerPage"]
	meta["totalCount"] = modelsV2.GetYoutubeCount(examinItemCodes, baseDateTimeQuery, queries)

	c.JSON(http.StatusOK, gin.H{
		"data": modelsV2.GetYoutubeList(
			examinItemCodes,
			results["start"],
			results["countPerPage"],
			baseDateTimeQuery,
			queries),
		"meta": meta,
	})
}
