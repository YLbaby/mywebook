package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"context"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicate
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   dao.UserDAO
	cache cache.RedisUserCache
}

func NewUserRepository(dao *dao.GORMUserDAO, c cache.RedisUserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	// repository负责dao层对象和domain对象的转换
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})

	// 操作缓存
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	// 先查缓存
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		// 必然有数据
		return u, nil
	}

	// 再查数据库
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = domain.User{
		Id:    ue.Id,
		Email: ue.Email,
	}

	// 开一个协程将用户信息写入缓存
	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			// 这里怎么办？
			// 打日志，做监控
		}
	}()
	return u, err
}
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	// repository负责dao层对象和domain对象的转换
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, err
}

func (r *CachedCodeRepository) domainToEntity(u domain.User) dao.User {

}
