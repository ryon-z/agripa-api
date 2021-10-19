package v2

import (
	"agripa-api/common"
	"fmt"
	"sort"
	"time"
)

// RetailPrice 소매 가격 정보
type RetailPrice struct {
	ExaminDate        string `gorm:"column:ExaminDate"`
	ExaminItemName    string `gorm:"column:ExaminItemName"`
	ExaminItemCode    string `gorm:"column:ExaminItemCode"`
	ExaminSpeciesName string `gorm:"column:ExaminSpeciesName"`
	ExaminSpeciesCode string `gorm:"column:ExaminSpeciesCode"`
	ExaminUnitName    string `gorm:"column:ExaminUnitName"`
	ExaminUnit        string `gorm:"column:ExaminUnit"`
	ExaminGradeName   string `gorm:"column:ExaminGradeName"`
	ExaminGradeCode   string `gorm:"column:ExaminGradeCode"`
	MinPrice          int    `gorm:"column:MinPrice"`
	MaxPrice          int    `gorm:"column:MaxPrice"`
}

// RetailPriceGraph 그래프용 소매 가격 정보
type RetailPriceGraph struct {
	StdItemCode string
	RangeLabel  []string
	GraphLine   []retailPriceGraphLine
}

type retailPriceGraphLine struct {
	ExaminItemName    string
	ExaminSpeciesName string
	ExaminUnitName    string
	ExaminGradeName   string
	PriceType         string
	GraphData         []retailPriceGraphData
}

type retailPriceGraphData struct {
	X string
	Y int
}

// GetRetailPriceLineGraph 소매 가격 꺾은선 그래프 정보 조회
func GetRetailPriceLineGraph(stdItemCode string, examinItemCodes []string) RetailPriceGraph {
	db := common.GetDB()

	var list []RetailPrice
	var graph RetailPriceGraph
	var rangeLabels []string
	var data = map[string]map[string]map[string]map[string]map[string][]retailPriceGraphData{}
	var itemNameMap = map[string]string{}
	var speciesNameMap = map[string]string{}
	var gradeNameMap = map[string]string{}

	examinItemCodesPhrase := common.GetOrPhrase(
		common.GetStringDefaultSlice(len(examinItemCodes), "ExaminItemCode"),
		examinItemCodes,
	)
	durationMonth := 1
	sqlQuery := fmt.Sprintf(`
		SELECT ExaminDate, ExaminItemName, ExaminItemCode, ExaminSpeciesName, ExaminSpeciesCode, 
			ExaminUnitName, ExaminUnit, ExaminGradeName, ExaminGradeCode, MinPrice, MaxPrice FROM MAFRA_RETAIL_PRICE
		WHERE ExaminDate > DATE_SUB(NOW(), INTERVAL %d MONTH)
		AND (%s)
		ORDER BY ExaminDate, ExaminItemCode, ExaminSpeciesCode, ExaminGradeCode ASC
	;`, durationMonth, examinItemCodesPhrase)

	db.Raw(sqlQuery).Scan(&list)

	for _, elem := range list {
		// 고유한 날짜(그래프의 X 축 라벨)만 추출
		if !common.IsStringInArray(elem.ExaminDate, rangeLabels) {
			rangeLabels = append(rangeLabels, elem.ExaminDate)
		}

		// 코드-이름 맵 구축
		itemCode := elem.ExaminItemCode
		itemNameMap[itemCode] = elem.ExaminItemName
		speciesCode := elem.ExaminSpeciesCode
		speciesNameMap[speciesCode] = elem.ExaminSpeciesName
		unitName := elem.ExaminUnitName
		gradeCode := elem.ExaminGradeCode
		gradeNameMap[gradeCode] = elem.ExaminGradeName

		// 그래프 X(날짜), Y(가격)값 추출
		for priceType, price := range map[string]int{"minPrice": elem.MinPrice} {
			graphData := retailPriceGraphData{
				X: elem.ExaminDate,
				Y: price,
			}

			if data[itemCode] == nil {
				data[itemCode] = make(map[string]map[string]map[string]map[string][]retailPriceGraphData)
			}
			if data[itemCode][speciesCode] == nil {
				data[itemCode][speciesCode] = make(map[string]map[string]map[string][]retailPriceGraphData)
			}
			if data[itemCode][speciesCode][unitName] == nil {
				data[itemCode][speciesCode][unitName] = make(map[string]map[string][]retailPriceGraphData)
			}
			if data[itemCode][speciesCode][unitName][gradeCode] == nil {
				data[itemCode][speciesCode][unitName][gradeCode] = make(map[string][]retailPriceGraphData)
			}

			data[itemCode][speciesCode][unitName][gradeCode][priceType] = append(
				data[itemCode][speciesCode][unitName][gradeCode][priceType], graphData)
		}
	}

	// json으로 출력하기 위해 struct 형태로 변환
	var graphLines []retailPriceGraphLine
	for itemCode, valueMap := range data {
		for speciesCode, valueMap2 := range valueMap {
			for unitName, valueMap3 := range valueMap2 {
				for gradeCode, valueMap4 := range valueMap3 {
					for priceType, graphData := range valueMap4 {
						var graphLine retailPriceGraphLine
						graphLine.ExaminItemName = itemNameMap[itemCode]
						graphLine.ExaminSpeciesName = speciesNameMap[speciesCode]
						graphLine.ExaminUnitName = unitName
						graphLine.ExaminGradeName = gradeNameMap[gradeCode]
						graphLine.PriceType = priceType
						graphLine.GraphData = graphData
						graphLines = append(graphLines, graphLine)
					}
				}
			}
		}
	}

	sort.Slice(graphLines, func(i, j int) bool {
		if graphLines[i].ExaminItemName < graphLines[j].ExaminItemName {
			return true
		}
		if graphLines[i].ExaminItemName > graphLines[j].ExaminItemName {
			return false
		}
		if graphLines[i].ExaminSpeciesName < graphLines[j].ExaminSpeciesName {
			return true
		}
		if graphLines[i].ExaminSpeciesName > graphLines[j].ExaminSpeciesName {
			return false
		}
		if graphLines[i].ExaminUnitName < graphLines[j].ExaminUnitName {
			return true
		}
		if graphLines[i].ExaminUnitName > graphLines[j].ExaminUnitName {
			return false
		}
		if graphLines[i].ExaminGradeName < graphLines[j].ExaminGradeName {
			return true
		}
		if graphLines[i].ExaminGradeName > graphLines[j].ExaminGradeName {
			return false
		}
		return graphLines[i].PriceType < graphLines[j].PriceType
	})

	graph.StdItemCode = stdItemCode
	graph.RangeLabel = rangeLabels
	graph.GraphLine = graphLines

	return graph
}

// GetRecentRetailPrice 최근 소매 가격 정보 조회
func GetRecentRetailPrice(examinItemCodes []string) []RetailPrice {
	db := common.GetDB()

	var list []RetailPrice

	examinItemCodesPhrase := common.GetOrPhrase(
		common.GetStringDefaultSlice(len(examinItemCodes), "ExaminItemCode"),
		examinItemCodes,
	)
	sqlQuery := fmt.Sprintf(`
		SELECT ExaminDate, ExaminItemName, ExaminItemCode, ExaminSpeciesName, ExaminSpeciesCode, 
			ExaminUnitName, ExaminUnit, ExaminGradeName, ExaminGradeCode, MinPrice, MaxPrice FROM MAFRA_RETAIL_PRICE
		WHERE (%s)
		AND ExaminDate = (
			SELECT ExaminDate FROM MAFRA_RETAIL_PRICE 
			WHERE (%s) 
			ORDER BY ExaminDate DESC 
			LIMIT 1
		)
	;`, examinItemCodesPhrase, examinItemCodesPhrase)
	fmt.Println(sqlQuery)
	db.Raw(sqlQuery).Scan(&list)

	return list
}

// GetPreviousYearRetailPrice 전년 소매 가격 정보 조회
func GetPreviousYearRetailPrice(examinItemCodes []string) []RetailPrice {
	db := common.GetDB()

	durationMonth := 3
	loc, _ := time.LoadLocation("Asia/Seoul")
	previousYearDate := time.Now().In(loc).AddDate(-1, 0, 0)
	maxDate := time.Now().In(loc).AddDate(-1, durationMonth, 0)
	previousYearDateString := common.GetDateString(previousYearDate, "-")
	maxDateString := common.GetDateString(maxDate, "-")

	var list []RetailPrice

	examinItemCodesPhrase := common.GetOrPhrase(
		common.GetStringDefaultSlice(len(examinItemCodes), "ExaminItemCode"),
		examinItemCodes,
	)
	sqlQuery := fmt.Sprintf(`
		SELECT ExaminDate, ExaminItemName, ExaminItemCode, ExaminSpeciesName, ExaminSpeciesCode, 
			ExaminUnitName, ExaminUnit, ExaminGradeName, ExaminGradeCode, MinPrice, MaxPrice FROM MAFRA_RETAIL_PRICE
		WHERE (%s) 
		AND ExaminDate  = (
			SELECT ExaminDate FROM MAFRA_RETAIL_PRICE
			WHERE (%s) 
			AND ExaminDate >= "%s" 
			AND ExaminDate < "%s" 
			ORDER BY ExaminDate ASC 
			LIMIT 1
		)
	;`, examinItemCodesPhrase, examinItemCodesPhrase, previousYearDateString, maxDateString)
	fmt.Println("sqlQuery", sqlQuery)
	db.Raw(sqlQuery).Scan(&list)

	return list
}
