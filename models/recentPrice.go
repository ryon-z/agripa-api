package models

import (
	"agripa-api/common"
	"time"
)

// RecentPrice 최근 가격 정보
type RecentPrice struct {
	ItemCode     int    `gorm:"column:ItemCode"`
	ItemKindCode string `gorm:"column:ItemKindCode"`
	GradeCode    string `gorm:"column:GradeCode"`
	AreaName     string `gorm:"column:AreaName"`
	MarketName   string `gorm:"column:MarketName"`
	LastShipDate string `gorm:"column:LastShipDate"`
	LastPrice    int    `gorm:"column:LastPrice"`
	Stat7days    StatPrice
	Stat30days   StatPrice
}

// StatPrice 가격 통계 정보
type StatPrice struct {
	Days     int
	AvgPrice float32 `gorm:"column:AvgPrice"`
	MinPrice int     `gorm:"column:MinPrice"`
	MaxPrice int     `gorm:"column:MaxPrice"`
	SumPrice int     `gorm:"column:SumPrice"`
	CntPrice int     `gorm:"column:CntPrice"`
}

// GetWholeRecentPrice 도매 최근 가격 정보
func GetWholeRecentPrice(itemCode string, itemKindCode string, gradeCode string) []RecentPrice {

	db := common.GetDB()

	var last RecentPrice
	var list []RecentPrice

	db.Raw(`SELECT ItemCode, ItemKindCode , GradeCode, max(ShipDate) as LastShipDate 
	FROM AGRI_WHOLE_PRICE
	where ItemCode = ?
	and ItemKindCode = ?
	and GradeCode = ?`, itemCode, itemKindCode, gradeCode).Scan(&last)

	db.Raw(`SELECT ItemCode, ItemKindCode , GradeCode ,AreaName, MarketName, max(ShipDate) as LastShipDate 
	FROM AGRI_WHOLE_PRICE
	where ItemCode = ?
	and ItemKindCode = ?
	and GradeCode = ?
	AND ShipDate BETWEEN DATE_SUB(?, INTERVAL 10 DAY) AND ?
	GROUP by AreaName, MarketName
	ORDER by AreaName desc, MarketName asc`, itemCode, itemKindCode, gradeCode, last.LastShipDate, last.LastShipDate).Scan(&list)

	for i, item := range list {
		db.Raw(`select ShipPrice as LastPrice
		FROM AGRI_WHOLE_PRICE 
		where ItemCode = ?
		and ItemKindCode = ?
		and GradeCode = ?
		and MarketName = ?
		and ShipDate = ?`, item.ItemCode, item.ItemKindCode, item.GradeCode, item.MarketName, item.LastShipDate).Row().Scan(&(list[i].LastPrice))

		list[i].Stat7days = GetWholeStatPrice(item.ItemCode, item.ItemKindCode, item.GradeCode, item.MarketName, item.LastShipDate, 7)
		list[i].Stat30days = GetWholeStatPrice(item.ItemCode, item.ItemKindCode, item.GradeCode, item.MarketName, item.LastShipDate, 30)
	}

	return list
}

// GetWholeStatPrice 도매 가격 통계 정보
func GetWholeStatPrice(itemCode int, itemKindCode string, gradeCode string, marketName string, baseDate string, days int) StatPrice {
	db := common.GetDB()

	var item StatPrice

	startDate, _ := time.Parse("2006-01-02", baseDate)
	startDate = startDate.AddDate(0, 0, days*-1)

	item.Days = days

	db.Raw(`select AVG(ShipPrice) as AvgPrice, min(ShipPrice) as MinPrice, max(ShipPrice) as MaxPrice, sum(ShipPrice) as SumPrice, count(ShipPrice) as CntPrice
	FROM AGRI_WHOLE_PRICE awp 
	where ItemCode = ?
	and ItemKindCode = ?
	and GradeCode = ?
	and MarketName = ?
	and ShipDate BETWEEN ? and ?;`, itemCode, itemKindCode, gradeCode, marketName, startDate.Format("2006-01-02"), baseDate).Scan(&item)

	return item
}

// GetRetailRecentPrice 소매 최근 가격 정보
func GetRetailRecentPrice(itemCode string, itemKindCode string, gradeCode string, areaName string) []RecentPrice {

	db := common.GetDB()

	var last RecentPrice
	var list []RecentPrice

	db.Raw(`SELECT ItemCode, ItemKindCode , GradeCode, max(ShipDate) as LastShipDate 
	FROM AGRI_RETAIL_PRICE
	where ItemCode = ?
	and ItemKindCode = ?
	and GradeCode = ?
	and AreaName = ?
	and MarketName NOT REGEXP '유통|마트'`, itemCode, itemKindCode, gradeCode, areaName).Scan(&last)

	db.Raw(`SELECT ItemCode, ItemKindCode , GradeCode ,AreaName, MarketName, max(ShipDate) as LastShipDate 
	FROM AGRI_RETAIL_PRICE
	where ItemCode = ?
	and ItemKindCode = ?
	and GradeCode = ?
	and AreaName = ?
	and MarketName NOT REGEXP '유통|마트'
	AND ShipDate BETWEEN DATE_SUB(?, INTERVAL 10 DAY) AND ?
	GROUP by MarketName
	ORDER by MarketName asc`, itemCode, itemKindCode, gradeCode, areaName, last.LastShipDate, last.LastShipDate).Scan(&list)

	for i, item := range list {
		db.Raw(`select ShipPrice as LastPrice
		FROM AGRI_RETAIL_PRICE 
		where ItemCode = ?
		and ItemKindCode = ?
		and GradeCode = ?
		and AreaName = ?
		and MarketName = ?
		and ShipDate = ?`, item.ItemCode, item.ItemKindCode, item.GradeCode, item.AreaName, item.MarketName, item.LastShipDate).Row().Scan(&(list[i].LastPrice))

		list[i].Stat7days = GetRetailStatPrice(item.ItemCode, item.ItemKindCode, item.GradeCode, item.AreaName, item.MarketName, item.LastShipDate, 7)
		list[i].Stat30days = GetRetailStatPrice(item.ItemCode, item.ItemKindCode, item.GradeCode, item.AreaName, item.MarketName, item.LastShipDate, 30)
	}

	return list
}

// GetRetailStatPrice 소매 가격 통계 정보
func GetRetailStatPrice(itemCode int, itemKindCode string, gradeCode string, areaName string, marketName string, baseDate string, days int) StatPrice {
	db := common.GetDB()

	var item StatPrice

	startDate, _ := time.Parse("2006-01-02", baseDate)
	startDate = startDate.AddDate(0, 0, days*-1)

	item.Days = days

	db.Raw(`select AVG(ShipPrice) as AvgPrice, min(ShipPrice) as MinPrice, max(ShipPrice) as MaxPrice, sum(ShipPrice) as SumPrice, count(ShipPrice) as CntPrice
	FROM AGRI_RETAIL_PRICE awp 
	where ItemCode = ?
	and ItemKindCode = ?
	and GradeCode = ?
	and AreaName = ?
	and MarketName = ?
	and ShipDate BETWEEN ? and ?;`, itemCode, itemKindCode, gradeCode, areaName, marketName, startDate, baseDate).Scan(&item)

	return item
}
