package models

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

// getQuery itemCode와 매칭되는 쿼리 조회
func getQuery(itemCode int) []Query {
	db := common.GetDB()

	var queries []Query

	sql := fmt.Sprintf(
		`SELECT Query FROM AGRIPA.AGRI_QUERY 
		WHERE %s = "youtube" and ItemCode = %d;
		`, "`Usage`", itemCode)
	db.Raw(sql).Scan(&queries)

	return queries
}

// getQueryWherePhrase query가 포함된 문자열만 리턴하는 조건 생성
func getQueryWherePhrase(queries []Query) string {
	var queryWherePhrases []string
	for _, query := range queries {
		phrase := fmt.Sprintf(` m.Title REGEXP ("%s") `, query.Query)
		queryWherePhrases = append(queryWherePhrases, phrase)
	}
	queryWherePhrase := strings.Join(queryWherePhrases, "OR")
	if queryWherePhrase != "" {
		queryWherePhrase = "AND (" + queryWherePhrase + ")"
	}

	return queryWherePhrase
}

// GetBlogList 블로그 리스트 조회
func GetBlogList(itemCode int, start int, countPerPage int, baseDateTime string) []Blog {
	db := common.GetDB()
	mediaTableName := Blog{}.TableName()
	queries := getQuery(itemCode)
	queryWherePhrase := getQueryWherePhrase(queries)

	var list []Blog

	sql := fmt.Sprintf(
		`SELECT q.ItemCode, m.Query, m.Title, m.Link, m.Bloggerlink, 
				m.Description, m.Bloggername, m.Postdate, m.Priority, m.IsDisplay
		FROM %s as m 
		JOIN AGRIPA.AGRI_QUERY as q 
		ON m.Query = q.Query 
		WHERE m.IsDisplay = 1 AND %s = "blog" 
		AND (q.ItemCode = %d or m.query = "000공통000") 
		AND m.Postdate <= "%s" AND m.Postdate > DATE_SUB("%s", INTERVAL 6 MONTH) 
		%s
		ORDER BY m.Postdate DESC 
		LIMIT %d, %d;
		`, mediaTableName, "q.`Usage`", itemCode, baseDateTime,
		baseDateTime, queryWherePhrase, start, countPerPage)
	db.Raw(sql).Scan(&list)

	return list
}

// GetBlogCount 블로그 레코드 수 조회
func GetBlogCount(itemCode int, baseDateTime string) int {
	queries := getQuery(itemCode)
	queryWherePhrase := getQueryWherePhrase(queries)

	sql := fmt.Sprintf(
		`SELECT count(q.ItemCode) 
		FROM %s as m 
		JOIN AGRIPA.AGRI_QUERY as q 
		ON m.Query = q.Query 
		WHERE m.IsDisplay = 1 AND %s = "blog" 
		AND (q.ItemCode = %d or m.query = "000공통000") 
		AND m.Postdate <= "%s" AND m.Postdate > DATE_SUB("%s", INTERVAL 6 MONTH)
		%s;
		`, Blog{}.TableName(), "q.`Usage`", itemCode, baseDateTime,
		baseDateTime, queryWherePhrase)

	return getCount(sql)
}

// GetNewsList 뉴스 리스트 조회
func GetNewsList(itemCode int, start int, countPerPage int, baseDateTime string) []News {
	db := common.GetDB()
	mediaTableName := News{}.TableName()
	queries := getQuery(itemCode)
	queryWherePhrase := getQueryWherePhrase(queries)

	var list []News

	sql := fmt.Sprintf(
		`SELECT q.ItemCode, m.Query, m.Title, m.Link, p.Name as Press, 
				m.Description, m.PubDate, m.Priority, m.IsDisplay
		from %s as m 
		JOIN AGRIPA.AGRI_QUERY as q 
		ON m.Query = q.Query 
		JOIN (SELECT Keyword, Name FROM AGRIPA.AGRI_PRESS WHERE IsDisplay = 1) as p 
		ON m.PressKeyword = p.Keyword 
		WHERE m.IsDisplay = 1 AND %s = "news" 
		AND (q.ItemCode = %d or m.query = "000공통000") 
		AND m.PubDate <= "%s" AND m.PubDate > DATE_SUB("%s", INTERVAL 6 MONTH)
		%s
		ORDER BY m.PubDate DESC 
		LIMIT %d, %d;
		`, mediaTableName, "q.`Usage`", itemCode, baseDateTime,
		baseDateTime, queryWherePhrase, start, countPerPage)
	db.Raw(sql).Scan(&list)

	return list
}

// GetNewsCount 뉴스 레코드 수 조회
func GetNewsCount(itemCode int, baseDateTime string) int {
	queries := getQuery(itemCode)
	queryWherePhrase := getQueryWherePhrase(queries)

	sql := fmt.Sprintf(
		`SELECT count(q.ItemCode)
		from %s as m 
		JOIN AGRIPA.AGRI_QUERY as q 
		ON m.Query = q.Query 
		JOIN (SELECT Keyword, Name FROM AGRIPA.AGRI_PRESS WHERE IsDisplay = 1) as p 
		ON m.PressKeyword = p.Keyword 
		WHERE m.IsDisplay = 1 AND %s = "news" 
		AND (q.ItemCode = %d or m.query = "000공통000") 
		AND m.PubDate <= "%s" AND m.PubDate > DATE_SUB("%s", INTERVAL 6 MONTH)
		%s;
		`, News{}.TableName(), "q.`Usage`", itemCode, baseDateTime,
		baseDateTime, queryWherePhrase)

	return getCount(sql)
}

// GetAllYoutubeList 모든 유튜브 비디오 리스트 조회
func GetAllYoutubeList(start int, countPerPage int, query string, baseDateTime string) []AllYoutube {
	db := common.GetDB()
	mediaTableName := AllYoutube{}.TableName()
	queryWherePhrase := ""
	if query != "" {
		queryWherePhrase = fmt.Sprintf(`AND Query = "%s" `, query)
	}

	var list []AllYoutube

	sql := fmt.Sprintf(
		`SELECT VideoID, Query, ChannelID, ChannelTitle, Title, 
				Description, ThumbnailURL, PublishedAt, Priority, IsDisplay
		FROM %s 
		WHERE IsDisplay = 1 
		%s 
		AND PublishedAt <= "%s"
		ORDER BY PublishedAt DESC  
		LIMIT %d, %d;
		`, mediaTableName, queryWherePhrase, baseDateTime, start, countPerPage)
	db.Raw(sql).Scan(&list)

	return list
}

// GetAllYoutubeCount 모든 유튜브 비디오 레코드 수 조회
func GetAllYoutubeCount(query string, baseDateTime string) int {
	queryWherePhrase := ""
	if query != "" {
		queryWherePhrase = fmt.Sprintf(`AND Query = "%s" `, query)
	}

	sql := fmt.Sprintf(
		`SELECT count(VideoID)
		FROM %s 
		WHERE IsDisplay = 1 
		%s 
		AND PublishedAt <= "%s";
		`, AllYoutube{}.TableName(), queryWherePhrase, baseDateTime)

	return getCount(sql)
}

// GetYoutubeList 유튜브 비디오 리스트 조회
func GetYoutubeList(itemCode int, start int, countPerPage int, baseDateTime string) []Youtube {
	db := common.GetDB()
	mediaTableName := Youtube{}.TableName()

	var list []Youtube

	queries := getQuery(itemCode)
	var queryWherePhrases []string
	for _, query := range queries {
		phrase := fmt.Sprintf(` m.Title REGEXP ("%s") `, query.Query)
		queryWherePhrases = append(queryWherePhrases, phrase)
	}
	queryWherePhrase := strings.Join(queryWherePhrases, "OR")
	if queryWherePhrase != "" {
		queryWherePhrase = "AND (" + queryWherePhrase + ")"
	}

	sql := fmt.Sprintf(
		`SELECT q.ItemCode, m.VideoID, m.Query, m.ChannelID, m.ChannelTitle, m.Title, 
				m.Description, m.ThumbnailURL, m.PublishedAt, m.Priority, m.IsDisplay
		FROM %s as m 
		JOIN AGRIPA.AGRI_QUERY as q 
		ON m.Query = q.Query 
		WHERE m.IsDisplay = 1 AND %s = "youtube" 
		AND (q.ItemCode = %d or m.query = "000공통000") 
		AND m.PublishedAt <= "%s" AND m.PublishedAt > DATE_SUB("%s", INTERVAL 6 MONTH) 
		%s
		ORDER BY m.PublishedAt DESC 
		LIMIT %d, %d;
		`, mediaTableName, "q.`Usage`", itemCode, baseDateTime,
		baseDateTime, queryWherePhrase, start, countPerPage)
	db.Raw(sql).Scan(&list)

	return list
}

// GetYoutubeCount 유튜브 비디오 레코드 수 조회
func GetYoutubeCount(itemCode int, baseDateTime string) int {
	queries := getQuery(itemCode)
	queryWherePhrase := getQueryWherePhrase(queries)

	sql := fmt.Sprintf(
		`SELECT count(q.ItemCode)
		FROM %s as m 
		JOIN AGRIPA.AGRI_QUERY as q 
		ON m.Query = q.Query 
		WHERE m.IsDisplay = 1 AND %s = "youtube" 
		AND (q.ItemCode = %d or m.query = "000공통000") 
		AND m.PublishedAt <= "%s" AND m.PublishedAt > DATE_SUB("%s", INTERVAL 6 MONTH) 
		%s
		`, Youtube{}.TableName(), "q.`Usage`", itemCode, baseDateTime,
		baseDateTime, queryWherePhrase)

	return getCount(sql)
}
