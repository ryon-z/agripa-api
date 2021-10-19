package v2

import (
	"agripa-api/common"
	"fmt"
	"strings"
)

// Blog 블로그 구조체
type Blog struct {
	ItemCode    string `gorm:"column:ItemCode"`
	Query       string `gorm:"column:Query"`
	Title       string `gorm:"column:Title"`
	Link        string `gorm:"column:Link"`
	Bloggerlink string `gorm:"column:Bloggerlink"`
	Description string `gorm:"column:Description"`
	Bloggername string `gorm:"column:Bloggername"`
	Postdate    string `gorm:"column:Postdate"`
	Priority    int    `gorm:"column:Priority"`
	IsDisplay   int    `gorm:"column:IsDisplay"`
}

// TableName 블로그 테이블 이름
func (Blog) TableName() string {
	return "AGRIPA.AGRI_BLOG"
}

// News 뉴스기사 구조체
type News struct {
	ItemCode    string `gorm:"column:ItemCode"`
	Query       string `gorm:"column:Query"`
	Title       string `gorm:"column:Title"`
	Link        string `gorm:"column:Link"`
	Press       string `gorm:"column:Press"`
	Description string `gorm:"column:Description"`
	PubDate     string `gorm:"column:PubDate"`
	Priority    int    `gorm:"column:Priority"`
	IsDisplay   int    `gorm:"column:IsDisplay"`
}

// TableName 뉴스 테이블 이름
func (News) TableName() string {
	return "AGRIPA.AGRI_NEWS"
}

// AllYoutube 가락동365 유튜브 비디오 구조체
type AllYoutube struct {
	VideoID      string `gorm:"column:VideoID"`
	Query        string `gorm:"column:Query"`
	ChannelID    string `gorm:"column:ChannelID"`
	ChannelTitle string `gorm:"column:ChannelTitle"`
	Title        string `gorm:"column:Title"`
	Description  string `gorm:"column:Description"`
	ThumbnailURL string `gorm:"column:ThumbnailURL"`
	PublishedAt  string `gorm:"column:PublishedAt"`
	IsDisplay    int    `gorm:"column:IsDisplay"`
}

// TableName 가락동 364 유튜브 영상 테이블 이름
func (AllYoutube) TableName() string {
	return "AGRIPA.AGRI_YOUTUBE"
}

// Youtube 유튜브 비디오 구조체
type Youtube struct {
	ItemCode     string `gorm:"column:ItemCode"`
	VideoID      string `gorm:"column:VideoID"`
	Query        string `gorm:"column:Query"`
	ChannelID    string `gorm:"column:ChannelID"`
	ChannelTitle string `gorm:"column:ChannelTitle"`
	Title        string `gorm:"column:Title"`
	Description  string `gorm:"column:Description"`
	ThumbnailURL string `gorm:"column:ThumbnailURL"`
	PublishedAt  string `gorm:"column:PublishedAt"`
	IsDisplay    int    `gorm:"column:IsDisplay"`
}

// TableName 유튜브 영상 테이블 이름
func (Youtube) TableName() string {
	return "AGRIPA.AGRI_YOUTUBE"
}

// Query 쿼리 구조체
type Query struct {
	Query string `gorm:"column:Query"`
}

// getCount record 수를 리턴
func getCount(sql string) int {
	db := common.GetDB()

	var count int
	db.Raw(sql).Row().Scan(&count)

	return count
}

// GetMediaQuery examinItemCode와 매칭되는 쿼리 조회
func GetMediaQuery(usage string, examinItemCode string) []string {
	db := common.GetDB()

	var queries []Query
	var result []string

	sql := fmt.Sprintf(
		`SELECT DISTINCT Query FROM AGRIPA.AGRI_QUERY 
		WHERE %s = "%s" and ItemCode = %s;
		`, "`Usage`", usage, examinItemCode)
	db.Raw(sql).Scan(&queries)

	for _, query := range queries {
		result = append(result, query.Query)
	}

	return result
}

// getMediaQueryWherePhrase query가 포함된 문자열만 리턴하는 조건 생성
func getMediaQueryWherePhrase(queries []string) string {
	queryWherePhrase := "["
	isEnd := func(i int) bool { return i == len(queries)-1 }
	for i, query := range queries {
		if isEnd(i) {
			queryWherePhrase += fmt.Sprintf("%s]", query)
		} else {
			queryWherePhrase += fmt.Sprintf("%s|", query)
		}
	}

	if len(queries) == 0 {
		return "1=2"
	}
	return fmt.Sprintf(` m.Title REGEXP ("%s") `, queryWherePhrase)
}

// getExaminItemCodesWherePhrase itemCode where 절 생성
func getExaminItemCodesWherePhrase(columnNames []string, examinItemCodes []string) string {
	var phrases []string
	for i := range examinItemCodes {
		columnName := columnNames[i]
		itemCode := examinItemCodes[i]
		phrases = append(phrases, fmt.Sprintf("%s = %s", columnName, itemCode))
	}

	if len(examinItemCodes) == 0 {
		return "1=2"
	}

	return strings.Join(phrases, " OR ")
}

// getSameColumnNames 같은 이름의 컬럼명 slice 획득
func getSameColumnNames(columName string, length int) []string {
	var result []string
	for i := 0; i < length; i++ {
		result = append(result, columName)
	}

	return result
}

// GetBlogList 블로그 리스트 조회
func GetBlogList(examinItemCodes []string, start int, countPerPage int, baseDateTime string, queries []string) []Blog {
	db := common.GetDB()
	mediaTableName := Blog{}.TableName()
	usage := "blog"
	columnNames := getSameColumnNames("q.ItemCode", len(examinItemCodes))
	examinItemCodesPhrase := getExaminItemCodesWherePhrase(columnNames, examinItemCodes)
	queryWherePhrase := getMediaQueryWherePhrase(queries)

	var list []Blog

	sql := fmt.Sprintf(
		`SELECT q.ItemCode, m.Query, m.Title, m.Link, m.Bloggerlink, 
				m.Description, m.Bloggername, m.Postdate, m.Priority, m.IsDisplay
		FROM %s as m 
		JOIN AGRIPA.AGRI_QUERY as q 
		ON m.Query = q.Query 
		WHERE m.IsDisplay = 1 AND %s = "%s"
		AND (%s OR m.query = "000공통000") 
		AND m.Postdate <= "%s" AND m.Postdate > DATE_SUB("%s", INTERVAL 6 MONTH) 
		AND %s
		ORDER BY m.Postdate DESC 
		LIMIT %d, %d;
		`, mediaTableName, "q.`Usage`", usage, examinItemCodesPhrase, baseDateTime,
		baseDateTime, queryWherePhrase, start, countPerPage)
	db.Raw(sql).Scan(&list)

	return list
}

// GetBlogCount 블로그 레코드 수 조회
func GetBlogCount(examinItemCodes []string, baseDateTime string, queries []string) int {
	columnNames := getSameColumnNames("q.ItemCode", len(examinItemCodes))
	examinItemCodesPhrase := getExaminItemCodesWherePhrase(columnNames, examinItemCodes)
	queryWherePhrase := getMediaQueryWherePhrase(queries)

	sql := fmt.Sprintf(
		`SELECT count(q.ItemCode) 
		FROM %s as m 
		JOIN AGRIPA.AGRI_QUERY as q 
		ON m.Query = q.Query 
		WHERE m.IsDisplay = 1 AND %s = "blog" 
		AND (%s OR m.query = "000공통000") 
		AND m.Postdate <= "%s" AND m.Postdate > DATE_SUB("%s", INTERVAL 6 MONTH)
		AND %s;
		`, Blog{}.TableName(), "q.`Usage`", examinItemCodesPhrase, baseDateTime,
		baseDateTime, queryWherePhrase)

	return getCount(sql)
}

// GetNewsList 뉴스 리스트 조회
func GetNewsList(examinItemCodes []string, start int, countPerPage int, baseDateTime string, queries []string) []News {
	db := common.GetDB()
	mediaTableName := News{}.TableName()
	usage := "news"
	columnNames := getSameColumnNames("q.ItemCode", len(examinItemCodes))
	examinItemCodesPhrase := getExaminItemCodesWherePhrase(columnNames, examinItemCodes)
	queryWherePhrase := getMediaQueryWherePhrase(queries)

	var list []News

	sql := fmt.Sprintf(
		`SELECT q.ItemCode, m.Query, m.Title, m.Link, p.Name as Press, 
				m.Description, m.PubDate, m.Priority, m.IsDisplay
		from %s as m 
		JOIN AGRIPA.AGRI_QUERY as q 
		ON m.Query = q.Query 
		JOIN (SELECT Keyword, Name FROM AGRIPA.AGRI_PRESS WHERE IsDisplay = 1) as p 
		ON m.PressKeyword = p.Keyword 
		WHERE m.IsDisplay = 1 AND %s = "%s" 
		AND (%s OR m.query = "000공통000") 
		AND m.PubDate <= "%s" AND m.PubDate > DATE_SUB("%s", INTERVAL 6 MONTH)
		AND %s
		ORDER BY m.PubDate DESC 
		LIMIT %d, %d;
		`, mediaTableName, "q.`Usage`", usage, examinItemCodesPhrase, baseDateTime,
		baseDateTime, queryWherePhrase, start, countPerPage)
	db.Raw(sql).Scan(&list)

	return list
}

// GetNewsCount 뉴스 레코드 수 조회
func GetNewsCount(examinItemCodes []string, baseDateTime string, queries []string) int {
	columnNames := getSameColumnNames("q.ItemCode", len(examinItemCodes))
	examinIemCodesPhrase := getExaminItemCodesWherePhrase(columnNames, examinItemCodes)
	queryWherePhrase := getMediaQueryWherePhrase(queries)

	sql := fmt.Sprintf(
		`SELECT count(q.ItemCode)
		from %s as m 
		JOIN AGRIPA.AGRI_QUERY as q 
		ON m.Query = q.Query 
		JOIN (SELECT Keyword, Name FROM AGRIPA.AGRI_PRESS WHERE IsDisplay = 1) as p 
		ON m.PressKeyword = p.Keyword 
		WHERE m.IsDisplay = 1 AND %s = "news" 
		AND (%s OR m.query = "000공통000") 
		AND m.PubDate <= "%s" AND m.PubDate > DATE_SUB("%s", INTERVAL 6 MONTH)
		AND %s;
		`, News{}.TableName(), "q.`Usage`", examinIemCodesPhrase, baseDateTime,
		baseDateTime, queryWherePhrase)

	return getCount(sql)
}

// GetYoutubeList 유튜브 비디오 리스트 조회
func GetYoutubeList(examinItemCodes []string, start int, countPerPage int, baseDateTime string, queries []string) []Youtube {
	db := common.GetDB()
	mediaTableName := Youtube{}.TableName()
	usage := "youtube"
	columnNames := getSameColumnNames("q.ItemCode", len(examinItemCodes))
	examinItemCodesPhrase := getExaminItemCodesWherePhrase(columnNames, examinItemCodes)
	queryWherePhrase := getMediaQueryWherePhrase(queries)

	var list []Youtube

	sql := fmt.Sprintf(
		`SELECT q.ItemCode, m.VideoID, m.Query, m.ChannelID, m.ChannelTitle, m.Title, 
				m.Description, m.ThumbnailURL, m.PublishedAt, m.Priority, m.IsDisplay
		FROM %s as m 
		JOIN AGRIPA.AGRI_QUERY as q 
		ON m.Query = q.Query 
		WHERE m.IsDisplay = 1 AND %s = "%s" 
		AND (%s or m.query = "000공통000") 
		AND m.PublishedAt <= "%s" AND m.PublishedAt > DATE_SUB("%s", INTERVAL 6 MONTH) 
		AND %s
		ORDER BY m.PublishedAt DESC 
		LIMIT %d, %d;
		`, mediaTableName, "q.`Usage`", usage, examinItemCodesPhrase, baseDateTime,
		baseDateTime, queryWherePhrase, start, countPerPage)
	db.Raw(sql).Scan(&list)
	fmt.Println(sql)

	return list
}

// GetYoutubeCount 유튜브 비디오 레코드 수 조회
func GetYoutubeCount(examinItemCodes []string, baseDateTime string, queries []string) int {
	columnNames := getSameColumnNames("q.ItemCode", len(examinItemCodes))
	examinIemCodesPhrase := getExaminItemCodesWherePhrase(columnNames, examinItemCodes)
	queryWherePhrase := getMediaQueryWherePhrase(queries)

	sql := fmt.Sprintf(
		`SELECT count(q.ItemCode)
		FROM %s as m 
		JOIN AGRIPA.AGRI_QUERY as q 
		ON m.Query = q.Query 
		WHERE m.IsDisplay = 1 AND %s = "youtube" 
		AND (%s or m.query = "000공통000") 
		AND m.PublishedAt <= "%s" AND m.PublishedAt > DATE_SUB("%s", INTERVAL 6 MONTH) 
		AND %s
		`, Youtube{}.TableName(), "q.`Usage`", examinIemCodesPhrase, baseDateTime,
		baseDateTime, queryWherePhrase)

	return getCount(sql)
}
