package domain

import "time"

// User 领域对象,是DDD中的聚合根
// BO(business object)
type User struct {
	Id       int64
	Email    string
	Password string
	Phone    string
	//Addr Address
	Ctime time.Time
}

type UserInfo struct {
	Id              int64
	NickName        string
	BrithDays       string
	PersonalProfile string
}
type Address struct {
}
