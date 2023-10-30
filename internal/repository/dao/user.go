package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("Email Already Exists")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	UpdateById(ctx context.Context, entity User) error
	FindById(ctx context.Context, uid int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	//FindByWechat(ctx context.Context, openId string) (User, error)
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

func (dao *GORMUserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.CreatedAt = now
	u.UpdatedAt = now

	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// email is duplicated (already exists)
			return ErrDuplicateEmail
		}
	}
	return err
}

func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

//func (dao *GORMUserDAO) FindByWechat(ctx context.Context, openId string) (User, error) {
//	var u User
//	err := dao.db.WithContext(ctx).Where("wechat_open_id=?", openId).First(&u).Error
//	return u, err
//}

func (dao *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone=?", phone).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) FindById(ctx context.Context, uid int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id=?", uid).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) UpdateById(ctx context.Context, entity User) error {
	return dao.db.WithContext(ctx).Model(&entity).Where("id=?", entity.Id).
		Updates(map[string]any{
			"updated_at": time.Now().UnixMilli(), // be careful about the naming convention GORM applies when mapping Go structs to database tables!
			"nickname":   entity.Nickname,
			"birthday":   entity.Birthday,
			"about_me":   entity.AboutMe,
		}).Error
}

type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	//Email    string `gorm:"unique"`
	//Email is a nullable column
	Email    sql.NullString `gorm:"unique"`
	Password string
	Nickname string         `gorm:"type=varchar(128)"`
	Birthday int64          // YYYY-MM-DD
	AboutMe  string         `gorm:"type=varchar(4096)"`
	Phone    sql.NullString `gorm:"unique"`

	// timezone，UTC 0 millisecond
	CreatedAt int64 //`gorm:"column:createdat"`
	UpdatedAt int64 //`gorm:"column:updatedat"`
}
