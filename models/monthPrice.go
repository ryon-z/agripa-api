package models

import (
	"agripa-api/common"
)

// MonthInfo 월 정보
type MonthInfo struct {
	ItemCode     int    `gorm:"column:ItemCode"`
	ItemKindCode string `gorm:"column:ItemKindCode"`
	GradeCode    string `gorm:"column:GradeCode"`
	AreaName     string `gorm:"column:AreaName"`
	MarketName   string `gorm:"column:MarketName"`
	BaseMonth    string `gorm:"column:BaseMonth"`
	ShipYear     int    `gorm:"column:ShipYear"`
	ShipMonth    int    `gorm:"column:ShipMonth"`

	MonthPriceList []MonthPrice
}

// MonthPrice 월 가격 정보
type MonthPrice struct {
	BaseMonth   string  `gorm:"column:BaseMonth"`
	AvgPrice    float32 `gorm:"column:AvgPrice"`
	MinPrice    int     `gorm:"column:MinPrice"`
	MaxPrice    int     `gorm:"column:MaxPrice"`
	SumPrice    int     `gorm:"column:SumPrice"`
	PriceCount  int     `gorm:"column:PriceCount"`
	MinShipDate string  `gorm:"column:MinShipDate"`
	MaxShipDate string  `gorm:"column:MaxShipDate"`
}

// GetWholeMonth 도매 월 가격 정보
func GetWholeMonth(itemCode string, itemKindCode string, gradeCode string) []MonthInfo {

	db := common.GetDB()

	var last MonthInfo
	var list []MonthInfo

	db.Raw(`SELECT ItemCode ,ItemKindCode ,GradeCode
	, MAX(BaseMonth) AS BaseMonth ,DATE_FORMAT(MAX(MaxShipDate), '%Y') AS ShipYear , DATE_FORMAT(MAX(MaxShipDate), '%c') AS ShipMonth
	FROM AGRI_WHOLE_PRICE_MONTH
	WHERE ItemCode = ?
	AND ItemKindCode = ?
	AND GradeCode = ?
	AND BaseMonth < DATE_FORMAT(NOW(),'%Y-%m')`, itemCode, itemKindCode, gradeCode).Scan(&last)

	db.Raw(`SELECT ItemCode ,ItemKindCode ,GradeCode ,AreaName, MarketName
	, MAX(BaseMonth) AS BaseMonth ,DATE_FORMAT(MAX(MaxShipDate), '%Y') AS ShipYear , DATE_FORMAT(MAX(MaxShipDate), '%c') AS ShipMonth
	FROM AGRI_WHOLE_PRICE_MONTH
	WHERE ItemCode = ?
	AND ItemKindCode = ?
	AND GradeCode = ?
	AND BaseMonth = ?
	GROUP BY AreaName, MarketName
	ORDER by AreaName desc, MarketName asc`, itemCode, itemKindCode, gradeCode, last.BaseMonth).Scan(&list)

	for i, item := range list {
		list[i].MonthPriceList = GetWholeMonthPriceList(item.ItemCode, item.ItemKindCode, item.GradeCode, item.MarketName, item.BaseMonth, item.ShipMonth)
	}

	return list
}

// GetWholeMonthPriceList 가격 통계 정보
func GetWholeMonthPriceList(itemCode int, itemKindCode string, gradeCode string, marketName string, baseMonth string, month int) []MonthPrice {
	db := common.GetDB()

	var list []MonthPrice

	db.Raw(`SELECT *
	FROM AGRI_WHOLE_PRICE_MONTH
	WHERE ItemCode = ?
	AND ItemKindCode = ?
	AND GradeCode = ?
	AND MarketName = ?
	AND ShipMonth = ?
	AND BaseMonth <= ?
	ORDER BY BaseMonth DESC;`, itemCode, itemKindCode, gradeCode, marketName, month, baseMonth).Scan(&list)

	return list
}

// GetRetailMonth 도매 월 가격 정보
func GetRetailMonth(itemCode string, itemKindCode string, gradeCode string, areaName string) []MonthInfo {

	db := common.GetDB()

	var last MonthInfo
	var list []MonthInfo

	db.Raw(`SELECT ItemCode ,ItemKindCode ,GradeCode
	, MAX(BaseMonth) AS BaseMonth ,DATE_FORMAT(MAX(MaxShipDate), '%Y') AS ShipYear , DATE_FORMAT(MAX(MaxShipDate), '%c') AS ShipMonth
	FROM AGRI_RETAIL_PRICE_MONTH
	WHERE ItemCode = ?
	AND ItemKindCode = ?
	AND GradeCode = ?
	AND AreaName = ?
	and MarketName NOT REGEXP '유통|마트'
	AND BaseMonth < DATE_FORMAT(NOW(),'%Y-%m')`, itemCode, itemKindCode, gradeCode, areaName).Scan(&last)

	db.Raw(`SELECT ItemCode ,ItemKindCode ,GradeCode ,AreaName, MarketName
	, MAX(BaseMonth) AS BaseMonth ,DATE_FORMAT(MAX(MaxShipDate), '%Y') AS ShipYear , DATE_FORMAT(MAX(MaxShipDate), '%c') AS ShipMonth
	FROM AGRI_RETAIL_PRICE_MONTH
	WHERE ItemCode = ?
	AND ItemKindCode = ?
	AND GradeCode = ?
	AND AreaName = ?
	and MarketName NOT REGEXP '유통|마트'
	AND BaseMonth = ?
	GROUP BY AreaName, MarketName
	ORDER by AreaName desc, MarketName asc`, itemCode, itemKindCode, gradeCode, areaName, last.BaseMonth).Scan(&list)

	for i, item := range list {
		list[i].MonthPriceList = GetRetailMonthPriceList(item.ItemCode, item.ItemKindCode, item.GradeCode, item.AreaName, item.MarketName, item.BaseMonth, item.ShipMonth)
	}

	return list
}

// GetRetailMonthPriceList 가격 통계 정보
func GetRetailMonthPriceList(itemCode int, itemKindCode string, gradeCode string, areaName string, marketName string, baseMonth string, month int) []MonthPrice {
	db := common.GetDB()

	var list []MonthPrice

	db.Raw(`SELECT *
	FROM AGRI_RETAIL_PRICE_MONTH
	WHERE ItemCode = ?
	AND ItemKindCode = ?
	AND GradeCode = ?
	AND AreaName = ?
	AND MarketName = ?
	AND ShipMonth = ?
	AND BaseMonth <= ?
	ORDER BY BaseMonth DESC;`, itemCode, itemKindCode, gradeCode, areaName, marketName, month, baseMonth).Scan(&list)

	return list
}
