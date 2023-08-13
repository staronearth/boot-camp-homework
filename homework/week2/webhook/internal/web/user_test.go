package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func TestEncrypt(t *testing.T) {
	bgpgyte, err := bcrypt.GenerateFromPassword([]byte("gaolin@123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword(bgpgyte, []byte("gaolin@123"))
	assert.NoError(t, err)
}

func TestBirthDays(t *testing.T) {
	birthday := "4000-12-31" // 替换为要验证的生日日期
	if isValidBirthday(birthday) {
		fmt.Println("生日有效")
	} else {
		fmt.Println("生日无效")
	}
}
func isValidBirthday(date string) bool {
	layout := "2006-01-02" // 指定日期格式

	// 解析日期字符串
	birthday, err := time.Parse(layout, date)
	if err != nil {
		return false
	}

	// 获取当前日期
	currentDate := time.Now()

	// 比较生日是否在当前日期之前
	if birthday.Before(currentDate) {
		return true
	}

	return false
}
