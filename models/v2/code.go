package v2

import (
	"agripa-api/common"
	"fmt"
)

// StdItemKeword 표준품목코드 키워드
type StdItemKeword struct {
	ItemCode string `gorm:"column:ItemCode"`
	Keyword  string `gorm:"column:Keyword"`
}

// ItemCodeMap 품목코드 맵
type ItemCodeMap struct {
	StdItemCode    string `gorm:"column:StdItemCode"`
	ExaminItemCode string `gorm:"column:ExaminItemCode"`
	HskPrdlstCode  string `gorm:"column:HskPrdlstCode"`
}

// GetStdItemKeyword 표준품목코드 키워드 조회
func GetStdItemKeyword(query string, limit int) []StdItemKeword {
	db := common.GetDB()

	var list []StdItemKeword
	sqlQuery := fmt.Sprintf(`
		SELECT ItemCode, Keyword
		FROM STD_ITEM_KEYWORD
		WHERE IsDisplay = 1
		AND Keyword LIKE "%%%s%%"
		ORDER BY NumSearch DESC, Priority DESC, ItemCode ASC
		LIMIT %d
	;`, query, limit)

	db.Raw(sqlQuery).Scan(&list)

	return list
}

// GetItemCodeMap 품목코드 맵 조회
func GetItemCodeMap(stdItemCode string) []ItemCodeMap {
	db := common.GetDB()

	var list []ItemCodeMap
	sqlQuery := fmt.Sprintf(`
		SELECT StdItemCode, ExaminItemCode, HskPrdlstCode
		FROM ITEM_MAPPING
		WHERE StdItemCode = "%s"
	;`, stdItemCode)

	db.Raw(sqlQuery).Scan(&list)

	return list
}
