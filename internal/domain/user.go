package domain

import "time"

type User struct {
	Id       int64
	Email    string
	Password string

	Nickname string
	Birthday time.Time // YYYY-MM-DD
	AboutMe  string

	Phone string
	// timezone，UTC 0 millisecond
	CreatedAt time.Time
	UpdatedAt time.Time

	//WechatInfo WechatInfo
}

type Address struct {
	Province string
	Region   string
}

//func (u User) ValidateEmail() bool {
// 在这里用正则表达式校验
//return u.Email
//}

// TodayIsBirthday
func (u User) TodayIsBirthday() bool {
	now := time.Now()
	return now.Month() == u.Birthday.Month() && now.Day() == u.Birthday.Day()
}
