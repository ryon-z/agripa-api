package controllers

import "github.com/gin-gonic/gin"

// APIResMeta API 응답 메타 정보
type APIResMeta struct {
	RequestMethod string
	RequestURL    string
	//RequestForm     map[string][]string
	//RequestPostForm map[string][]string
	//RequestParam    gin.Params
	RequestHost string

	UserAgent string
	Referer   string
}

//SetAPIResMeta 메타 정보 설정
func SetAPIResMeta(c *gin.Context) *APIResMeta {
	var meta APIResMeta

	meta.RequestMethod = c.Request.Method
	meta.RequestURL = c.Request.RequestURI
	//meta.RequestForm = c.Request.Form
	//meta.RequestPostForm = c.Request.PostForm
	//meta.RequestParam = c.Params
	meta.RequestHost = c.Request.Host

	meta.UserAgent = c.Request.UserAgent()
	meta.Referer = c.Request.Referer()

	return &meta
}
