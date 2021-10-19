package v2

import (
	"agripa-api/common"
)

// Question 문의사항
type Question struct {
	RegiDate string `gorm:"column:RegiDate" form:"regiDate" json:"regiDate"`
	Name     string `gorm:"column:Name" form:"name" json:"name"`
	Email    string `gorm:"column:Email" form:"email" json:"email"`
	Question string `gorm:"column:Question" form:"question" json:"question"`
}

// TableName Question 테이블 명
func (Question) TableName() string {
	return "AGRI_QUESTION"
}

// InsertQuestion 문의사항 등록
func InsertQuestion(question Question) error {
	db := common.GetDB()

	result := db.Create(&question)

	return result.Error
}
