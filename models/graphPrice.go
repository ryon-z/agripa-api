package models

import (
	"agripa-api/common"
	"fmt"
)

// Graph 월 정보
type Graph struct {
	ItemCode     int    `gorm:"column:ItemCode"`
	ItemKindCode string `gorm:"column:ItemKindCode"`
	GradeCode    string `gorm:"column:GradeCode"`
	AreaName     string `gorm:"column:AreaName"`
	LastShipDate string `gorm:"column:LastShipDate"`

	RangeType  string
	RangeLabel []string
	RangeMin   int `gorm:"column:MinPrice"`
	RangeMax   int `gorm:"column:MaxPrice"`
	RangeStep  int

	GraphLine []GraphLine
}

// GraphLine 그래프 라인 정보
type GraphLine struct {
	AreaName   string `gorm:"column:AreaName"`
	MarketName string `gorm:"column:MarketName"`

	GraphData []GraphData
}

// GraphData 월 가격 정보
type GraphData struct {
	Range string  `gorm:"column:Label" json:"x"`
	Price float32 `gorm:"column:Price" json:"y"`
}

// GetWholeGraphPrice 도매 그래프용 가격 정보
func GetWholeGraphPrice(itemCode string, itemKindCode string, gradeCode string, rangeType string) Graph {

	db := common.GetDB()

	var root Graph
	var list []GraphLine

	db.Raw(`SELECT ItemCode, ItemKindCode , GradeCode
	, max(ShipDate) as LastShipDate, MIN(ShipPrice) AS MinPrice, MAX(ShipPrice) AS MaxPrice
	FROM AGRI_WHOLE_PRICE
	where ItemCode = ?
	and ItemKindCode = ?
	and GradeCode = ?
	AND ShipPrice > 0`, itemCode, itemKindCode, gradeCode).Scan(&root)

	root.RangeType = rangeType

	i := root.RangeMin / 10000

	if i > 0 {
		root.RangeMin = i * 10000
	} else {
		root.RangeMin = (root.RangeMin / 1000) * 1000
	}

	j := (root.RangeMax + 5000) / 10000

	if j > 0 {
		root.RangeMax = j * 10000
	} else {
		root.RangeMax = ((root.RangeMax + 500) / 1000) * 1000
	}

	d := root.RangeMax - root.RangeMin

	if d >= 500000 {
		root.RangeStep = 500000
	} else if d >= 100000 {
		root.RangeStep = 100000
	} else if d >= 50000 {
		root.RangeStep = 50000
	} else if d >= 10000 {
		root.RangeStep = 10000
	} else if d >= 5000 {
		root.RangeStep = 5000
	} else if d >= 1000 {
		root.RangeStep = 1000
	} else if d >= 500 {
		root.RangeStep = 500
	} else {
		root.RangeStep = 100
	}

	root.RangeLabel = GetWholeGraphLabelList(root.ItemCode, root.ItemKindCode, root.GradeCode, root.LastShipDate, root.RangeType)

	db.Raw(`SELECT ItemCode, ItemKindCode , GradeCode ,AreaName, MarketName, max(ShipDate) as LastShipDate
	FROM AGRI_WHOLE_PRICE
	where ItemCode = ?
	and ItemKindCode = ?
	and GradeCode = ?
	AND ShipDate > DATE_FORMAT(DATE_SUB(?, INTERVAL 10 DAY), '%Y-%m-%d')
	GROUP by AreaName, MarketName
	ORDER by AreaName desc, MarketName asc`, itemCode, itemKindCode, gradeCode, root.LastShipDate).Scan(&list)

	root.GraphLine = list

	for i, item := range root.GraphLine {
		root.GraphLine[i].GraphData = GetWholeGraphPriceList(root.ItemCode, root.ItemKindCode, root.GradeCode, item.MarketName, root.LastShipDate, rangeType)
	}

	return root
}

// GetWholeGraphLabelList 도매 그래프 라벨 리스트
func GetWholeGraphLabelList(itemCode int, itemKindCode string, gradeCode string, LastShipDate string, graphType string) []string {
	db := common.GetDB()

	var list []string

	var sql string

	if graphType == "day" {
		dateFormat := "'%Y-%m-%d'"
		rangeVal := 30

		sql = fmt.Sprintf(`SELECT ShipDate AS Label
		FROM AGRI_WHOLE_PRICE
		WHERE ItemCode = ?
		AND ItemKindCode = ?
		AND GradeCode = ?
		AND ShipDate BETWEEN DATE_FORMAT(DATE_SUB(?, INTERVAL %d %s), %s) AND DATE_FORMAT(?, %s)
		GROUP BY ShipDate
		ORDER BY 1 ASC`, rangeVal, graphType, dateFormat, dateFormat)
	} else if graphType == "month" {
		dateFormat := "'%Y-%m'"
		rangeVal := 12

		sql = fmt.Sprintf(`SELECT BaseMonth AS Label
		FROM AGRI_WHOLE_PRICE_MONTH
		WHERE ItemCode = ?
		AND ItemKindCode = ?
		AND GradeCode = ?
		AND BaseMonth BETWEEN DATE_FORMAT(DATE_SUB(?, INTERVAL %d %s), %s) AND DATE_FORMAT(?, %s)
		GROUP BY BaseMonth
		ORDER BY 1 ASC`, rangeVal, graphType, dateFormat, dateFormat)
	} else if graphType == "year" {
		dateFormat := "'%Y'"
		rangeVal := 10

		sql = fmt.Sprintf(`SELECT ShipYear AS Label
		FROM AGRI_WHOLE_PRICE_MONTH
		WHERE ItemCode = ?
		AND ItemKindCode = ?
		AND GradeCode = ?
		AND ShipYear BETWEEN DATE_FORMAT(DATE_SUB(?, INTERVAL %d %s), %s) AND DATE_FORMAT(?, %s)
		GROUP BY ShipYear
		ORDER BY 1 ASC`, rangeVal, graphType, dateFormat, dateFormat)
	}

	db.Raw(sql, itemCode, itemKindCode, gradeCode, LastShipDate, LastShipDate).Pluck("Label", &list)

	return list
}

// GetWholeGraphPriceList 가격 통계 정보
func GetWholeGraphPriceList(itemCode int, itemKindCode string, gradeCode string, marketName string, LastShipDate string, graphType string) []GraphData {
	db := common.GetDB()

	var list []GraphData

	var sql string

	if graphType == "day" {
		dateFormat := "'%Y-%m-%d'"
		rangeVal := 30

		sql = fmt.Sprintf(`SELECT ShipDate AS Label, ShipPrice AS Price
		FROM AGRI_WHOLE_PRICE
		WHERE ItemCode = ?
		AND ItemKindCode = ?
		AND GradeCode = ?
		AND MarketName = ?
		AND ShipDate BETWEEN DATE_FORMAT(DATE_SUB(?, INTERVAL %d %s), %s) AND DATE_FORMAT(?, %s)
		ORDER BY 1 ASC`, rangeVal, graphType, dateFormat, dateFormat)
	} else if graphType == "month" {
		dateFormat := "'%Y-%m'"
		rangeVal := 12

		sql = fmt.Sprintf(`SELECT BaseMonth AS Label, AvgPrice AS Price
		FROM AGRI_WHOLE_PRICE_MONTH
		WHERE ItemCode = ?
		AND ItemKindCode = ?
		AND GradeCode = ?
		AND MarketName = ?
		AND BaseMonth BETWEEN DATE_FORMAT(DATE_SUB(?, INTERVAL %d %s), %s) AND DATE_FORMAT(?, %s)
		ORDER BY 1 ASC`, rangeVal, graphType, dateFormat, dateFormat)
	} else if graphType == "year" {
		dateFormat := "'%Y'"
		rangeVal := 10

		sql = fmt.Sprintf(`SELECT ShipYear AS Label, AVG(AvgPrice) AS Price
		FROM AGRI_WHOLE_PRICE_MONTH
		WHERE ItemCode = ?
		AND ItemKindCode = ?
		AND GradeCode = ?
		AND MarketName = ?
		AND ShipYear BETWEEN DATE_FORMAT(DATE_SUB(?, INTERVAL %d %s), %s) AND DATE_FORMAT(?, %s)
		GROUP BY ShipYear
		ORDER BY 1 ASC`, rangeVal, graphType, dateFormat, dateFormat)
	}

	db.Raw(sql, itemCode, itemKindCode, gradeCode, marketName, LastShipDate, LastShipDate).Scan(&list)

	return list
}

// GetRetailGraphPrice 소매 그래프용 가격 정보
func GetRetailGraphPrice(itemCode string, itemKindCode string, gradeCode string, areaName string, rangeType string) Graph {

	db := common.GetDB()

	var root Graph
	var list []GraphLine

	db.Raw(`SELECT ItemCode, ItemKindCode , GradeCode, AreaName
	, max(ShipDate) as LastShipDate, MIN(ShipPrice) AS MinPrice, MAX(ShipPrice) AS MaxPrice
	FROM AGRI_RETAIL_PRICE
	where ItemCode = ?
	and ItemKindCode = ?
	and GradeCode = ?
	and AreaName = ?
	and MarketName NOT REGEXP '유통|마트'
	AND ShipPrice > 0`, itemCode, itemKindCode, gradeCode, areaName).Scan(&root)

	root.RangeType = rangeType
	i := root.RangeMin / 10000

	if i > 0 {
		root.RangeMin = i * 10000
	} else {
		root.RangeMin = (root.RangeMin / 1000) * 1000
	}

	j := (root.RangeMax + 5000) / 10000

	if j > 0 {
		root.RangeMax = j * 10000
	} else {
		root.RangeMax = ((root.RangeMax + 500) / 1000) * 1000
	}

	d := root.RangeMax - root.RangeMin

	if d >= 500000 {
		root.RangeStep = 500000
	} else if d >= 100000 {
		root.RangeStep = 100000
	} else if d >= 50000 {
		root.RangeStep = 50000
	} else if d >= 10000 {
		root.RangeStep = 10000
	} else if d >= 5000 {
		root.RangeStep = 5000
	} else if d >= 1000 {
		root.RangeStep = 1000
	} else if d >= 500 {
		root.RangeStep = 500
	} else {
		root.RangeStep = 100
	}

	root.RangeLabel = GetRetailGraphLabelList(root.ItemCode, root.ItemKindCode, root.GradeCode, root.AreaName, root.LastShipDate, root.RangeType)

	db.Raw(`SELECT ItemCode, ItemKindCode , GradeCode ,AreaName, MarketName, max(ShipDate) as LastShipDate
	FROM AGRI_RETAIL_PRICE
	where ItemCode = ?
	and ItemKindCode = ?
	and GradeCode = ?
	and AreaName = ?
	and MarketName NOT REGEXP '유통|마트'
	AND ShipDate > DATE_FORMAT(DATE_SUB(?, INTERVAL 10 DAY), '%Y-%m-%d')
	GROUP by AreaName, MarketName
	ORDER by AreaName desc, MarketName asc`, itemCode, itemKindCode, gradeCode, areaName, root.LastShipDate).Scan(&list)

	root.GraphLine = list

	for i, item := range root.GraphLine {
		root.GraphLine[i].GraphData = GetRetailGraphPriceList(root.ItemCode, root.ItemKindCode, root.GradeCode, item.AreaName, item.MarketName, root.LastShipDate, rangeType)
	}

	return root
}

// GetRetailGraphLabelList 소매 그래프 라벨 리스트
func GetRetailGraphLabelList(itemCode int, itemKindCode string, gradeCode string, areaName string, LastShipDate string, graphType string) []string {
	db := common.GetDB()

	var list []string

	var sql string

	if graphType == "day" {
		dateFormat := "'%Y-%m-%d'"
		rangeVal := 30

		sql = fmt.Sprintf(`SELECT ShipDate AS Label
		FROM AGRI_RETAIL_PRICE
		WHERE ItemCode = ?
		AND ItemKindCode = ?
		AND GradeCode = ?
		AND AreaName = ?
		and MarketName NOT REGEXP '유통|마트'
		AND ShipDate BETWEEN DATE_FORMAT(DATE_SUB(?, INTERVAL %d %s), %s) AND DATE_FORMAT(?, %s)
		GROUP BY ShipDate
		ORDER BY 1 ASC`, rangeVal, graphType, dateFormat, dateFormat)
	} else if graphType == "month" {
		dateFormat := "'%Y-%m'"
		rangeVal := 12

		sql = fmt.Sprintf(`SELECT BaseMonth AS Label
		FROM AGRI_RETAIL_PRICE_MONTH
		WHERE ItemCode = ?
		AND ItemKindCode = ?
		AND GradeCode = ?
		AND AreaName = ?
		and MarketName NOT REGEXP '유통|마트'
		AND BaseMonth BETWEEN DATE_FORMAT(DATE_SUB(?, INTERVAL %d %s), %s) AND DATE_FORMAT(?, %s)
		GROUP BY BaseMonth
		ORDER BY 1 ASC`, rangeVal, graphType, dateFormat, dateFormat)
	} else if graphType == "year" {
		dateFormat := "'%Y'"
		rangeVal := 10

		sql = fmt.Sprintf(`SELECT ShipYear AS Label
		FROM AGRI_RETAIL_PRICE_MONTH
		WHERE ItemCode = ?
		AND ItemKindCode = ?
		AND GradeCode = ?
		AND AreaName = ?
		and MarketName NOT REGEXP '유통|마트'
		AND ShipYear BETWEEN DATE_FORMAT(DATE_SUB(?, INTERVAL %d %s), %s) AND DATE_FORMAT(?, %s)
		GROUP BY ShipYear
		ORDER BY 1 ASC`, rangeVal, graphType, dateFormat, dateFormat)
	}

	db.Raw(sql, itemCode, itemKindCode, gradeCode, areaName, LastShipDate, LastShipDate).Pluck("Label", &list)

	return list
}

// GetRetailGraphPriceList 가격 통계 정보
func GetRetailGraphPriceList(itemCode int, itemKindCode string, gradeCode string, areaName string, marketName string, LastShipDate string, graphType string) []GraphData {
	db := common.GetDB()

	var list []GraphData

	var sql string

	if graphType == "day" {
		dateFormat := "'%Y-%m-%d'"
		rangeVal := 30

		sql = fmt.Sprintf(`SELECT ShipDate AS Label, ShipPrice AS Price
		FROM AGRI_RETAIL_PRICE
		WHERE ItemCode = ?
		AND ItemKindCode = ?
		AND GradeCode = ?
		AND AreaName = ?
		AND MarketName = ?
		AND ShipDate BETWEEN DATE_FORMAT(DATE_SUB(?, INTERVAL %d %s), %s) AND DATE_FORMAT(?, %s)
		ORDER BY 1 ASC`, rangeVal, graphType, dateFormat, dateFormat)
	} else if graphType == "month" {
		dateFormat := "'%Y-%m'"
		rangeVal := 12

		sql = fmt.Sprintf(`SELECT BaseMonth AS Label, AvgPrice AS Price
		FROM AGRI_RETAIL_PRICE_MONTH
		WHERE ItemCode = ?
		AND ItemKindCode = ?
		AND GradeCode = ?
		AND AreaName = ?
		AND MarketName = ?
		AND BaseMonth BETWEEN DATE_FORMAT(DATE_SUB(?, INTERVAL %d %s), %s) AND DATE_FORMAT(?, %s)
		ORDER BY 1 ASC`, rangeVal, graphType, dateFormat, dateFormat)
	} else if graphType == "year" {
		dateFormat := "'%Y'"
		rangeVal := 10

		sql = fmt.Sprintf(`SELECT ShipYear AS Label, AVG(AvgPrice) AS Price
		FROM AGRI_RETAIL_PRICE_MONTH
		WHERE ItemCode = ?
		AND ItemKindCode = ?
		AND GradeCode = ?
		AND AreaName = ?
		AND MarketName = ?
		AND ShipYear BETWEEN DATE_FORMAT(DATE_SUB(?, INTERVAL %d %s), %s) AND DATE_FORMAT(?, %s)
		GROUP BY ShipYear
		ORDER BY 1 ASC`, rangeVal, graphType, dateFormat, dateFormat)
	}

	db.Raw(sql, itemCode, itemKindCode, gradeCode, areaName, marketName, LastShipDate, LastShipDate).Scan(&list)

	return list
}
