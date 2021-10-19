package v2

import (
	"agripa-api/common"
	"fmt"
	"sort"
)

// TradeInfo 수출입 정보
type TradeInfo struct {
	HskPrdlstCode string `gorm:"column:HskPrdlstCode"`
	BaseDate      string `gorm:"column:BaseDate"`
	Weight        string `gorm:"column:Weight"`
	Amount        string `gorm:"column:Amount"`
}

// TradeInfoGraph 그래프용 수출입 정보
type TradeInfoGraph struct {
	StdItemCode string
	RangeLabel  []string
	GraphLine   []tradeInfoGraphLine
}

type tradeInfoGraphLine struct {
	HskPrdlstCode string
	TradeType     string
	GraphData     []tradeInfoGraphData
}

type tradeInfoGraphData struct {
	X string
	Y string
}

// TradeRecentTop3 수출입 최근 중량 TOP 3 정보
type TradeRecentTop3 struct {
	StdItemCode     string `gorm:"column:StdItemCode"`
	StdItemName     string `gorm:"column:StdItemName"`
	HskPrdlstCode   string `gorm:"column:HskPrdlstCode"`
	BaseDate        string `gorm:"column:BaseDate"`
	ConvertedWeight string `gorm:"column:ConvertedWeight"`
}

// GetImportationInfo 수입 정보
func GetImportationInfo(stdItemCode string, hskPrdlstCodes []string) TradeInfoGraph {
	db := common.GetDB()

	var list []TradeInfo
	var graph TradeInfoGraph
	var rangeLabels []string
	var data = map[string]map[string][]tradeInfoGraphData{}

	hskPrdlstCodesPhrase := common.GetOrPhrase(
		common.GetStringDefaultSlice(len(hskPrdlstCodes), "HskPrdlstCode"),
		hskPrdlstCodes,
	)
	sqlQuery := fmt.Sprintf(`
		SELECT HskPrdlstCode, BaseDate, Weight, Amount
		FROM TRADE
		WHERE (%s)
		AND TradeType = "수입"
		ORDER BY BaseDate
	;`, hskPrdlstCodesPhrase)

	db.Raw(sqlQuery).Scan(&list)

	for _, elem := range list {
		// 고유한 날짜(그래프의 X 축 라벨)만 추출
		if !common.IsStringInArray(elem.BaseDate, rangeLabels) {
			rangeLabels = append(rangeLabels, elem.BaseDate)
		}

		hskPrdlstCode := elem.HskPrdlstCode
		for tradeType, value := range map[string]string{"중량(kg)": elem.Weight, "금액(달러)": elem.Amount} {
			// 그래프 X(날짜), Y(수출입 정보)값 추출
			graphData := tradeInfoGraphData{
				X: elem.BaseDate,
				Y: value,
			}
			if data[hskPrdlstCode] == nil {
				data[hskPrdlstCode] = make(map[string][]tradeInfoGraphData)
			}
			data[hskPrdlstCode][tradeType] = append(data[hskPrdlstCode][tradeType], graphData)
		}
	}

	// json으로 출력하기 위해 struct 형태로 변환
	var graphLines []tradeInfoGraphLine
	for hskPrdlstCode, valueMap := range data {
		for tradeType, graphData := range valueMap {
			var graphLine tradeInfoGraphLine
			graphLine.HskPrdlstCode = hskPrdlstCode
			graphLine.TradeType = tradeType
			graphLine.GraphData = graphData
			graphLines = append(graphLines, graphLine)
		}
	}

	sort.Slice(graphLines, func(i, j int) bool {
		if graphLines[i].HskPrdlstCode < graphLines[j].HskPrdlstCode {
			return true
		}
		if graphLines[i].HskPrdlstCode > graphLines[j].HskPrdlstCode {
			return false
		}

		return graphLines[i].TradeType < graphLines[j].TradeType
	})

	graph.StdItemCode = stdItemCode
	graph.RangeLabel = rangeLabels
	graph.GraphLine = graphLines

	return graph
}

// GetExportationInfo 수출 정보
func GetExportationInfo(stdItemCode string, hskPrdlstCodes []string) TradeInfoGraph {
	db := common.GetDB()

	var list []TradeInfo
	var graph TradeInfoGraph
	var rangeLabels []string
	var data = map[string]map[string][]tradeInfoGraphData{}

	hskPrdlstCodesPhrase := common.GetOrPhrase(
		common.GetStringDefaultSlice(len(hskPrdlstCodes), "HskPrdlstCode"),
		hskPrdlstCodes,
	)
	sqlQuery := fmt.Sprintf(`
		SELECT HskPrdlstCode, BaseDate, Weight, Amount
		FROM TRADE
		WHERE (%s)
		AND TradeType = "수출"
		ORDER BY BaseDate
	;`, hskPrdlstCodesPhrase)

	db.Raw(sqlQuery).Scan(&list)

	for _, elem := range list {
		// 고유한 날짜(그래프의 X 축 라벨)만 추출
		if !common.IsStringInArray(elem.BaseDate, rangeLabels) {
			rangeLabels = append(rangeLabels, elem.BaseDate)
		}

		hskPrdlstCode := elem.HskPrdlstCode
		for tradeType, value := range map[string]string{"중량(kg)": elem.Weight, "금액(달러)": elem.Amount} {
			// 그래프 X(날짜), Y(수출입 정보)값 추출
			graphData := tradeInfoGraphData{
				X: elem.BaseDate,
				Y: value,
			}
			if data[hskPrdlstCode] == nil {
				data[hskPrdlstCode] = make(map[string][]tradeInfoGraphData)
			}
			data[hskPrdlstCode][tradeType] = append(data[hskPrdlstCode][tradeType], graphData)
		}
	}

	// json으로 출력하기 위해 struct 형태로 변환
	var graphLines []tradeInfoGraphLine
	for hskPrdlstCode, valueMap := range data {
		for tradeType, graphData := range valueMap {
			var graphLine tradeInfoGraphLine
			graphLine.HskPrdlstCode = hskPrdlstCode
			graphLine.TradeType = tradeType
			graphLine.GraphData = graphData
			graphLines = append(graphLines, graphLine)
		}
	}

	sort.Slice(graphLines, func(i, j int) bool {
		if graphLines[i].HskPrdlstCode < graphLines[j].HskPrdlstCode {
			return true
		}
		if graphLines[i].HskPrdlstCode > graphLines[j].HskPrdlstCode {
			return false
		}

		return graphLines[i].TradeType < graphLines[j].TradeType
	})

	graph.StdItemCode = stdItemCode
	graph.RangeLabel = rangeLabels
	graph.GraphLine = graphLines

	return graph
}

// GetTradeRecentTop3  수출입 최근 중량 TOP 3 정보 조회
func GetTradeRecentTop3(tradeType string) []TradeRecentTop3 {
	if !common.IsStringInArray(tradeType, []string{"수출", "수입"}) {
		return []TradeRecentTop3{}
	}

	db := common.GetDB()

	var list []TradeRecentTop3
	sqlQuery := fmt.Sprintf(`
		SELECT B.StdItemCode, C.ItemName AS StdItemName, A.*
		FROM (
			SELECT HskPrdlstCode, BaseDate, TRUNCATE(sum(Weight) * 0.00110231, 3) AS ConvertedWeight
			FROM TRADE
			WHERE BaseDate = (SELECT BaseDate FROM TRADE order by BaseDate DESC LIMIT 1)
			AND TradeType = "%s"
			GROUP BY HskPrdlstCode
			ORDER BY Weight DESC
			LIMIT 6
		) AS A
		JOIN ITEM_MAPPING AS B
		ON A.HskPrdlstCode = B.HskPrdlstCode
		JOIN STD_ITEM_CODE AS C
		ON B.StdItemCode = C.ItemCode
		GROUP BY A.HskPrdlstCode
		ORDER BY ConvertedWeight DESC 
		LIMIT 3
	;`, tradeType)

	db.Raw(sqlQuery).Scan(&list)

	return list
}
