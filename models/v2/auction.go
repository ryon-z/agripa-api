package v2

import (
	"agripa-api/common"
	"fmt"
	"sort"
)

// AdjAuctionQuantity 정산 경락 거래량
type AdjAuctionQuantity struct {
	AuctionDate    string `gorm:"column:AuctionDate"`
	StdItemName    string `gorm:"column:StdItemName"`
	StdItemCode    string `gorm:"column:StdItemCode"`
	StdSpeciesName string `gorm:"column:StdSpeciesName"`
	StdSpeciesCode string `gorm:"column:StdSpeciesCode"`
	Quantity       string `gorm:"column:Quantity"`
}

// AdjAuctionQuantityBarGraph 그래프용 정산 경락 바차트  거래량
type AdjAuctionQuantityBarGraph struct {
	StdItemCode string
	RangeLabel  []string
	GraphBar    []adjAuctionQuantityGraphBar
}

type adjAuctionQuantityGraphBar struct {
	StdItemName    string
	StdSpeciesName string
	GraphData      []string
}

// AdjAuctionQuantityRecentTop3 정산 경락 거래량 최근 TOP 3
type AdjAuctionQuantityRecentTop3 struct {
	AuctionDate string `gorm:"column:AuctionDate"`
	StdItemName string `gorm:"column:StdItemName"`
	StdItemCode string `gorm:"column:StdItemCode"`
	AccQy       string `gorm:"column:AccQy"`
}

// GetAdjAuctionQuantityBarGraph 정산 경락 거래량 바 차트
func GetAdjAuctionQuantityBarGraph(stdItemCode string) AdjAuctionQuantityBarGraph {
	db := common.GetDB()

	var list []AdjAuctionQuantity
	var graph AdjAuctionQuantityBarGraph
	var rangeLabels []string
	var data = map[string]map[string]map[string]string{}
	var itemNameMap = map[string]string{}
	var speciesNameMap = map[string]string{}

	durationMonth := 1
	sqlQuery := fmt.Sprintf(`
		SELECT AuctionDate, StdItemName, StdItemCode, StdSpeciesName, 
			StdSpeciesCode, Quantity
		FROM MAFRA_ADJ_AUCTION_QUANTITY
		WHERE StdItemCode = "%s"
		AND AuctionDate > DATE_SUB(NOW(), INTERVAL %d MONTH)
		ORDER BY AuctionDate, StdItemCode, StdSpeciesCode
	;`, stdItemCode, durationMonth)

	db.Raw(sqlQuery).Scan(&list)

	for _, elem := range list {
		// 고유한 날짜(그래프의 X 축 라벨)만 추출
		if !common.IsStringInArray(elem.AuctionDate, rangeLabels) {
			rangeLabels = append(rangeLabels, elem.AuctionDate)
		}
	}

	for _, elem := range list {
		itemCode := elem.StdItemCode
		itemNameMap[itemCode] = elem.StdItemName
		speciesCode := elem.StdSpeciesCode
		speciesNameMap[speciesCode] = elem.StdSpeciesName
		if data[itemCode] == nil {
			data[itemCode] = make(map[string]map[string]string)
		}
		if data[itemCode][speciesCode] == nil {
			data[itemCode][speciesCode] = make(map[string]string)
			for _, date := range rangeLabels {
				data[itemCode][speciesCode][date] = "0"
			}
		}
		data[itemCode][speciesCode][elem.AuctionDate] = elem.Quantity
	}

	// json으로 출력하기 위해 struct 형태로 변환
	var graphBars []adjAuctionQuantityGraphBar
	for itemCode, valueMap := range data {
		for speciesCode, valueMap2 := range valueMap {
			var graphBar adjAuctionQuantityGraphBar
			graphBar.StdItemName = itemNameMap[itemCode]
			graphBar.StdSpeciesName = speciesNameMap[speciesCode]
			for index := range rangeLabels {
				date := rangeLabels[index]
				quantity := valueMap2[date]
				graphBar.GraphData = append(graphBar.GraphData, quantity)
			}
			graphBars = append(graphBars, graphBar)
		}
	}

	sort.Slice(graphBars, func(i, j int) bool {
		if graphBars[i].StdItemName < graphBars[j].StdItemName {
			return true
		}
		if graphBars[i].StdItemName > graphBars[j].StdItemName {
			return false
		}
		return graphBars[i].StdSpeciesName < graphBars[j].StdSpeciesName
	})

	graph.StdItemCode = stdItemCode
	graph.RangeLabel = rangeLabels
	graph.GraphBar = graphBars

	return graph
}

// GetAdjAuctionQuantityRecentTop3 정산 경락 거래량 최근 TOP 3
func GetAdjAuctionQuantityRecentTop3() []AdjAuctionQuantityRecentTop3 {
	db := common.GetDB()

	var list []AdjAuctionQuantityRecentTop3
	sqlQuery := fmt.Sprintf(`
		SELECT AuctionDate, StdItemName, StdItemCode, TRUNCATE(sum(Quantity) * 0.00110231, 3) AS AccQy 
		FROM MAFRA_ADJ_AUCTION_QUANTITY AS A 
		WHERE AuctionDate = (
			SELECT AuctionDate FROM MAFRA_ADJ_AUCTION_QUANTITY 
			ORDER BY AuctionDate DESC LIMIT 1)
		GROUP BY AuctionDate, StdItemCode 
		ORDER BY AccQy desc
		LIMIT 3
	;`)

	db.Raw(sqlQuery).Scan(&list)

	return list
}
