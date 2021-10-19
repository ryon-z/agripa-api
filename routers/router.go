package routers

import (
	md "agripa-api/controllers/media"
	rc "agripa-api/controllers/retail"
	v2a "agripa-api/controllers/v2/auction"
	v2c "agripa-api/controllers/v2/code"
	v2m "agripa-api/controllers/v2/media"
	v2q "agripa-api/controllers/v2/question"
	v2r "agripa-api/controllers/v2/retail"
	v2t "agripa-api/controllers/v2/trade"
	v2w "agripa-api/controllers/v2/whole"
	wc "agripa-api/controllers/whole"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetRouter 라우트 설정
func SetRouter(r *gin.Engine) *gin.Engine {

	// 공통
	common := r.Group("/common")
	{
		common.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"data": "ok",
			})
		})
	}

	// 도매
	whole := r.Group("/whole")
	{
		code := whole.Group("")
		{
			code.GET("/category", wc.CategoryList)
			code.GET("/item/:catecode", wc.ItemList)
			code.GET("/kind/:itemcode", wc.ItemKindList)
			code.GET("/grade/:itemcode/:kindcode", wc.ItemGradeList)
		}

		price := whole.Group("/price")
		{
			price.GET("/recent/:itemcode/:kindcode/:gradecode", wc.RecentPriceList)
			price.GET("/month/:itemcode/:kindcode/:gradecode", wc.MonthPriceList)
			price.GET("/graph/:itemcode/:kindcode/:gradecode/:graphtype", wc.GraphPriceList)
		}
	}

	// 소매
	retail := r.Group("/retail")
	{
		code := retail.Group("")
		{
			code.GET("/category", rc.CategoryList)
			code.GET("/item/:catecode", rc.ItemList)
			code.GET("/kind/:itemcode", rc.ItemKindList)
			code.GET("/grade/:itemcode/:kindcode", rc.ItemGradeList)
			code.GET("/grade/:itemcode/:kindcode/:gradecode", rc.AreaList)
		}

		price := retail.Group("/price")
		{
			price.GET("/recent/:itemcode/:kindcode/:gradecode/:areaname", rc.RecentPriceList)
			price.GET("/month/:itemcode/:kindcode/:gradecode/:areaname", rc.MonthPriceList)
			price.GET("/graph/:itemcode/:kindcode/:gradecode/:areaname/:graphtype", rc.GraphPriceList)
		}
	}

	// 미디어
	media := r.Group("/media")
	{
		mediaDetail := media.Group("")
		{
			mediaDetail.GET("/blog", md.BlogList)
			mediaDetail.GET("/news", md.NewsList)
			mediaDetail.GET("/youtube", md.YoutubeList)
			mediaDetail.GET("/all-youtube", md.AllYoutubeList)
		}
	}

	// API version 2
	v2 := r.Group("/v2")
	{
		// 코드
		code := v2.Group("/code")
		{
			code.GET("/std-item-codes", v2c.StdItemKeywordList)
			code.GET("/item-code-map", v2c.ItemCodeMapList)
		}

		// 도매
		whole := v2.Group("/whole")
		{
			whole.GET("/price/recent/:std-item-code", v2w.RecentPriceList)
			whole.GET("/price/previous-year/:std-item-code", v2w.PreviousYearPriceList)
			whole.GET("/price/line-graph/:std-item-code", v2w.PriceLineGraphList)
		}

		// 소매
		retail := v2.Group("/retail")
		{
			retail.GET("/price/recent/:std-item-code", v2r.RecentPriceList)
			retail.GET("/price/previous-year/:std-item-code", v2r.PreviousYearPriceList)
			retail.GET("/price/line-graph/:std-item-code", v2r.PriceLineGraphList)
		}

		// 경락
		auction := v2.Group("/auction")
		{
			// 정산
			adjAuction := auction.Group("/adj")
			{
				adjAuction.GET("/quantity/bar-graph/:std-item-code", v2a.AdjAuctionQuantityBarGraphList)
				adjAuction.GET("/quantity/recent-top-3", v2a.AdjAuctionQuantityRecentTop3List)
			}
		}

		// 수출입
		trade := v2.Group("/trade")
		{
			trade.GET("/importation/line-graph/:std-item-code", v2t.ImportationInfoList)
			trade.GET("/exportation/line-graph/:std-item-code", v2t.ExportationInfoList)
			trade.GET("/importation/recent-top-3/", v2t.ImportationRecentTop3List)
			trade.GET("/exportation/recent-top-3/", v2t.ExportationRecentTop3List)
		}

		// 미디어
		media := v2.Group("/media")
		{
			media.GET("/blog", v2m.BlogList)
			media.GET("/news", v2m.NewsList)
			media.GET("/youtube", v2m.YoutubeList)
		}

		// 문의사항
		v2.POST("/question", v2q.QuestionPost)
	}

	return r
}
