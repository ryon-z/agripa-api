package common

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// IsStringInArray target string이 입력 받은 배열에 포함되는지 확인
func IsStringInArray(target string, arr []string) bool {
	for _, elem := range arr {
		if elem == target {
			return true
		}
	}

	return false
}

// GetOrPhrase SQLQuery에 사용되는 OR 문 생성
func GetOrPhrase(columnNames []string, data interface{}) string {
	var result []string
	switch data.(type) {
	case []string:
		for i, value := range data.([]string) {
			row := fmt.Sprintf("%s = \"%s\"", columnNames[i], value)
			result = append(result, row)
		}
	}

	if len(result) == 0 {
		return "1=2"
	}

	return strings.Join(result, " OR ")
}

// GetStringDefaultSlice 값을 가진 slice 얻기
func GetStringDefaultSlice(length int, defaultValue string) []string {
	result := make([]string, length)
	for i := range result {
		result[i] = defaultValue
	}

	return result
}

// ErrorInfo 에러 정보 구조체
type ErrorInfo struct {
	Status  bool
	Message string
}

// ConvertStringToInt 숫자string을 int로 변경하여 리턴
func ConvertStringToInt(inputs map[string]string) (map[string]int, map[string]string) {
	var results map[string]int
	results = make(map[string]int)
	var errors map[string]string
	errors = make(map[string]string)

	for name, value := range inputs {
		result, err := strconv.Atoi(value)
		if err != nil {
			errors[name] = fmt.Sprintf(
				"%s가 숫자형이 아닙니다. %s: %s",
				name, name, value,
			)
		} else {
			errors[name] = ""
		}
		results[name] = result
	}
	return results, errors
}

// IsError 에러인지 판별하여 에러 정보 구조체 리턴
func IsError(inputError interface{}) ErrorInfo {
	var result ErrorInfo
	result.Status = false
	result.Message = ""

	switch inputError.(type) {
	case map[string]string:
		for _, errorString := range inputError.(map[string]string) {
			if errorString != "" {
				result.Status = true
				result.Message = errorString
			}
		}
	case error:
		if inputError != nil {
			result.Status = true
			result.Message = inputError.(error).Error()
		}
	}

	return result
}

// GetDateString 입력 받은 구분자를 사용하여 만든 날짜 문자열 리턴
func GetDateString(date time.Time, sep string) string {
	return fmt.Sprintf("%02d%s%02d%s%02d",
		date.Year(), sep, date.Month(), sep, date.Day())
}

// IsDateString 올바른 날짜 문자열 인지 확인
func IsDateString(candidate string, sep string) bool {
	if len(candidate) != 10 {
		return false
	}

	if strings.Contains(sep, candidate) {
		return false
	}

	splited := strings.Split(candidate, sep)
	if len(splited) != 3 {
		return false
	}

	if len(splited[0]) != 4 &&
		len(splited[1]) != 2 &&
		len(splited[1]) != 2 {
		return false
	}

	return true
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// IsValidEmail 이메일 유효성 확인
func IsValidEmail(email string) bool {

	if len(email) < 3 && len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}
