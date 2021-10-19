package models

import (
	"agripa-api/common"
	"fmt"
	"strings"
)

// Category 카테고리 정보
type Category struct {
	CateCode int    `gorm:"column:CateCode"`
	CateName string `gorm:"column:CateName"`
}

// Item 품목 정보
type Item struct {
	ItemCode int    `gorm:"column:ItemCode"`
	ItemName string `gorm:"column:ItemName"`
}

// ItemKind 품종 정보
type ItemKind struct {
	ItemCode     int    `gorm:"column:ItemCode"`
	ItemName     string `gorm:"column:ItemName"`
	ItemKindCode string `gorm:"column:ItemKindCode"`
	ItemKindName string `gorm:"column:ItemKindName"`
}

// ItemGrade 등급 정보
type ItemGrade struct {
	GradeCode  string  `gorm:"column:GradeCode"`
	GradeName  string  `gorm:"column:GradeName"`
	GradeType  int     `gorm:"column:GradeType"`
	GradeOrder int     `gorm:"column:GradeOrder"`
	ShipAmt    float32 `gorm:"column:ShipAmt"`
	ShipUnit   string  `gorm:"column:ShipUnit"`
}

// Area 지역
type Area struct {
	AreaName string `gorm:"column:AreaName"`
}

// GetCategoryList 카테고리 리스트 조회
func GetCategoryList(shipType string) []Category {
	db := common.GetDB()

	var list []Category

	db.Raw(`SELECT ag.CateCode, ag.CateName
	FROM AGRI_ITEM_MNG aim
	INNER JOIN AGRI_CATEGORY ag on aim.CateCode = ag.CateCode
	WHERE aim.ShipType = ?
	AND ag.IsDisplay = 1
	GROUP BY ag.CateCode, ag.CateName`, shipType).Scan(&list)

	return list
}

// GetItemList 품목 리스트 조회
func GetItemList(shipType string, cateCode string) []Item {
	db := common.GetDB()

	var list []Item

	db.Raw(`SELECT aim.ItemCode, ai.ItemName
	FROM AGRI_ITEM_MNG aim
	inner join AGRI_ITEM ai on ai.ItemCode = aim.ItemCode and ai.ItemKindCode = aim.ItemKindCode
	where aim.isDisplay = 1
	and aim.ShipType = ?
	and aim.CateCode = ?
	GROUP BY aim.ItemCode`, shipType, cateCode).Scan(&list)

	return list
}

// GetItemKindList 품종 리스트 조회
func GetItemKindList(shipType string, itemCode string) []ItemKind {
	db := common.GetDB()

	var list []ItemKind

	db.Raw(`SELECT aim.ItemCode ,ai.ItemName ,aim.ItemKindCode ,ai.ItemKindName 
	FROM AGRI_ITEM_MNG aim
	inner join AGRI_ITEM ai on ai.ItemCode = aim.ItemCode and ai.ItemKindCode = aim.ItemKindCode
	where aim.isDisplay = 1
	and aim.ShipType = ?
	and aim.ItemCode = ?
	GROUP BY aim.ItemCode ,aim.ItemKindCode`, shipType, itemCode).Scan(&list)

	return list
}

// GetItemGradeList 등급 리스트 조회
func GetItemGradeList(shipType string, itemCode string, kindCode string) []ItemGrade {
	db := common.GetDB()

	var list []ItemGrade

	db.Raw(`SELECT aim.GradeCode , ag.GradeName ,ag.GradeType ,ag.GradeOrder, aim.ShipAmt, aim.ShipUnit
	FROM AGRI_ITEM_MNG aim
	inner join AGRI_GRADE ag on ag.GradeCode = aim.GradeCode
	where aim.isDisplay = 1
	and aim.ShipType = ?
	and aim.ItemCode = ?
	and aim.ItemKindCode = ?
	GROUP BY aim.GradeCode
	ORDER BY GradeType, GradeOrder`, shipType, itemCode, kindCode).Scan(&list)

	return list
}

// GetAreaList 지역 리스트 조회
func GetAreaList(ShipType string, itemCode string, kindCode string, gradeCode string) []Area {
	db := common.GetDB()

	var list []Area

	sql := fmt.Sprintf(`SELECT aa.AreaName, aa.SortOrder
	FROM AGRI_%s_PRICE arp 
	INNER JOIN AGRI_AREA aa ON arp.AreaName = aa.AreaName
	WHERE ItemCode = ?
	AND ItemKindCode = ?
	AND GradeCode = ?
	AND MarketName NOT REGEXP '유통|마트'
	GROUP BY aa.AreaName, aa.SortOrder
	ORDER BY aa.SortOrder, aa.AreaName`, strings.ToUpper(ShipType))

	db.Raw(sql, itemCode, kindCode, gradeCode).Scan(&list)

	return list
}
