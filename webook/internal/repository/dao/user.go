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
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	// 存毫秒数
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	// 保持链路 WithContext
	err := dao.db.WithContext(ctx).Create(&u).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			// 邮箱错误
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (dao *UserDAO) FindById(ctx context.Context, Id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", Id).First(&u).Error
	return u, err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

// User 直接对应数据库表结构，有些人叫做entity, 有些人叫做 model, 也有些人叫PO(Persistent Object)
type User struct {
	Id int64 `gorm:"primaryKey, autoIncrement"`
	// 全部用户邮箱唯一
	Email    string `gorm:"unique"`
	Password string

	// 创建时间，毫秒数
	Ctime int64
	// 更新时间，毫秒数
	Utime int64
}
