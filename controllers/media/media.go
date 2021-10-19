package media

import (
	"agripa-api/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func convertStringToInt(inputs map[string]string) (map[string]int, map[string]string) {
	var results map[string]int
	results = make(map[string]int)
	var errors map[string]string
	errors = make(map[string]string)

	for name, value := range inputs {
		result, err := strconv.Atoi(value)
		if err != nil {
			errors[name] = fmt.Sprintf(
				"%s가 숫자형이 아닙니다. %s: %s",
				name, name, value,
			)
		} else {
			errors[name] = ""
		}
		results[name] = result
	}
	return results, errors
}

// BlogList 품목별 블로그 조회
// @Summary 품목별 블로그 조회(기준시간으로부터 6개월까지만 조회 가능)
// @Description 품목별 블로그 조회(기준시간으로부터 6개월까지만 조회 가능)
// @Tags 미디어
// @Accept  json
// @Produce  json
// @Param itemCode query int true "품목코드"
// @Param baseDateTime query string true "기준시간(YYYY-MM-DD HH:mm:ss)"
// @Param start query int false "조회 시작 index(허용값: 최소 0)" mininum(0)
// @Param countPerPage query int false "페이지 당 자료 수(허용값: 최소 0, 최대 50)" mininum(0) maxinum(50)
// @Router /media/blog [get]
// @Success 200 {object} models.Blog
func BlogList(c *gin.Context) {
	const maxCountPerPage int = 50
	itemCodeQuery := c.Query("itemCode")
	baseDateTimeQuery := c.Query("baseDateTime") // YYYY-MM-DD HH:mm:ss
	startQuery := c.DefaultQuery("start", "0")
	countPerPageQuery := c.DefaultQuery("countPerPage", "10")

	results, errors := convertStringToInt(
		map[string]string{
			"itemCode":     itemCodeQuery,
			"start":        startQuery,
			"countPerPage": countPerPageQuery})

	for _, errorString := range errors {
		if errorString != "" {
			c.JSON(
				http.StatusBadRequest,
				errorString,
			)
			return
		}
	}

	if results["countPerPage"] > maxCountPerPage {
		results["countPerPage"] = maxCountPerPage
	}

	var meta map[string]interface{}
	meta = make(map[string]interface{})
	meta["startIndex"] = results["start"]
	meta["countPerPage"] = results["countPerPage"]
	meta["totalCount"] = models.GetBlogCount(results["itemCode"], baseDateTimeQuery)

	c.JSON(http.StatusOK, gin.H{
		"data": models.GetBlogList(
			results["itemCode"],
			results["start"],
			results["countPerPage"],
			baseDateTimeQuery),
		"meta": meta,
	})
}

// NewsList 품목별 뉴스 조회
// @Summary 품목별 뉴스 조회(기준시간으로부터 6개월까지만 조회 가능)
// @Description 품목별 뉴스 조회(기준시간으로부터 6개월까지만 조회 가능)
// @Tags 미디어
// @Accept  json
// @Produce  json
// @Param itemCode query int true "품목코드"
// @Param baseDateTime query string true "기준시간(YYYY-MM-DD HH:mm:ss)"
// @Param start query int false "조회 시작 index(허용값: 최소 0)" mininum(0)
// @Param countPerPage query int false "페이지 당 자료 수(허용값: 최소 0, 최대 50)" mininum(0) maxinum(50)
// @Router /media/news [get]
// @Success 200 {object} models.News
func NewsList(c *gin.Context) {
	const maxCountPerPage int = 50
	itemCodeQuery := c.Query("itemCode")
	baseDateTimeQuery := c.Query("baseDateTime") // YYYY-MM-DD HH:mm:ss
	startQuery := c.DefaultQuery("start", "0")
	countPerPage := c.DefaultQuery("countPerPage", "10")

	results, errors := convertStringToInt(
		map[string]string{
			"itemCode":     itemCodeQuery,
			"start":        startQuery,
			"countPerPage": countPerPage})

	for _, errorString := range errors {
		if errorString != "" {
			c.JSON(
				http.StatusBadRequest,
				errorString,
			)
			return
		}
	}

	if results["countPerPage"] > maxCountPerPage {
		results["countPerPage"] = maxCountPerPage
	}

	var meta map[string]interface{}
	meta = make(map[string]interface{})
	meta["startIndex"] = results["start"]
	meta["countPerPage"] = results["countPerPage"]
	meta["totalCount"] = models.GetNewsCount(results["itemCode"], baseDateTimeQuery)

	c.JSON(http.StatusOK, gin.H{
		"data": models.GetNewsList(
			results["itemCode"],
			results["start"],
			results["countPerPage"],
			baseDateTimeQuery),
		"meta": meta,
	})
}

// AllYoutubeList 모든 유튜브 영상 조회
// @Summary 모든 유튜브 영상 조회
// @Description 모든 유튜브 영상 조회
// @Tags 미디어
// @Accept  json
// @Produce  json
// @Param baseDateTime query string true "기준시간(YYYY-MM-DD HH:mm:ss)"
// @Param start query int false "조회 시작 index(허용값: 최소 0)" mininum(0)
// @Param countPerPage query int false "페이지 당 자료 수(허용값: 최소 0, 최대 50)" mininum(0) maxinum(50)
// @Param query query string false "검색어(미입력 시 모든 유튜브 검색, 가락시장365만 보고 싶은 경우 '000공통000'입력)"
// @Router /media/all-youtube [get]
// @Success 200 {object} models.AllYoutube
func AllYoutubeList(c *gin.Context) {
	const maxCountPerPage int = 50
	baseDateTimeQuery := c.Query("baseDateTime") // YYYY-MM-DD HH:mm:ss
	startQuery := c.DefaultQuery("start", "0")
	countPerPage := c.DefaultQuery("countPerPage", "10")
	queryQuery := c.Query("query")

	results, errors := convertStringToInt(
		map[string]string{"start": startQuery, "countPerPage": countPerPage})

	for _, errorString := range errors {
		if errorString != "" {
			c.JSON(
				http.StatusBadRequest,
				errorString,
			)
			return
		}
	}

	if results["countPerPage"] > maxCountPerPage {
		results["countPerPage"] = maxCountPerPage
	}

	var meta map[string]interface{}
	meta = make(map[string]interface{})
	meta["startIndex"] = results["start"]
	meta["countPerPage"] = results["countPerPage"]
	meta["totalCount"] = models.GetAllYoutubeCount(queryQuery, baseDateTimeQuery)

	c.JSON(http.StatusOK, gin.H{
		"data": models.GetAllYoutubeList(
			results["start"],
			results["countPerPage"],
			queryQuery,
			baseDateTimeQuery),
		"meta": meta,
	})
}

// YoutubeList 품목별 유튜브 조회
// @Summary 품목별 유튜브 조회(기준시간으로부터 6개월까지만 조회 가능)
// @Description 품목별 유튜브 조회(기준시간으로부터 6개월까지만 조회 가능)
// @Tags 미디어
// @Accept  json
// @Produce  json
// @Param itemCode query int true "품목코드"
// @Param baseDateTime query string true "기준시간(YYYY-MM-DD HH:mm:ss)"
// @Param start query int false "조회 시작 index(허용값: 최소 0)" mininum(0)
// @Param countPerPage query int false "페이지 당 자료 수(허용값: 최소 0, 최대 50)" mininum(0) maxinum(50)
// @Router /media/youtube [get]
// @Success 200 {object} models.Youtube
func YoutubeList(c *gin.Context) {
	const maxCountPerPage int = 50
	itemCodeQuery := c.Query("itemCode")
	baseDateTimeQuery := c.Query("baseDateTime") // YYYY-MM-DD HH:mm:ss
	startQuery := c.DefaultQuery("start", "0")
	countPerPage := c.DefaultQuery("countPerPage", "10")

	results, errors := convertStringToInt(
		map[string]string{
			"itemCode":     itemCodeQuery,
			"start":        startQuery,
			"countPerPage": countPerPage})

	for _, errorString := range errors {
		if errorString != "" {
			c.JSON(
				http.StatusBadRequest,
				errorString,
			)
			return
		}
	}

	if results["countPerPage"] > maxCountPerPage {
		results["countPerPage"] = maxCountPerPage
	}

	var meta map[string]interface{}
	meta = make(map[string]interface{})
	meta["startIndex"] = results["start"]
	meta["countPerPage"] = results["countPerPage"]
	meta["totalCount"] = models.GetYoutubeCount(results["itemCode"], baseDateTimeQuery)

	c.JSON(http.StatusOK, gin.H{
		"data": models.GetYoutubeList(
			results["itemCode"],
			results["start"],
			results["countPerPage"],
			baseDateTimeQuery),
		"meta": meta,
	})
}
