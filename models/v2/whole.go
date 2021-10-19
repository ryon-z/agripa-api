package v2

import (
	"agripa-api/common"
	"fmt"
	"sort"
	"time"
)

// WholePrice 도매 가격 정보
type WholePrice struct {
	ExaminDate        string `gorm:"column:ExaminDate"`
	ExaminItemName    string `gorm:"column:ExaminItemName"`
	ExaminItemCode    string `gorm:"column:ExaminItemCode"`
	ExaminSpeciesName string `gorm:"column:ExaminSpeciesName"`
	ExaminSpeciesCode string `gorm:"column:ExaminSpeciesCode"`
	ExaminUnitName    string `gorm:"column:ExaminUnitName"`
	ExaminUnit        string `gorm:"column:ExaminUnit"`
	ExaminGradeName   string `gorm:"column:ExaminGradeName"`
	ExaminGradeCode   string `gorm:"column:ExaminGradeCode"`
	Price             int    `gorm:"column:Price"`
}

// WholePriceGraph 그래프용 도매 가격 정보
type WholePriceGraph struct {
	StdItemCode string
	RangeLabel  []string
	GraphLine   []wholePriceGraphLine
}

type wholePriceGraphLine struct {
	ExaminItemName    string
	ExaminSpeciesName string
	ExaminUnitName    string
	ExaminGradeName   string
	GraphData         []wholePriceGraphData
}

type wholePriceGraphData struct {
	X string
	Y int
}

// GetWholePriceLineGraph 도매 가격 꺾은선 그래프 정보 조회
func GetWholePriceLineGraph(stdItemCode string, examinItemCodes []string) WholePriceGraph {
	db := common.GetDB()

	var list []WholePrice
	var graph WholePriceGraph
	var rangeLabels []string
	var data = map[string]map[string]map[string]map[string][]wholePriceGraphData{}
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
			ExaminUnitName, ExaminUnit, ExaminGradeName, ExaminGradeCode, Price FROM MAFRA_WHOLE_PRICE
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

		// 그래프 X(날짜), Y(가격)값 추출
		graphData := wholePriceGraphData{
			X: elem.ExaminDate,
			Y: elem.Price,
		}

		itemCode := elem.ExaminItemCode
		itemNameMap[itemCode] = elem.ExaminItemName
		speciesCode := elem.ExaminSpeciesCode
		speciesNameMap[speciesCode] = elem.ExaminSpeciesName
		unitName := elem.ExaminUnitName
		gradeCode := elem.ExaminGradeCode
		gradeNameMap[gradeCode] = elem.ExaminGradeName
		if data[itemCode] == nil {
			data[itemCode] = make(map[string]map[string]map[string][]wholePriceGraphData)
		}
		if data[itemCode][speciesCode] == nil {
			data[itemCode][speciesCode] = make(map[string]map[string][]wholePriceGraphData)
		}
		if data[itemCode][speciesCode][unitName] == nil {
			data[itemCode][speciesCode][unitName] = make(map[string][]wholePriceGraphData)
		}
		data[itemCode][speciesCode][unitName][gradeCode] = append(data[itemCode][speciesCode][unitName][gradeCode], graphData)
	}

	// json으로 출력하기 위해 struct 형태로 변환
	var graphLines []wholePriceGraphLine
	for itemCode, valueMap := range data {
		for speciesCode, valueMap2 := range valueMap {
			for unitName, valueMap3 := range valueMap2 {
				for gradeCode, graphData := range valueMap3 {
					var graphLine wholePriceGraphLine
					graphLine.ExaminItemName = itemNameMap[itemCode]
					graphLine.ExaminSpeciesName = speciesNameMap[speciesCode]
					graphLine.ExaminUnitName = unitName
					graphLine.ExaminGradeName = gradeNameMap[gradeCode]
					graphLine.GraphData = graphData
					graphLines = append(graphLines, graphLine)
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
		return graphLines[i].ExaminGradeName < graphLines[j].ExaminGradeName
	})

	graph.StdItemCode = stdItemCode
	graph.RangeLabel = rangeLabels
	graph.GraphLine = graphLines

	return graph
}

// GetRecentWholePrice 최근 도매 가격 정보 조회
func GetRecentWholePrice(examinItemCodes []string) []WholePrice {
	db := common.GetDB()

	examinItemCodesPhrase := common.GetOrPhrase(
		common.GetStringDefaultSlice(len(examinItemCodes), "ExaminItemCode"),
		examinItemCodes,
	)
	var list []WholePrice
	sqlQuery := fmt.Sprintf(`
		SELECT ExaminDate, ExaminItemName, ExaminItemCode, ExaminSpeciesName, ExaminSpeciesCode, 
			ExaminUnitName, ExaminUnit, ExaminGradeName, ExaminGradeCode, Price FROM MAFRA_WHOLE_PRICE
		WHERE (%s) 
		AND ExaminDate = (
			SELECT ExaminDate FROM MAFRA_WHOLE_PRICE 
			WHERE (%s) 
			ORDER BY ExaminDate DESC 
			LIMIT 1
		)
	;`, examinItemCodesPhrase, examinItemCodesPhrase)

	db.Raw(sqlQuery).Scan(&list)

	return list
}

// GetPreviousYearWholePrice 전년 도매 가격 정보 조회
func GetPreviousYearWholePrice(examinItemCodes []string) []WholePrice {
	db := common.GetDB()

	durationMonth := 3
	loc, _ := time.LoadLocation("Asia/Seoul")
	previousYearDate := time.Now().In(loc).AddDate(-1, 0, 0)
	maxDate := time.Now().In(loc).AddDate(-1, durationMonth, 0)
	previousYearDateString := common.GetDateString(previousYearDate, "-")
	maxDateString := common.GetDateString(maxDate, "-")

	var list []WholePrice

	examinItemCodesPhrase := common.GetOrPhrase(
		common.GetStringDefaultSlice(len(examinItemCodes), "ExaminItemCode"),
		examinItemCodes,
	)
	sqlQuery := fmt.Sprintf(`
		SELECT ExaminDate, ExaminItemName, ExaminItemCode, ExaminSpeciesName, ExaminSpeciesCode, 
			ExaminUnitName, ExaminUnit, ExaminGradeName, ExaminGradeCode, Price FROM MAFRA_WHOLE_PRICE
		WHERE (%s)
		AND ExaminDate  = (
			SELECT ExaminDate FROM MAFRA_WHOLE_PRICE
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
