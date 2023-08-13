package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
	ErrUserInfoNotFound   = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}
func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	//dao.db.WithContext(ctx).First(&u, "email=?", email).Error
	return u, err
}
func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	//存储毫秒数
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			//邮箱冲突
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (dao *UserDAO) FindUserTableById(ctx context.Context, idx int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id=?", idx).First(&u).Error
	return u, err
}

func (dao *UserDAO) FindUserInfoTableById(ctx context.Context, idx int64) (UserInfo, error) {
	var ui UserInfo
	err := dao.db.WithContext(ctx).Where("id=?", idx).First(&ui).Error
	return ui, err
}

func (dao *UserDAO) InsertUserInfo(ctx context.Context, uinfo UserInfo) error {
	//存储毫秒数
	now := time.Now().UnixMilli()
	uinfo.Ctime = now
	uinfo.Utime = now
	err := dao.db.WithContext(ctx).Create(&uinfo).Error
	return err
}

func (dao *UserDAO) UpdateUserInfo(ctx context.Context, uinfo UserInfo) error {
	now := time.Now().UnixMilli()
	uinfo.Ctime = now
	uinfo.Utime = now
	err := dao.db.WithContext(ctx).Where("id=?", uinfo.Id).Updates(&uinfo).Error
	return err
}

// User User直接对应数据库表,
// 有些叫做entity,有些叫做model,有些叫做PO(persistent object)
type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 全局唯一
	Email    string `gorm:"unique"`
	Password string

	//往这里添加
	//创建时间,毫秒数
	Ctime int64
	//更新时间,毫秒数
	Utime int64
}

type UserInfo struct {
	Id int64 `gorm:"unique"`

	NickName        string
	BrithDays       string
	PersonalProfile string

	Ctime int64
	//更新时间,毫秒数
	Utime int64
}
type UserDetail struct {
}

//type Address struct {
//	Id     int64
//	UserId int64
//}
